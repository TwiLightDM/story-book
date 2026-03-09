package shopservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"gorm.io/gorm"
)

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) Create(ctx context.Context, shop *entities.Shop) error {
	err := r.db.WithContext(ctx).Create(shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrShopAlreadyExists
		}
		return err
	}

	return nil
}

func (r *shopRepository) ReadAll(ctx context.Context) ([]entities.Shop, error) {
	var shops []entities.Shop
	if err := r.db.
		WithContext(ctx).
		Model(&entities.Shop{}).
		Find(&shops).Error; err != nil {

		return nil, err
	}

	return shops, nil
}

func (r *shopRepository) ReadById(ctx context.Context, id string) (*entities.Shop, error) {
	var shop entities.Shop

	if err := r.db.
		WithContext(ctx).
		Model(&entities.Shop{}).
		Where("id = ?", id).
		First(&shop).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShopNotFound
		}
		return nil, err
	}

	return &shop, nil
}

func (r *shopRepository) Update(ctx context.Context, shop *entities.Shop) (*entities.Shop, error) {
	var updatedShop entities.Shop
	res := r.db.
		WithContext(ctx).
		Model(&entities.Shop{}).
		Where("id = ?", shop.Id).
		Updates(shop).
		Scan(&updatedShop)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrShopNotFound
		}

		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, ErrShopNotFound
	}

	return &updatedShop, nil
}

func (r *shopRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&entities.Shop{Id: id})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrShopNotFound
	}

	return nil
}
