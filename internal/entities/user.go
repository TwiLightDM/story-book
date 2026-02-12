package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        string
	Name      string
	Surname   string
	Email     string
	Phone     string
	Password  string
	Salt      string
	Role      string
	Question  string
	Answer    string
	Points    int
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}
