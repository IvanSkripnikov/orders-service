package helpers

import (
	"encoding/json"
	"net/http"
	"strings"

	"orders-service/models"

	"github.com/IvanSkripnikov/go-gormdb"
)

func GetOrdersList(w http.ResponseWriter, _ *http.Request) {
	category := "/v1/orders/list"
	var orders []models.Order

	db := gormdb.GetClient(models.ServiceDatabase)
	err := db.Find(&orders).Error
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"response": orders,
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

	if !isExists("SELECT * FROM orders WHERE id = ?", order.ID) {
		FormatResponse(w, http.StatusNotFound, category)
		return
	}

	db := gormdb.GetClient(models.ServiceDatabase)
	err := db.Where("id = ?", order.ID).First(&order).Error
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"response": order,
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

	currentTimestamp := GetCurrentTimestamp()
	order.Created = int(currentTimestamp)
	order.Updated = int(currentTimestamp)

	db := gormdb.GetClient(models.ServiceDatabase)
	err = db.Create(&order).Error
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"response": "success",
	}
	SendResponse(w, data, category, http.StatusOK)
}
