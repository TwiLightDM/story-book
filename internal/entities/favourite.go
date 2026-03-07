package entities

import (
	"time"

	"gorm.io/gorm"
)

type Favourite struct {
	Id        string
	BookId    string
	UserId    string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
