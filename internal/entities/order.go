package entities

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	Id           string
	Status       string
	DeliveryType string
	Address      *string
	Cost         float64
	Points       int
	ShopId       *string
	UserId       string
	CodeId       *string
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt
	Books        []OrderBook `gorm:"foreignKey:OrderId"`
	Shop         Shop        `gorm:"foreignKey:ShopId;references:Id"`
}
