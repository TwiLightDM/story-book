package encryptservice

import "errors"

var (
	ErrGeneratingSalt  = errors.New("error generating salt")
	ErrHashingPassword = errors.New("error hashing password")
	ErrInvalidPassword = errors.New("invalid password")
)
