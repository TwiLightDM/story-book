package orderservice

import (
	"context"
	"errors"
	"story-book/internal/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *entities.Order, code string) error {
	tx := r.db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var carts []entities.Cart
	if err := tx.
		Where("user_id = ?", order.UserId).
		Find(&carts).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(carts) == 0 {
		tx.Rollback()
		return ErrEmptyCart
	}

	totalCost := 0.0
	orderBooks := make([]entities.OrderBook, 0, len(carts))

	for _, cart := range carts {

		var book entities.Book

		if err := tx.
			Where("id = ?", cart.BookId).
			First(&book).Error; err != nil {
			tx.Rollback()
			return err
		}

		if book.Amount < cart.Amount {
			tx.Rollback()
			return ErrAmountOfBooks
		}

		price := book.Cost

		if book.Discount != nil {
			price = price - price*float64(*book.Discount)/100
		}

		totalCost += price * float64(cart.Amount)

		orderBooks = append(orderBooks, entities.OrderBook{
			Id:          uuid.NewString(),
			OrderId:     order.Id,
			BookId:      book.Id,
			Amount:      cart.Amount,
			PriceForOne: price,
		})

		if err := tx.Model(&entities.Book{}).
			Where("id = ?", book.Id).
			Update("amount", book.Amount-cart.Amount).Error; err != nil {

			tx.Rollback()
			return err
		}
	}

	if code != "" {
		var promo entities.Code

		if err := tx.
			Where("code = ?", code).
			First(&promo).
			Error; err != nil {
			tx.Rollback()
			return ErrCodeNotFound
		}

		if promo.ExpiredAt != nil && promo.ExpiredAt.Before(time.Now()) {
			tx.Rollback()
			return ErrCodeExpired
		}

		if promo.AmountOfUsage != nil {

			if *promo.AmountOfUsage == 0 {
				tx.Rollback()
				return ErrAmountOfUsage
			}

			newAmount := *promo.AmountOfUsage - 1
			order.CodeId = &promo.Id

			if err := tx.Model(&entities.Code{}).
				Where("id = ?", promo.Id).
				Update("amount_of_usage", newAmount).Error; err != nil {

				tx.Rollback()
				return err
			}
		}

		totalCost = totalCost - totalCost*float64(promo.Percent)/100

		order.CodeId = &promo.Id
	}

	order.Cost = totalCost
	order.Points = int(totalCost * 0.10)

	if err := tx.
		Create(order).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.
		Create(&orderBooks).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.
		Where("user_id = ?", order.UserId).
		Delete(&entities.Cart{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) ReadAll(ctx context.Context, limit, offset int, userId string, startDate, endDate time.Time) ([]entities.Order, error) {
	var orders []entities.Order

	if err := r.db.
		WithContext(ctx).
		Model(&entities.Order{}).
		Where("orders.user_id = ? and orders.created_at between ? and ?", userId, startDate, endDate).
		Limit(limit).
		Offset(offset).
		Preload("Shop").
		Preload("Books").
		Preload("Books.Book").
		Preload("Books.Book.Genres").
		Find(&orders).Error; err != nil {

		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) ReadById(ctx context.Context, id string) (*entities.Order, error) {
	var order entities.Order

	if err := r.db.
		WithContext(ctx).
		Model(&entities.Order{}).
		Where("orders.id = ?", id).
		Preload("Shop").
		Preload("Books").
		Preload("Books.Book").
		Preload("Books.Book.Genres").
		First(&order).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) ReadByIdAndUserId(ctx context.Context, id, userId string) (*entities.Order, error) {
	var order entities.Order

	if err := r.db.
		WithContext(ctx).
		Model(&entities.Order{}).
		Where("orders.id = ? and orders.user_id = ?", id, userId).
		Preload("Shop").
		Preload("Books").
		Preload("Books.Book").
		Preload("Books.Book.Genres").
		First(&order).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *entities.Order) error {
	res := r.db.
		WithContext(ctx).
		Model(&entities.Order{}).
		Where("id = ?", order.Id).
		Updates(order)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}

		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrOrderNotFound
	}

	return nil
}

func (r *orderRepository) UpdateAndPay(ctx context.Context, order *entities.Order) error {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	res := tx.
		Model(&entities.Order{}).
		Where("id = ?", order.Id).
		Updates(order).
		Scan(&order)

	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	if res.RowsAffected == 0 {
		tx.Rollback()
		return ErrOrderNotFound
	}

	res = tx.
		Model(&entities.User{}).
		Where("id = ?", order.UserId).
		Update("points", gorm.Expr("points + ?", order.Points))

	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&entities.Order{Id: id})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return ErrOrderNotFound
	}

	return nil
}
