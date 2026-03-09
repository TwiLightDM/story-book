package shopservice

import "errors"

var (
	ErrShopNotFound      = errors.New("shop not found")
	ErrShopAlreadyExists = errors.New("shop already exists")
)
