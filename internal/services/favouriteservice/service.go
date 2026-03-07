package favouriteservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"github.com/google/uuid"
)

type FavouriteRepository interface {
	Create(ctx context.Context, favourite *entities.Favourite) error
	ReadAllByUserId(ctx context.Context, userId string, limit, offset int) ([]entities.Book, error)
	ReadByUserIdAndBookId(ctx context.Context, userId, bookId string) (*entities.Favourite, error)
	Delete(ctx context.Context, userId, bookId string) error
}

type favouriteService struct {
	repo FavouriteRepository
}

func NewFavouriteService(repo FavouriteRepository) FavouriteService {
	return &favouriteService{repo: repo}
}

func (s *favouriteService) CreateFavourite(ctx context.Context, favourite *entities.Favourite) error {
	fav, err := s.repo.ReadByUserIdAndBookId(ctx, favourite.UserId, favourite.BookId)
	if err != nil {
		if !errors.Is(err, ErrFavouriteNotFound) {
			return err
		}
	}

	if fav != nil {
		return ErrFavouriteAlreadyExists
	}

	favourite.Id = uuid.NewString()

	err = s.repo.Create(ctx, favourite)
	if err != nil {
		return err
	}

	return nil
}

func (s *favouriteService) ReadFavourites(ctx context.Context, userId string, limit, offset int) ([]entities.Book, error) {
	books, err := s.repo.ReadAllByUserId(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *favouriteService) DeleteFavourite(ctx context.Context, userId, bookId string) error {
	err := s.repo.Delete(ctx, userId, bookId)
	if err != nil {
		return err
	}

	return nil
}
