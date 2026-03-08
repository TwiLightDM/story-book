package genreservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"gorm.io/gorm"
)

type genreRepository struct {
	db *gorm.DB
}

func NewGenreRepository(db *gorm.DB) GenreRepository {
	return &genreRepository{db: db}
}

func (r *genreRepository) Create(ctx context.Context, genre *entities.Genre) error {
	return r.db.WithContext(ctx).Create(genre).Error
}

func (r *genreRepository) ReadByGenreAndBookId(ctx context.Context, genre, bookId string) (*entities.Genre, error) {
	var gen entities.Genre
	if err := r.db.
		WithContext(ctx).
		Where("genre = ? and book_id = ?", genre, bookId).
		First(&gen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGenreNotFound
		}
		return nil, err
	}
	return &gen, nil
}

func (r *genreRepository) Delete(ctx context.Context, genre, bookId string) error {
	result := r.db.WithContext(ctx).
		Where("genre = ? and book_id = ?", genre, bookId).
		Delete(&entities.Genre{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrGenreNotFound
	}

	return nil
}
