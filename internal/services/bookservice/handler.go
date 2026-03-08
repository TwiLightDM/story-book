package bookservice

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

type BookService interface {
	CreateBook(ctx context.Context, book *entities.Book) (*entities.Book, error)
	ReadBooks(ctx context.Context, limit, offset int) ([]entities.Book, error)
	ReedBookById(ctx context.Context, id string) (*entities.Book, error)
	UpdateBook(ctx context.Context, book *entities.Book) (*entities.Book, error)
	DeleteBook(ctx context.Context, id string) error
}

type BookHandler struct {
	service BookService
}

func NewBookHandler(service BookService) *BookHandler {
	return &BookHandler{service: service}
}

// CreateBook
// @Summary Создать книгу
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.BookRequest true "Данные книги"
// @Success 201 {object} dto.BookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books [post]
func (h *BookHandler) CreateBook(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	var request dto.BookRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	var image []byte
	var mime string
	if request.Image != nil {
		var err error
		image, mime, err = helperservice.FromStringToBytes(*request.Image)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid image"})
		}
	}

	if request.Amount < 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid amount"})
	}

	book := &entities.Book{
		Title:     request.Title,
		Author:    request.Author,
		Year:      request.Year,
		Cost:      request.Cost,
		Publisher: request.Publisher,
		Amount:    request.Amount,
		ImageData: image,
		ImageMime: mime,
	}

	if request.Discount != nil {
		book.Discount = request.Discount
	}

	if request.Description != nil {
		book.Description = request.Description
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	book, err := h.service.CreateBook(ctx, book)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.BookResponse{
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

// ReadBook
// @Summary Получить книгу по ID
// @Tags books
// @Param id path string true "ID книги"
// @Produce json
// @Success 200 {object} dto.BookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) ReadBook(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	book, err := h.service.ReedBookById(context.Background(), id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	genres := make([]dto.GenreResponse, 0, len(book.Genres))
	for _, genre := range book.Genres {
		genres = append(genres, dto.GenreResponse{
			Genre: genre.Genre,
		})
	}

	return c.JSON(http.StatusOK, dto.BookResponse{
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
		Genres:      genres,
	})
}

// ReadBooks
// @Summary Получить книги
// @Tags books
// @Produce json
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Param offset query int false "Количество пропускаемых записей (по умолчанию 0)"
// @Success 200 {object} dto.BookListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books [get]
func (h *BookHandler) ReadBooks(c echo.Context) error {
	limit, offset, err := helperservice.GetLimitAndOffset(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	}

	response, err := h.service.ReadBooks(context.Background(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	books := make([]dto.BookResponse, 0, len(response))

	for _, book := range response {
		genres := make([]dto.GenreResponse, 0, len(book.Genres))
		for _, genre := range book.Genres {
			genres = append(genres, dto.GenreResponse{
				Genre: genre.Genre,
			})
		}

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
			Genres:      genres,
		})
	}

	return c.JSON(http.StatusOK, dto.BookListResponse{
		Books: books,
	})
}

// UpdateBook
// @Summary Обновить книгу
// @Tags books
// @Security BearerAuth
// @Param id path string true "ID книги"
// @Accept json
// @Produce json
// @Param request body dto.BookRequest true "Данные книги"
// @Success 200 {object} dto.BookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	var request dto.BookRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	var image []byte
	var mime string
	if request.Image != nil {
		var err error
		image, mime, err = helperservice.FromStringToBytes(*request.Image)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid image"})
		}
	}

	book := &entities.Book{
		Id:        id,
		Title:     request.Title,
		Author:    request.Author,
		Year:      request.Year,
		Cost:      request.Cost,
		Publisher: request.Publisher,
		Amount:    request.Amount,
		ImageData: image,
		ImageMime: mime,
	}

	if request.Discount != nil {
		book.Discount = request.Discount
	}

	if request.Description != nil {
		book.Description = request.Description
	}

	book, err := h.service.UpdateBook(context.Background(), book)

	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	genres := make([]dto.GenreResponse, 0, len(book.Genres))
	for _, genre := range book.Genres {
		genres = append(genres, dto.GenreResponse{
			Genre: genre.Genre,
		})
	}

	return c.JSON(http.StatusOK, dto.BookResponse{
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
		Genres:      genres,
	})
}

// DeleteBook
// @Summary Удалить книгу
// @Tags books
// @Security BearerAuth
// @Param id path string true "ID книги"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.service.DeleteBook(ctx, id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "book not found"})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
