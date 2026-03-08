package ratingservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"gorm.io/gorm"
)

type ratingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) RatingRepository {
	return &ratingRepository{db: db}
}

func (r *ratingRepository) Create(ctx context.Context, rating *entities.Rating) error {
	return r.db.WithContext(ctx).Create(rating).Error
}

func (r *ratingRepository) ReadByUserIdAndBookId(ctx context.Context, userId, bookId string) (*entities.Rating, error) {
	var rating entities.Rating
	if err := r.db.
		WithContext(ctx).
		Where("user_id = ? and book_id = ?", userId, bookId).
		First(&rating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRatingNotFound
		}
		return nil, err
	}
	return &rating, nil
}

func (r *ratingRepository) Update(ctx context.Context, rating *entities.Rating) error {
	return r.db.
		WithContext(ctx).
		Model(&entities.Rating{}).
		Where("user_id = ? and book_id = ?", rating.UserId, rating.BookId).
		Updates(rating).Error
}

func (r *ratingRepository) Delete(ctx context.Context, userId, bookId string) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? and book_id = ?", userId, bookId).
		Delete(&entities.Rating{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrRatingNotFound
	}

	return nil
}
