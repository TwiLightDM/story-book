package entities

import (
	"time"

	"gorm.io/gorm"
)

type Code struct {
	Id            string
	Code          string
	Percent       int
	AmountOfUsage *int
	ExpiredAt     *time.Time
	CreatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
