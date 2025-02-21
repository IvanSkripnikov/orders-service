package helpers

import (
	"encoding/json"
	"net/http"
	"strings"

	"orders-service/logger"
	"orders-service/models"
)

func GetOrdersList(w http.ResponseWriter, _ *http.Request) {
	category := "/v1/orders/list"
	var orders []models.Order

	query := "SELECT id, user_id, item_id, price, created, updated, completed FROM orders WHERE id > 0"
	rows, err := DB.Query(query)
	if err != nil {
		logger.Error(err.Error())
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	for rows.Next() {
		order := models.Order{}
		if err = rows.Scan(&order.ID, &order.UserID, &order.ItemID, &order.Price, &order.Created, &order.Updated, &order.Completed); err != nil {
			logger.Error(err.Error())
			continue
		}
		orders = append(orders, order)
	}

	data := ResponseData{
		"data": orders,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	category := "/v1/orders/get"
	var order models.Order

	order.ID, _ = getIDFromRequestString(strings.TrimSpace(r.URL.Path))
	if order.ID == 0 {
		FormatResponse(w, http.StatusUnprocessableEntity, category)
		return
	}

	if !isExists("SELECT * FROM users WHERE id = ?", order.ID) {
		FormatResponse(w, http.StatusNotFound, category)
		return
	}

	query := "SELECT id, user_id, item_id, price, created, updated, completed FROM orders WHERE id = ?"
	rows, err := DB.Prepare(query)

	if checkError(w, err, category) {
		return
	}

	defer func() {
		_ = rows.Close()
	}()

	err = rows.QueryRow(order.ID).Scan(&order.ID, &order.UserID, &order.ItemID, &order.Price, &order.Created, &order.Updated, &order.Completed)
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"data": order,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	category := "/v1/orders/create"
	var order models.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if checkError(w, err, category) {
		return
	}

	query := "INSERT INTO users (user_id, item_id, price, created, updated) VALUES (?, ?, ?, ?, ?)"
	currentTimestamp := GetCurrentTimestamp()
	rows, err := DB.Query(query, order.UserID, order.ItemID, order.Price, currentTimestamp, currentTimestamp)

	if checkError(w, err, category) {
		return
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	data := ResponseData{
		"message": "Order successfully created!",
	}
	SendResponse(w, data, category, http.StatusOK)
}
