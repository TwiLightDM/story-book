package bookservice

import (
	"context"
	"errors"
	"fmt"
	"story-book/internal/entities"

	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, book *entities.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *bookRepository) ReadAll(ctx context.Context, limit, offset int) ([]entities.Book, error) {
	var books []entities.Book

	if err := r.db.
		WithContext(ctx).
		Model(&entities.Book{}).
		Select("books.*, COALESCE(AVG(ratings.stars), 0) as rating").
		Joins("LEFT JOIN ratings ON ratings.book_id = books.id").
		Group("books.id").
		Limit(limit).
		Offset(offset).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("genre", "book_id")
		}).
		Find(&books).Error; err != nil {

		return nil, err
	}

	return books, nil
}

func (r *bookRepository) ReadById(ctx context.Context, id string) (*entities.Book, error) {
	var book entities.Book

	if err := r.db.
		WithContext(ctx).
		Model(&entities.Book{}).
		Select("books.*, COALESCE(AVG(ratings.stars), 0) as rating").
		Joins("LEFT JOIN ratings ON ratings.book_id = books.id").
		Where("books.id = ?", id).
		Group("books.id").
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("genre", "book_id")
		}).
		First(&book).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	fmt.Println(book.Rating)
	fmt.Println("///////")

	return &book, nil
}

func (r *bookRepository) Update(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	var updatedBook entities.Book
	res := r.db.
		WithContext(ctx).
		Model(&entities.Book{}).
		Where("id = ?", book.Id).
		Updates(book).
		Scan(&updatedBook)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}

		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, ErrBookNotFound
	}

	return &updatedBook, nil
}

func (r *bookRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&entities.Book{Id: id})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrBookNotFound
	}

	return nil
}
