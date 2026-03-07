package favouriteservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"gorm.io/gorm"
)

type favouriteRepository struct {
	db *gorm.DB
}

func NewFavouriteRepository(db *gorm.DB) FavouriteRepository {
	return &favouriteRepository{db: db}
}

func (r *favouriteRepository) Create(ctx context.Context, favourite *entities.Favourite) error {
	return r.db.WithContext(ctx).Create(favourite).Error
}

func (r *favouriteRepository) ReadAllByUserId(ctx context.Context, userId string, limit, offset int) ([]entities.Book, error) {
	var books []entities.Book

	subQuery := r.db.
		Model(&entities.Favourite{}).
		Select("book_id").
		Where("user_id = ? AND deleted_at IS NULL", userId)

	result := r.db.WithContext(ctx).
		Where("id IN (?) AND deleted_at IS NULL", subQuery).
		Limit(limit).
		Offset(offset).
		Find(&books)

	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil
}

func (r *favouriteRepository) ReadByUserIdAndBookId(ctx context.Context, userId, bookId string) (*entities.Favourite, error) {
	var favourite entities.Favourite
	if err := r.db.
		WithContext(ctx).
		Where("user_id = ? and book_id = ?", userId, bookId).
		First(&favourite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFavouriteNotFound
		}
		return nil, err
	}
	return &favourite, nil
}

func (r *favouriteRepository) Delete(ctx context.Context, userId, bookId string) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? and book_id = ?", userId, bookId).
		Delete(&entities.Favourite{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrFavouriteNotFound
	}

	return nil
}
