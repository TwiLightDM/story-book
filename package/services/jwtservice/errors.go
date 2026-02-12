package jwtservice

import "errors"

var (
	ErrLifetimeIsOver          = errors.New("lifetime is over")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrSignToken               = errors.New("error signing token")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
)
