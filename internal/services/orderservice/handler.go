package orderservice

import (
	"context"
	"errors"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"
	"story-book/package/services/helperservice"
	"time"

	"github.com/labstack/echo/v4"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *entities.Order, code string) error
	ReadOrders(ctx context.Context, limit, offset int, userId string, startDate, endDate time.Time) ([]entities.Order, error)
	ReedOrderById(ctx context.Context, id, userId string) (*entities.Order, error)
	UpdateOrder(ctx context.Context, order *entities.Order) error
	PayOrder(ctx context.Context, order *entities.Order) error
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder
// @Summary Создать заказ
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.OrderRequest true "Данные заказа"
// @Success 201
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var request dto.OrderRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	userId := c.Get("id").(string)

	order := &entities.Order{
		Address: helperservice.Validate(&request.Address),
		ShopId:  helperservice.Validate(&request.ShopId),
		UserId:  userId,
	}
	err := h.service.CreateOrder(context.Background(), order, helperservice.Validate(request.Code))
	if err != nil {
		if errors.Is(err, ErrCodeNotFound) || errors.Is(err, ErrCodeExpired) || errors.Is(err, ErrAmountOfBooks) || errors.Is(err, ErrAmountOfUsage) {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		}

		if errors.Is(err, ErrEmptyCart) || errors.Is(err, ErrDeliveryDestinationRequired) || errors.Is(err, ErrDeliveryDestinationConflict) {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

// ReadOrder
// @Summary Получить закад по ID
// @Tags orders
// @Security BearerAuth
// @Param id path string true "ID заказа"
// @Produce json
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) ReadOrder(c echo.Context) error {
	userId := c.Get("id").(string)
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	order, err := h.service.ReedOrderById(context.Background(), id, userId)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	orderBooks := make([]dto.OrderBooksResponse, 0, len(order.Books))
	for _, ob := range order.Books {
		orderBooks = append(orderBooks, dto.OrderBooksResponse{
			Amount:      ob.Amount,
			PriceForOne: ob.PriceForOne,
			Book: dto.BookResponse{
				Id:        ob.Book.Id,
				Title:     ob.Book.Title,
				Author:    ob.Book.Author,
				Year:      ob.Book.Year,
				Publisher: ob.Book.Publisher,
				Image:     helperservice.FromBytesToString(ob.Book.ImageData, ob.Book.ImageMime),
				Rating:    ob.Book.Rating,
			},
		})
	}

	response := dto.OrderResponse{
		Id:           order.Id,
		Status:       order.Status,
		DeliveryType: order.DeliveryType,
		Address:      helperservice.Validate(order.Address),
		Cost:         order.Cost,
		Points:       order.Points,
		CreatedAt:    order.CreatedAt.Format("2006-01-02"),
		Shop: dto.ShopResponse{
			Id:      order.Shop.Id,
			Name:    order.Shop.Name,
			Address: order.Shop.Address,
		},
		OrderBooks: orderBooks,
	}

	return c.JSON(http.StatusOK, response)
}

// ReadOrders
// @Summary Получить заказы
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param start_date query string false "Фильтр: начальная дата (формат YYYY-MM-DD)"
// @Param end_date   query string false "Фильтр: конечная дата (формат YYYY-MM-DD)"
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Param offset query int false "Количество пропускаемых записей (по умолчанию 0)"
// @Success 200 {object} dto.OrderListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) ReadOrders(c echo.Context) error {
	userId := c.Get("id").(string)

	limit, offset, err := helperservice.GetLimitAndOffset(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	}

	startDate, endDate, err := helperservice.GetStartAndEndDate(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	}

	ordersFromDB, err := h.service.ReadOrders(context.Background(), limit, offset, userId, startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	ordersDTO := make([]dto.OrderResponse, 0, len(ordersFromDB))

	for _, order := range ordersFromDB {
		orderBooks := make([]dto.OrderBooksResponse, 0, len(order.Books))
		for _, ob := range order.Books {
			orderBooks = append(orderBooks, dto.OrderBooksResponse{
				Amount:      ob.Amount,
				PriceForOne: ob.PriceForOne,
				Book: dto.BookResponse{
					Id:        ob.Book.Id,
					Title:     ob.Book.Title,
					Author:    ob.Book.Author,
					Year:      ob.Book.Year,
					Publisher: ob.Book.Publisher,
					Image:     helperservice.FromBytesToString(ob.Book.ImageData, ob.Book.ImageMime),
					Rating:    ob.Book.Rating,
				},
			})
		}

		ordersDTO = append(ordersDTO, dto.OrderResponse{
			Id:           order.Id,
			Status:       order.Status,
			DeliveryType: order.DeliveryType,
			Address:      helperservice.Validate(order.Address),
			Cost:         order.Cost,
			Points:       order.Points,
			CreatedAt:    order.CreatedAt.Format("2006-01-02"),
			Shop: dto.ShopResponse{
				Id:      order.Shop.Id,
				Name:    order.Shop.Name,
				Address: order.Shop.Address,
			},
			OrderBooks: orderBooks,
		})
	}

	return c.JSON(http.StatusOK, ordersDTO)
}

// UpdateOrder
// @Summary Поменять статус заказа
// @Tags orders
// @Security BearerAuth
// @Param id path string true "ID заказа"
// @Accept json
// @Produce json
// @Param request body dto.StatusRequest true "Данные статуса"
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/{id} [patch]
func (h *OrderHandler) UpdateOrder(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	var request dto.StatusRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	order := &entities.Order{
		Id:     id,
		Status: request.Status,
	}

	err := h.service.UpdateOrder(context.Background(), order)

	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// PayOrder
// @Summary Оплатить заказ (поменять статус на "paid")
// @Tags orders
// @Security BearerAuth
// @Param id path string true "ID заказа"
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /orders/pay/{id} [patch]
func (h *OrderHandler) PayOrder(c echo.Context) error {
	var request dto.StatusRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	order := &entities.Order{
		Id:     id,
		Status: request.Status,
	}

	err := h.service.PayOrder(context.Background(), order)

	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}

		if errors.Is(err, ErrFailedToPay) || errors.Is(err, ErrAlreadyPaid) {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
