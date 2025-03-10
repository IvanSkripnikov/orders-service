package helpers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"orders-service/models"

	"github.com/IvanSkripnikov/go-gormdb"
	"github.com/IvanSkripnikov/go-logger"
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

func GetOrdersByUser(w http.ResponseWriter, r *http.Request) {
	category := "/v1/orders/get-by-user"
	var orders []models.Order

	userID, _ := getIDFromRequestString(strings.TrimSpace(r.URL.Path))
	if userID == 0 {
		FormatResponse(w, http.StatusUnprocessableEntity, category)
		return
	}

	db := gormdb.GetClient(models.ServiceDatabase)
	err := db.Where("user_id = ?", userID).Find(&orders).Error
	if checkError(w, err, category) {
		return
	}

	data := ResponseData{
		"response": orders,
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

	response := "success"
	newOrderParams := models.OrderParams{UserID: order.UserID, Price: order.Price, ItemID: order.ItemID, Volume: order.Volume, OrderID: order.ID}
	if orderCreateSaga(newOrderParams) {
		response = "failure"
	} else {
		err = db.Model(&order).Update("status", models.StatusCompleted).Error
		if checkError(w, err, category) {
			response = "failure"
		}
	}

	data := ResponseData{
		"response": response,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func orderCreateSaga(orderParams models.OrderParams) bool {
	// 1. списываем средаства у пользователя
	newPayment := models.PaymentParams{UserID: orderParams.UserID, Amount: orderParams.Price}
	paymentResponse, err := CreateQueryWithScalarResponse(http.MethodPut, Config.PaymentServiceUrl+"/v1/account/buy", newPayment)
	if paymentResponse != models.Success || err != nil {
		logger.Errorf("Unsuccessful payment: %v", newPayment)

		// отправить сообщение в redis
		messageData := map[string]interface{}{
			"title":       "Unsuccessful payment",
			"description": "Failure payment " + strconv.FormatFloat(float64(orderParams.Price), 'f', -1, 32),
			"user":        orderParams.UserID,
			"category":    "deal",
		}
		SendNotification(messageData)
		return false
	}

	// 2. проверяем наличие товара на складе
	newItemBook := models.BookingItem{ItemID: orderParams.ItemID, Volume: orderParams.Volume}
	itemBookResponse, err := CreateQueryWithScalarResponse(http.MethodPost, Config.WarehouseServiceUrl+"/v1/warehouses/book-item", newItemBook)
	if itemBookResponse != models.Success || err != nil {
		logger.Errorf("Unsuccessful booking: %v", newItemBook)

		// в случае неудачи делаем возврат средств
		_, err = CreateQueryWithScalarResponse(http.MethodPut, Config.PaymentServiceUrl+"/v1/account/rollback", newPayment)
		if err != nil {
			logger.Errorf("Unsuccessful rollback: %v", err)
		}

		// отправить сообщение в redis
		messageData := map[string]interface{}{
			"title":       "Not enough items on warehouse",
			"description": "Requested item now not available on our warehouses. Please try again later",
			"user":        orderParams.UserID,
			"category":    "deal",
		}
		SendNotification(messageData)

		return false
	}

	// 3. проверяем наличие свободных курьеров
	newCourierBook := models.BookingCourier{OrderID: orderParams.OrderID}
	courierBookResponse, err := CreateQueryWithScalarResponse(http.MethodPost, Config.DeliveryServiceUrl+"/v1/couriers/book", newCourierBook)
	if courierBookResponse != models.Success || err != nil {
		logger.Errorf("Unsuccessful booking delivery: %v", newCourierBook)

		// в случае неудачи делаем возврат средств и снимаем бронь со склада
		_, err = CreateQueryWithScalarResponse(http.MethodPut, Config.PaymentServiceUrl+"/v1/account/rollback", newPayment)
		if err != nil {
			logger.Errorf("Unsuccessful rollback: %v", err)
		}
		_, err = CreateQueryWithScalarResponse(http.MethodPost, Config.WarehouseServiceUrl+"/v1/warehouses/rollback-book", newItemBook)
		if err != nil {
			logger.Errorf("Unsuccessful rollback booking: %v", err)
		}

		// отправить сообщение в redis
		messageData := map[string]interface{}{
			"title":       "Not enough couriers for delivery",
			"description": "All couriers is busy now. Please try again later",
			"user":        orderParams.UserID,
			"category":    "deal",
		}
		SendNotification(messageData)

		return false
	}

	return true
}
