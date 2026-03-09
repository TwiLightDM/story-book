package cartservice

import (
	"context"
	"errors"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"
	"story-book/package/services/helperservice"

	"github.com/labstack/echo/v4"
)

type CartService interface {
	CreateCart(ctx context.Context, cart *entities.Cart) error
	ReadCarts(ctx context.Context, userId string, limit, offset int) ([]entities.Book, error)
	UpdateCart(ctx context.Context, cart *entities.Cart) error
	DeleteCart(ctx context.Context, cart, bookId string) error
}

type CartHandler struct {
	service CartService
}

func NewCartHandler(service CartService) *CartHandler {
	return &CartHandler{service: service}
}

// CreateCart
// @Summary Добавить книгу в корзину
// @Tags carts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CartRequest true "Данные для корзины"
// @Success 201 "Created"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /carts [post]
func (h *CartHandler) CreateCart(c echo.Context) error {
	userId := c.Get("id").(string)

	var request dto.CartRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	cart := &entities.Cart{
		Amount: request.Amount,
		BookId: request.BookId,
		UserId: userId,
	}

	err := h.service.CreateCart(context.Background(), cart)
	if err != nil {
		switch {
		case errors.Is(err, ErrCartNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})

		case errors.Is(err, ErrBookInCartAlreadyExists):
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})

		case errors.Is(err, ErrNotEnoughBooks):
			return c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})

		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
	}

	return c.NoContent(http.StatusCreated)
}

// ReadCarts
// @Summary Получить список книг в корзине пользователя
// @Tags carts
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Param offset query int false "Количество пропускаемых записей (по умолчанию 0)"
// @Success 200 {object} dto.BookListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /carts [get]
func (h *CartHandler) ReadCarts(c echo.Context) error {
	limit, offset, err := helperservice.GetLimitAndOffset(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	}

	userId := c.Get("id").(string)

	response, err := h.service.ReadCarts(context.Background(), userId, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	books := make([]dto.BookResponse, 0, len(response))
	for _, book := range response {
		books = append(books, dto.BookResponse{
			Id:          book.Id,
			Title:       book.Title,
			Author:      book.Author,
			Year:        book.Year,
			Cost:        book.Cost,
			Discount:    helperservice.Validate(book.Discount),
			Publisher:   book.Publisher,
			Description: helperservice.Validate(book.Description),
			Amount:      book.Amount,
			Image:       helperservice.FromBytesToString(book.ImageData, book.ImageMime),
			Rating:      book.Rating,
		})
	}

	return c.JSON(http.StatusOK, books)
}

// UpdateCart
// @Summary Обновить количество книг в корзине
// @Tags carts
// @Security BearerAuth
// @Param book_id path string true "ID книги"
// @Accept json
// @Produce json
// @Param request body dto.CartRequest true "Данные магазина"
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /carts/{book_id} [patch]
func (h *CartHandler) UpdateCart(c echo.Context) error {
	userId := c.Get("id").(string)

	var request dto.CartRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	bookId := c.Param("book_id")
	if bookId == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	cart := &entities.Cart{
		Amount: request.Amount,
		BookId: bookId,
		UserId: userId,
	}

	err := h.service.UpdateCart(context.Background(), cart)

	if err != nil {
		switch {
		case errors.Is(err, ErrCartNotFound):
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})

		case errors.Is(err, ErrNotEnoughBooks):
			return c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})

		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
	}

	return c.NoContent(http.StatusOK)
}

// DeleteCart
// @Summary Удалить книгу из корзины
// @Tags carts
// @Security BearerAuth
// @Accept json
// @Param book_id path string true "ID книги"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /carts/{book_id} [delete]
func (h *CartHandler) DeleteCart(c echo.Context) error {
	userId := c.Get("id").(string)

	bookId := c.Param("book_id")
	if bookId == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	err := h.service.DeleteCart(context.Background(), userId, bookId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
