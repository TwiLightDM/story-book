package userservice

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrWrongAnswer       = errors.New("wrong answer")
	ErrUserNotFounded    = errors.New("user not founded")
	ErrUserAlreadyExists = errors.New("user already exists")
)
