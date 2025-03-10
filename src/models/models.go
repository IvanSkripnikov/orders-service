package models

const ServiceDatabase = "OrdersService"
const Success = "success"

type PaymentParams struct {
	UserID int     `json:"userId"`
	Amount float32 `json:"amount"`
}

type BookingItem struct {
	ItemID int `json:"id"`
	Volume int `json:"volume"`
}

type BookingCourier struct {
	OrderID int `json:"orderId"`
}

type OrderParams struct {
	UserID  int     `json:"userId"`
	Price   float32 `json:"price"`
	ItemID  int     `json:"id"`
	Volume  int     `json:"volume"`
	OrderID int     `json:"orderId"`
}

type Redis struct {
	Address  string
	Port     string
	Password string
	DB       int
	Stream   string
}
