package userservice

import "errors"

var (
	ErrWrongAnswer       = errors.New("wrong answer")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
