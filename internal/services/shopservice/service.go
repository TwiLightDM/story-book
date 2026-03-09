package shopservice

import (
	"context"
	"fmt"
	"story-book/internal/entities"

	"github.com/google/uuid"
)

type ShopRepository interface {
	Create(ctx context.Context, shop *entities.Shop) error
	ReadAll(ctx context.Context) ([]entities.Shop, error)
	ReadById(ctx context.Context, id string) (*entities.Shop, error)
	Update(ctx context.Context, shop *entities.Shop) (*entities.Shop, error)
	Delete(ctx context.Context, id string) error
}

type shopService struct {
	repo ShopRepository
}

func NewShopService(repo ShopRepository) ShopService {
	return &shopService{repo: repo}
}

func (s *shopService) CreateShop(ctx context.Context, shop *entities.Shop) (*entities.Shop, error) {
	shop.Id = uuid.NewString()

	err := s.repo.Create(ctx, shop)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

func (s *shopService) ReedShopById(ctx context.Context, id string) (*entities.Shop, error) {
	shop, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

func (s *shopService) ReadShops(ctx context.Context) ([]entities.Shop, error) {
	fmt.Println(uuid.NewString())
	shops, err := s.repo.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return shops, nil
}

func (s *shopService) UpdateShop(ctx context.Context, shop *entities.Shop) (*entities.Shop, error) {
	var err error

	updatedShop, err := s.repo.Update(ctx, shop)
	if err != nil {
		return nil, err
	}

	return updatedShop, nil
}

func (s *shopService) DeleteShop(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
