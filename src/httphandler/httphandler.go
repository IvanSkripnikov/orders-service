package httphandler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"orders-service/helpers"
	"orders-service/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitHTTPServer() {
	// подключение роутов
	http.HandleFunc("/", Serve)
	// подключение prometheus
	http.Handle("/auth/metrics", promhttp.Handler())

	err := http.ListenAndServe(":8080", nil) //nolint:gosec
	if err != nil {
		errMessage := fmt.Sprintf("Can't init HTTP server: %v", err)

		logger.Error(errMessage)
	}
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func Serve(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var allow []string
	found := false

	for _, routeUnit := range routes {
		matches := routeUnit.regex.FindStringSubmatch(strings.TrimSpace(r.URL.Path))

		if len(matches) > 0 {
			if r.Method != routeUnit.method {
				allow = append(allow, routeUnit.method)
				continue
			}
			found = true
			routeUnit.handler(w, r)
		}
	}

	defer func() {
		method := r.Method
		elapsed := time.Since(start).Seconds()
		helpers.RequestsByMethodTotal.WithLabelValues(method).Inc()
		helpers.RequestsTotal.Inc()
		helpers.RequestLatencySummary.Observe(elapsed)
		helpers.RequestLatencyHistogram.WithLabelValues("").Observe(elapsed)
	}()

	if !found && len(allow) == 0 {
		helpers.FormatResponse(w, http.StatusNotFound, "middleware")

		return
	}
}

func GetHTTPHandler() *http.ServeMux {
	httpHandler := http.NewServeMux()

	for _, routeUnit := range routes {
		httpHandler.HandleFunc(handleRegexp(routeUnit.regex), routeUnit.handler)
	}

	return httpHandler
}

func handleRegexp(regExp *regexp.Regexp) string {
	expr := regExp.String()[1 : len(regExp.String())-1]

	var result string

	if strings.Count(expr, "/") > 1 {
		parts := strings.Split(expr, "/")

		parts = parts[:len(parts)-1]

		result = strings.Join(parts, "/") + "/"
	} else {
		result = expr
	}

	return result
}
