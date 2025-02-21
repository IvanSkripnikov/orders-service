package httphandler

import (
	"net/http"
	"regexp"

	"orders-service/controllers"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

var routes = []route{
	// system
	newRoute(http.MethodGet, "/health", controllers.HealthCheck),
	// orders
	newRoute(http.MethodGet, "/v1/orders/list", controllers.GetOrdersListV1),
	newRoute(http.MethodGet, "/v1/orders/get/([0-9]+)", controllers.GetOrderV1),
	newRoute(http.MethodPost, "/v1/orders/create", controllers.CreateOrderV1),
}
