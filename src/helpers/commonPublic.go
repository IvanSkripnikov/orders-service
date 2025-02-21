package helpers

import (
	"net/http"
	"time"
)

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func FormatResponse(w http.ResponseWriter, httpStatus int, category string) {
	w.WriteHeader(httpStatus)

	data := ResponseData{
		"error": "Unsuccessfull request",
	}
	SendResponse(w, data, category, httpStatus)
}
