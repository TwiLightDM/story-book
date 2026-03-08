package genreservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"github.com/google/uuid"
)

type GenreRepository interface {
	Create(ctx context.Context, genre *entities.Genre) error
	ReadByGenreAndBookId(ctx context.Context, genre, bookId string) (*entities.Genre, error)
	Delete(ctx context.Context, genre, bookId string) error
}

type genreService struct {
	repo GenreRepository
}

func NewGenreService(repo GenreRepository) GenreService {
	return &genreService{repo: repo}
}

func (s *genreService) CreateGenre(ctx context.Context, genre *entities.Genre) error {
	gen, err := s.repo.ReadByGenreAndBookId(ctx, genre.Genre, genre.BookId)
	if err != nil {
		if !errors.Is(err, ErrGenreNotFound) {
			return err
		}
	}

	if gen != nil {
		return ErrGenreAlreadyExists
	}

	genre.Id = uuid.NewString()

	err = s.repo.Create(ctx, genre)
	if err != nil {
		return err
	}

	return nil
}

func (s *genreService) DeleteGenre(ctx context.Context, genre, bookId string) error {
	err := s.repo.Delete(ctx, genre, bookId)
	if err != nil {
		return err
	}

	return nil
}
