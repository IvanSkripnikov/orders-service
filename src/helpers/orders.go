package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"orders-service/models"

	"github.com/IvanSkripnikov/go-logger"
	"gorm.io/gorm"
)

func GetOrdersList(w http.ResponseWriter, _ *http.Request) {
	category := "/v1/orders/list"
	var orders []models.Order

	err := GormDB.Find(&orders).Error
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

	err := GormDB.Where("id = ?", order.ID).First(&order).Error
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

	err := GormDB.Where("user_id = ?", userID).Find(&orders).Error
	if checkError(w, err, category) && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	data := ResponseData{
		"response": orders,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	category := "/v1/orders/create"
	var orderRequest models.OrderRequest

	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if checkError(w, err, category) {
		return
	}

	orderRequest.Created = GetCurrentDate()
	orderRequest.Updated = GetCurrentDate()
	orderRequest.Status = models.StatusNew

	// рассчитываем возможную скидку
	resp, err := CreateQueryWithResponse(http.MethodPut, Config.LoyaltyServiceUrl+"/v1/loyalty/apply-for-order", orderRequest)
	if checkError(w, err, category) {
		return
	}
	orderResponse := getOrderFromResponse(resp)

	// проверяем, был ли такой заказ ранее
	var orderDB models.Order
	err = GormDB.Where("request_id = ?", orderRequest.RequestID).First(&orderDB).Error

	// такого заказа раньше не было, создаём новый, иначе возвращаем его результат
	var uniqueOrder models.UniqueOrder
	response := "success"
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Debugf("Create new order: %v", orderResponse)
		var order models.Order
		order.ID = orderResponse.ID
		order.UserID = orderResponse.UserID
		order.ItemID = orderResponse.ItemID
		order.Volume = orderResponse.Volume
		order.Price = orderResponse.Price
		order.Created = orderResponse.Created
		order.Updated = orderResponse.Updated
		order.Status = orderResponse.Status
		order.RequestID = orderResponse.RequestID

		err = GormDB.Create(&order).Error
		if checkError(w, err, category) {
			return
		}

		newOrderParams := models.OrderParams{UserID: order.UserID, Price: order.Price, ItemID: order.ItemID, Volume: order.Volume, OrderID: order.ID}
		if !orderCreateSaga(newOrderParams) {
			GormDB.Model(&order).Update("status", models.StatusCanceled)
			response = "failure"
		} else {
			err = GormDB.Model(&order).Update("status", models.StatusCompleted).Error
			if checkError(w, err, category) {
				response = "failure"
			}
		}

		// создаём запись с уникальным ID от пользователя
		uniqueOrder.RequestID = order.RequestID
		uniqueOrder.Response = response
		err = GormDB.Create(&uniqueOrder).Error
		if err != nil {
			logger.Errorf("Cant create unique order record %v", err)
		}

		_, err = CreateQueryWithResponse(http.MethodDelete, Config.LoyaltyServiceUrl+"/v1/loyalty/remove-for-user/"+strconv.Itoa(order.UserID), orderResponse.LoyaltyID)
		if checkError(w, err, category) {
			logger.Errorf("Cant deactivate used loyalty %v", err)
		}
		if orderRequest.Certificate != "" {
			_, err = CreateQueryWithResponse(http.MethodDelete, Config.LoyaltyServiceUrl+"/v1/loyalty/remove-certificate", orderRequest)
			if checkError(w, err, category) {
				logger.Errorf("Cant deactivate certificate %v", err)
			}
		}
	} else {
		logger.Debugf("Order already exists: %v", orderResponse)
		err = GormDB.Where("request_id = ?", orderRequest.RequestID).First(&uniqueOrder).Error
		if checkError(w, err, category) {
			return
		}
		response = uniqueOrder.Response
	}

	data := ResponseData{
		"response": response,
	}
	SendResponse(w, data, category, http.StatusOK)
}

func orderCreateSaga(orderParams models.OrderParams) bool {
	// 1. списываем средаства у пользователя
	newPayment := models.PaymentParams{UserID: orderParams.UserID, Amount: orderParams.Price}
	paymentResponse, err := CreateQueryWithResponse(http.MethodPut, Config.PaymentServiceUrl+"/v1/payment/pay", newPayment)
	if paymentResponse != models.Success || err != nil {
		logger.Errorf("Unsuccessful payment: %v, paymentRespoonse: %v, error: %v", newPayment, paymentResponse, err)

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
	itemBookResponse, err := CreateQueryWithResponse(http.MethodPost, Config.WarehouseServiceUrl+"/v1/warehouses/book-item", newItemBook)
	if itemBookResponse != models.Success || err != nil {
		logger.Errorf("Unsuccessful booking: %v, itemBookRespoonse: %v, error: %v", newItemBook, itemBookResponse, err)

		// в случае неудачи делаем возврат средств
		_, err = CreateQueryWithResponse(http.MethodPut, Config.PaymentServiceUrl+"/v1/payment/rollback", newPayment)
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
	courierBookResponse, err := CreateQueryWithResponse(http.MethodPost, Config.DeliveryServiceUrl+"/v1/couriers/book", newCourierBook)
	if courierBookResponse != models.Success || err != nil {
		logger.Errorf("Unsuccessful booking delivery: %v, courierBookRespoonse: %v, error: %v", newCourierBook, courierBookResponse, err)

		// в случае неудачи делаем возврат средств и снимаем бронь со склада
		_, err = CreateQueryWithResponse(http.MethodPut, Config.PaymentServiceUrl+"/v1/payment/rollback", newPayment)
		if err != nil {
			logger.Errorf("Unsuccessful rollback: %v", err)
		}
		_, err = CreateQueryWithResponse(http.MethodPost, Config.WarehouseServiceUrl+"/v1/warehouses/rollback-book", newItemBook)
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

func getOrderFromResponse(response any) models.OrderResponse {
	userMap, ok := response.(map[string]any)
	if !ok {
		logger.Errorf("Error asserting item to map[string]interface{}")
	}

	var loyaltyIDs = []int{}
	if userMap["loyaltyId"] != nil {
		loyaltyIDResponse := userMap["loyaltyId"].([]any)
		for _, loyalty := range loyaltyIDResponse {
			loyaltyIDs = append(loyaltyIDs, int(loyalty.(float64)))
		}
	}

	// Создаем экземпляр Order и заполняем его данными
	order := models.OrderResponse{
		ID:          int(userMap["id"].(float64)),
		UserID:      int(userMap["userId"].(float64)),
		ItemID:      int(userMap["itemId"].(float64)),
		Volume:      int(userMap["volume"].(float64)),
		Price:       float32(userMap["price"].(float64)),
		Created:     userMap["created"].(string),
		Updated:     userMap["updated"].(string),
		Status:      int(userMap["id"].(float64)),
		RequestID:   userMap["requestId"].(string),
		Promocode:   userMap["promocode"].(string),
		Certificate: userMap["certificate"].(string),
		LoyaltyID:   loyaltyIDs,
	}

	return order
}
