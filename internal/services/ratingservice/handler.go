package ratingservice

import (
	"context"
	"errors"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"

	"github.com/labstack/echo/v4"
)

type RatingService interface {
	CreateRating(ctx context.Context, rating *entities.Rating) error
	DeleteRating(ctx context.Context, rating, bookId string) error
}

type RatingHandler struct {
	service RatingService
}

func NewRatingHandler(service RatingService) *RatingHandler {
	return &RatingHandler{service: service}
}

// CreateRating
// @Summary Оценить книгу
// @Tags ratings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.RatingRequest true "Данные жанра"
// @Success 201 "Created"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ratings [post]
func (h *RatingHandler) CreateRating(c echo.Context) error {
	userId := c.Get("id").(string)

	var request dto.RatingRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	rating := &entities.Rating{
		Stars:  request.Stars,
		BookId: request.BookId,
		UserId: userId,
	}

	err := h.service.CreateRating(context.Background(), rating)
	if err != nil {
		if errors.Is(err, ErrUnsupportedRating) {
			return c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

// DeleteRating
// @Summary Удалить оценку
// @Tags ratings
// @Security BearerAuth
// @Accept json
// @Param book_id path string true "ID книги"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ratings/{book_id} [delete]
func (h *RatingHandler) DeleteRating(c echo.Context) error {
	userId := c.Get("id").(string)

	bookId := c.Param("book_id")
	if bookId == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	err := h.service.DeleteRating(context.Background(), userId, bookId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
