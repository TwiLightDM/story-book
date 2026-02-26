package bookservice

import (
	"context"
	"story-book/internal/entities"

	"github.com/google/uuid"
)

type BookRepository interface {
	Create(ctx context.Context, book *entities.Book) error
	ReadAll(ctx context.Context) ([]entities.Book, error)
	ReadById(ctx context.Context, id string) (*entities.Book, error)
	Update(ctx context.Context, book *entities.Book) (*entities.Book, error)
	Delete(ctx context.Context, id string) error
}

type bookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) CreateBook(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	book.Id = uuid.NewString()

	err := s.repo.Create(ctx, book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) ReedBookById(ctx context.Context, id string) (*entities.Book, error) {
	book, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) ReadBooks(ctx context.Context) ([]entities.Book, error) {
	books, err := s.repo.ReadAll(ctx)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *bookService) UpdateBook(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	var err error

	updatedBook, err := s.repo.Update(ctx, book)
	if err != nil {
		return nil, err
	}

	return updatedBook, nil
}

func (s *bookService) DeleteBook(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
