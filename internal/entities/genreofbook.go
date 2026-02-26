package entities

import (
	"time"

	"gorm.io/gorm"
)

type GenreOfBook struct {
	Id        string
	Genre     string
	BookId    string
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
