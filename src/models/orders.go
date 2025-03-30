package models

const StatusNew = 1
const StatusCompleted = 2
const StatusCanceled = 3

type Order struct {
	ID        int     `gorm:"index;type:int" json:"id"`
	UserID    int     `gorm:"index;type:int" json:"userId"`
	ItemID    int     `gorm:"index;type:int" json:"itemId"`
	Volume    int     `gorm:"type:int" json:"volume"`
	Price     float32 `gorm:"type:float" json:"price"`
	Created   string  `gorm:"type:text" json:"created"`
	Updated   string  `gorm:"type:text" json:"updated"`
	Status    int     `gorm:"index;type:tinyint;default:1" json:"status"`
	RequestID string  `gorm:"index;type:string" json:"requestId"`
}

func (s Order) TableName() string { return "orders" }

type UniqueOrder struct {
	ID        int    `gorm:"index;type:int" json:"id"`
	RequestID string `gorm:"index;type:string" json:"requestId"`
	Response  string `gorm:"type:string" json:"response"`
}

func (s UniqueOrder) TableName() string { return "unique_orders" }

type OrderRequest struct {
	ID          int     `json:"id"`
	UserID      int     `json:"userId"`
	ItemID      int     `json:"itemId"`
	Volume      int     `json:"volume"`
	Price       float32 `json:"price"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
	Status      int     `json:"status"`
	RequestID   string  `json:"requestId"`
	Promocode   string  `json:"promocode"`
	Certificate string  `json:"certificate"`
}

type OrderResponse struct {
	ID          int     `json:"id"`
	UserID      int     `json:"userId"`
	ItemID      int     `json:"itemId"`
	Volume      int     `json:"volume"`
	Price       float32 `json:"price"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
	Status      int     `json:"status"`
	RequestID   string  `json:"requestId"`
	Promocode   string  `json:"promocode"`
	Certificate string  `json:"certificate"`
	LoyaltyID   []int   `json:"loyaltyId"`
}
