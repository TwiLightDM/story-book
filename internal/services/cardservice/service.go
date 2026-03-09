package cardservice

import (
	"context"
	"fmt"
	"story-book/internal/entities"

	"github.com/google/uuid"
)

type CardRepository interface {
	Create(ctx context.Context, card *entities.Card) error
	ReadAll(ctx context.Context) ([]entities.Card, error)
	Delete(ctx context.Context, id string) error
}

type cardService struct {
	repo CardRepository
}

func NewCardService(repo CardRepository) CardService {
	return &cardService{repo: repo}
}

func (s *cardService) CreateCard(ctx context.Context, card *entities.Card) (*entities.Card, error) {
	card.Id = uuid.NewString()

	err := s.repo.Create(ctx, card)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (s *cardService) ReadCards(ctx context.Context) ([]entities.Card, error) {
	fmt.Println(uuid.NewString())
	cards, err := s.repo.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (s *cardService) DeleteCard(ctx context.Context, numberOfCard string) error {
	err := s.repo.Delete(ctx, numberOfCard)
	if err != nil {
		return err
	}

	return nil
}
