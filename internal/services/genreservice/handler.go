package genreservice

import (
	"context"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"
	"time"

	"github.com/labstack/echo/v4"
)

type GenreService interface {
	CreateGenre(ctx context.Context, genre *entities.Genre) error
	DeleteGenre(ctx context.Context, genre, bookId string) error
}

type GenreHandler struct {
	service GenreService
}

func NewGenreHandler(service GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

// CreateGenre
// @Summary Добавить жанр к книге
// @Tags genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.GenreRequest true "Данные жанра"
// @Success 201 "Created"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /genres [post]
func (h *GenreHandler) CreateGenre(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	var request dto.GenreRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	genre := &entities.Genre{
		Genre:  request.Genre,
		BookId: request.BookId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.service.CreateGenre(ctx, genre)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

// DeleteGenre
// @Summary Удалить жанр книги
// @Tags genres
// @Security BearerAuth
// @Accept json
// @Param book_id path string true "ID жанра"
// @Param request body dto.GenreRequest true "Данные жанра"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /genres/{book_id} [delete]
func (h *GenreHandler) DeleteGenre(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	bookId := c.Param("book_id")
	if bookId == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	var request dto.GenreRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.service.DeleteGenre(ctx, request.Genre, bookId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
