package cartservice

import "errors"

var (
	ErrCartNotFound            = errors.New("cart not found")
	ErrBookInCartAlreadyExists = errors.New("book in cart already exists")
	ErrBookNotFound            = errors.New("book not found")
	ErrNotEnoughBooks          = errors.New("not enough books")
)
