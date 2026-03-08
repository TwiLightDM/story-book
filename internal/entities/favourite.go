package entities

import (
	"time"
)

type Favourite struct {
	Id        string
	BookId    string
	UserId    string
	CreatedAt time.Time
}
