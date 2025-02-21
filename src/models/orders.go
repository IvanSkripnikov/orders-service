package models

type Order struct {
	ID        int     `json:"id"`
	UserID    int     `json:"userId"`
	ItemID    int     `json:"itemId"`
	Price     float32 `json:"price"`
	Created   int     `json:"created"`
	Updated   int     `json:"updated"`
	Completed uint8   `json:"completed"`
}
