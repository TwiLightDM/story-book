package favouriteservice

import (
	"context"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"
	"story-book/package/services/helperservice"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type FavouriteService interface {
	CreateFavourite(ctx context.Context, favourite *entities.Favourite) error
	ReadFavourites(ctx context.Context, userId string, limit, offset int) ([]entities.Book, error)
	DeleteFavourite(ctx context.Context, userId, bookId string) error
}

type FavouriteHandler struct {
	service FavouriteService
}

func NewFavouriteHandler(service FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{service: service}
}

// CreateFavourite
// @Summary Добавить книгу в избранное
// @Tags favourites
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.FavouriteRequest true "Данные избранного"
// @Success 201 "Created"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /favourites [post]
func (h *FavouriteHandler) CreateFavourite(c echo.Context) error {
	userId := c.Get("id").(string)

	var request dto.FavouriteRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	favourite := &entities.Favourite{
		BookId: request.BookId,
		UserId: userId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.service.CreateFavourite(ctx, favourite)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

// ReadFavourites
// @Summary Получить список избранных книг пользователя
// @Tags favourites
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Param offset query int false "Количество пропускаемых записей (по умолчанию 0)"
// @Success 200 {object} dto.BookListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /favourites [get]
func (h *FavouriteHandler) ReadFavourites(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10
	offset := 0

	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid limit"})
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid offset"})
		}
	}

	userId := c.Get("id").(string)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := h.service.ReadFavourites(ctx, userId, limit, offset)
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
		})
	}

	return c.JSON(http.StatusOK, books)
}

// DeleteFavourite
// @Summary Удалить книгу из избранного
// @Tags favourites
// @Security BearerAuth
// @Param book_id path string true "ID книги"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /favourites/{book_id} [delete]
func (h *FavouriteHandler) DeleteFavourite(c echo.Context) error {
	userId := c.Get("id").(string)

	bookId := c.Param("book_id")
	if bookId == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.service.DeleteFavourite(ctx, userId, bookId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
