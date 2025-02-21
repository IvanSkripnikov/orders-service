package controllers

import (
	"net/http"

	"orders-service/helpers"
)

func GetOrdersListV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.GetOrdersList(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/orders/list")
	}
}

func GetOrderV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		helpers.GetOrder(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/orders/get")
	}
}

func CreateOrderV1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		helpers.CreateOrder(w, r)
	default:
		helpers.FormatResponse(w, http.StatusMethodNotAllowed, "/v1/orders/create")
	}
}
