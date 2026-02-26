package entities

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	Id          string
	Title       string
	Author      string
	Year        int
	Cost        float64
	Discount    *int
	Publisher   string
	Description *string
	Amount      int
	ImageData   []byte
	ImageMime   string
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
