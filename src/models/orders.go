package models

const StatusNew = 1
const StatusCompleted = 2
const StatusCanceled = 3

type Order struct {
	ID        int     `gorm:"index;type:int" json:"id"`
	UserID    int     `gorm:"index;type:int" json:"userId"`
	ItemID    int     `gorm:"index;type:int" json:"itemId"`
	Price     float32 `gorm:"type:float" json:"price"`
	Created   int     `gorm:"index;type:bigint" json:"created"`
	Updated   int     `gorm:"index;type:bigint" json:"updated"`
	Completed uint8   `gorm:"index;type:tinyint;default:1" json:"completed"`
}

func (s Order) TableName() string { return "orders" }
