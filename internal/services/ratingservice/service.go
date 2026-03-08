package ratingservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"github.com/google/uuid"
)

type RatingRepository interface {
	Create(ctx context.Context, rating *entities.Rating) error
	ReadByUserIdAndBookId(ctx context.Context, userId, bookId string) (*entities.Rating, error)
	Update(ctx context.Context, rating *entities.Rating) error
	Delete(ctx context.Context, rating, bookId string) error
}

type ratingService struct {
	repo RatingRepository
}

func NewRatingService(repo RatingRepository) RatingService {
	return &ratingService{repo: repo}
}

func (s *ratingService) CreateRating(ctx context.Context, rating *entities.Rating) error {
	r, err := s.repo.ReadByUserIdAndBookId(ctx, rating.UserId, rating.BookId)
	if err != nil {
		if !errors.Is(err, ErrRatingNotFound) {
			return err
		}
	}

	if r != nil {
		r.Stars = rating.Stars
		return s.repo.Update(ctx, r)
	}

	if rating.Stars < 0 || rating.Stars > 5 {
		return ErrUnsupportedRating
	}

	rating.Id = uuid.NewString()

	err = s.repo.Create(ctx, rating)
	if err != nil {
		return err
	}

	return nil
}

func (s *ratingService) DeleteRating(ctx context.Context, rating, bookId string) error {
	err := s.repo.Delete(ctx, rating, bookId)
	if err != nil {
		return err
	}

	return nil
}
