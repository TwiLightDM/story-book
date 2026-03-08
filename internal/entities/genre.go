package entities

import (
	"time"

	"gorm.io/gorm"
)

type Genre struct {
	Id        string
	Genre     string
	BookId    string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
