package cartservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(ctx context.Context, cart *entities.Cart) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var bookAmount int

	if err := tx.Model(&entities.Book{}).
		Select("amount").
		Where("id = ?", cart.BookId).
		Take(&bookAmount).Error; err != nil {

		tx.Rollback()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return err
	}

	if cart.Amount > bookAmount {
		tx.Rollback()
		return ErrNotEnoughBooks
	}

	if err := tx.Create(cart).Error; err != nil {
		tx.Rollback()

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrBookInCartAlreadyExists
		}
		return err
	}

	return tx.Commit().Error
}

func (r *cartRepository) ReadAllByUserId(ctx context.Context, userId string, limit, offset int) ([]entities.Book, error) {
	var books []entities.Book

	subQuery := r.db.
		Model(&entities.Cart{}).
		Select("book_id").
		Where("user_id = ?", userId)

	result := r.db.WithContext(ctx).
		Model(&entities.Book{}).
		Where("books.id IN (?) AND books.deleted_at IS NULL", subQuery).
		Select("books.*, COALESCE(AVG(ratings.stars), 0) as rating").
		Joins("LEFT JOIN ratings ON ratings.book_id = books.id").
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("genre", "book_id")
		}).
		Group("books.id").
		Limit(limit).
		Offset(offset).
		Find(&books)

	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil
}

func (r *cartRepository) Update(ctx context.Context, cart *entities.Cart) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var bookAmount int
	if err := tx.Model(&entities.Book{}).
		Select("amount").
		Where("id = ?", cart.BookId).
		Take(&bookAmount).Error; err != nil {

		tx.Rollback()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return err
	}

	if cart.Amount > bookAmount {
		tx.Rollback()
		return ErrNotEnoughBooks
	}

	result := tx.
		Model(&entities.Cart{}).
		Where("user_id = ? AND book_id = ?", cart.UserId, cart.BookId).
		Updates(cart)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return ErrCartNotFound
	}

	return tx.Commit().Error
}

func (r *cartRepository) Delete(ctx context.Context, userId, bookId string) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? and book_id = ?", userId, bookId).
		Delete(&entities.Cart{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrCartNotFound
	}

	return nil
}
