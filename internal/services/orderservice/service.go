package orderservice

import (
	"context"
	"math/rand"
	"story-book/internal/entities"
	"time"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order, code string) error
	ReadAll(ctx context.Context, limit, offset int, userId string, startDate, endDate time.Time) ([]entities.Order, error)
	ReadById(ctx context.Context, id string) (*entities.Order, error)
	ReadByIdAndUserId(ctx context.Context, id, userId string) (*entities.Order, error)
	Update(ctx context.Context, order *entities.Order) error
	UpdateAndPay(ctx context.Context, order *entities.Order) error
	Delete(ctx context.Context, id string) error
}

type orderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(ctx context.Context, order *entities.Order, code string) error {
	if order.ShopId == nil && order.Address == nil {
		return ErrDeliveryDestinationRequired
	}
	if order.ShopId != nil && order.Address != nil {
		return ErrDeliveryDestinationConflict
	}

	order.Id = uuid.NewString()
	order.Status = "created"

	if order.ShopId != nil {
		order.DeliveryType = "pickup"
	}

	if order.Address != nil {
		order.DeliveryType = "delivery"
	}

	return s.repo.Create(ctx, order, code)
}

func (s *orderService) ReedOrderById(ctx context.Context, id, userId string) (*entities.Order, error) {
	order, err := s.repo.ReadByIdAndUserId(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) ReadOrders(ctx context.Context, limit, offset int, userId string, startDate, endDate time.Time) ([]entities.Order, error) {
	orders, err := s.repo.ReadAll(ctx, limit, offset, userId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, order *entities.Order) error {
	err := s.repo.Update(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *orderService) PayOrder(ctx context.Context, order *entities.Order) error {
	o, err := s.repo.ReadById(ctx, order.Id)
	if err != nil {
		return err
	}

	if o.Status != "created" {
		return ErrAlreadyPaid
	}

	if rand.Intn(100) < 70 {
		return ErrFailedToPay
	}

	order.Status = "paid"

	err = s.repo.UpdateAndPay(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *orderService) DeleteOrder(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
