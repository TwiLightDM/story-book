package cardservice

import "errors"

var (
	ErrCardNotFound      = errors.New("card not found")
	ErrCardAlreadyExists = errors.New("card already exists")
)
