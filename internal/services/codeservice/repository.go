package codeservice

import (
	"context"
	"story-book/internal/entities"

	"gorm.io/gorm"
)

type codeRepository struct {
	db *gorm.DB
}

func NewCodeRepository(db *gorm.DB) CodeRepository {
	return &codeRepository{db: db}
}

func (r *codeRepository) Create(ctx context.Context, code *entities.Code) error {
	return r.db.WithContext(ctx).Create(code).Error
}

func (r *codeRepository) ReadAll(ctx context.Context, limit, offset int) ([]entities.Code, error) {
	var codes []entities.Code
	if err := r.db.
		WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&codes).Error; err != nil {

		return nil, err
	}
	return codes, nil
}

func (r *codeRepository) Delete(ctx context.Context, code string) error {
	result := r.db.WithContext(ctx).Delete(&entities.Code{Code: code})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrCodeNotFound
	}

	return nil
}
