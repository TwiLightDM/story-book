package cartservice

import (
	"context"
	"story-book/internal/entities"

	"github.com/google/uuid"
)

type CartRepository interface {
	Create(ctx context.Context, cart *entities.Cart) error
	ReadAllByUserId(ctx context.Context, userId string, limit, offset int) ([]entities.Book, error)
	Update(ctx context.Context, cart *entities.Cart) error
	Delete(ctx context.Context, cart, bookId string) error
}

type cartService struct {
	repo CartRepository
}

func NewCartService(repo CartRepository) CartService {
	return &cartService{repo: repo}
}

func (s *cartService) CreateCart(ctx context.Context, cart *entities.Cart) error {
	cart.Id = uuid.NewString()

	err := s.repo.Create(ctx, cart)
	if err != nil {
		return err
	}

	return nil
}

func (s *cartService) ReadCarts(ctx context.Context, userId string, limit, offset int) ([]entities.Book, error) {
	books, err := s.repo.ReadAllByUserId(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *cartService) UpdateCart(ctx context.Context, cart *entities.Cart) error {
	return s.repo.Update(ctx, cart)
}

func (s *cartService) DeleteCart(ctx context.Context, userId, bookId string) error {
	err := s.repo.Delete(ctx, userId, bookId)
	if err != nil {
		return err
	}

	return nil
}
