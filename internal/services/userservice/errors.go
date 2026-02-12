package userservice

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrWrongAnswer  = errors.New("wrong answer")
)
