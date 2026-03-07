package favouriteservice

import "errors"

var (
	ErrFavouriteNotFound      = errors.New("favourites not found")
	ErrFavouriteAlreadyExists = errors.New("favourite already exists")
)
