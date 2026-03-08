package ratingservice

import "errors"

var (
	ErrRatingNotFound    = errors.New("rating not found")
	ErrUnsupportedRating = errors.New("unsupported rating")
)
