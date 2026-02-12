package validateservice

import "errors"

var (
	ErrBadEmail    = errors.New("bad email")
	ErrBadPassword = errors.New("bad password")
)
