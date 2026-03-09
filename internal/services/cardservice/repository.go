package cardservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"gorm.io/gorm"
)

type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) Create(ctx context.Context, card *entities.Card) error {
	err := r.db.WithContext(ctx).Create(card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrCardAlreadyExists
		}
	}

	return nil
}

func (r *cardRepository) ReadAll(ctx context.Context) ([]entities.Card, error) {
	var cards []entities.Card
	if err := r.db.
		WithContext(ctx).
		Model(&entities.Card{}).
		Find(&cards).Error; err != nil {

		return nil, err
	}

	return cards, nil
}

func (r *cardRepository) Delete(ctx context.Context, numberOfCard string) error {
	res := r.db.WithContext(ctx).Delete(&entities.Card{NumberOfCard: numberOfCard})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrCardNotFound
	}

	return nil
}
