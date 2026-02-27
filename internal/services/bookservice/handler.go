package bookservice

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type BookService interface {
	CreateBook(ctx context.Context, book *entities.Book) (*entities.Book, error)
	ReadBooks(ctx context.Context, page, limit int) ([]entities.Book, error)
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
// @Success 200 {object} dto.BookResponse
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
		image, mime, err = fromStringToBytes(*request.Image)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid image"})
		}
	}

	book := &entities.Book{
		Title:     request.Title,
		Author:    request.Author,
		Year:      request.Year,
		Cost:      request.Cost,
		Publisher: request.Publisher,
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

	return c.JSON(http.StatusOK, dto.BookResponse{
		Id:          book.Id,
		Title:       book.Title,
		Author:      book.Author,
		Year:        book.Year,
		Cost:        book.Cost,
		Discount:    validate(book.Discount),
		Publisher:   book.Publisher,
		Description: validate(book.Description),
		Amount:      book.Amount,
		Image:       fromBytesToString(book.ImageData, book.ImageMime),
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	book, err := h.service.ReedBookById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "book not found"})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.BookResponse{
		Id:          book.Id,
		Title:       book.Title,
		Author:      book.Author,
		Year:        book.Year,
		Cost:        book.Cost,
		Discount:    validate(book.Discount),
		Publisher:   book.Publisher,
		Description: validate(book.Description),
		Amount:      book.Amount,
		Image:       fromBytesToString(book.ImageData, book.ImageMime),
	})
}

// ReadBooks
// @Summary Получить книги
// @Tags books
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Success 200 {object} dto.BookResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books [get]
func (h *BookHandler) ReadBooks(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page := 1
	limit := 10

	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid page"})
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid limit"})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := h.service.ReadBooks(ctx, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	books := make([]dto.BookResponse, len(response))
	for _, book := range response {
		books = append(books, dto.BookResponse{
			Id:          book.Id,
			Title:       book.Title,
			Author:      book.Author,
			Year:        book.Year,
			Cost:        book.Cost,
			Discount:    validate(book.Discount),
			Publisher:   book.Publisher,
			Description: validate(book.Description),
			Amount:      book.Amount,
			Image:       fromBytesToString(book.ImageData, book.ImageMime),
		})
	}

	return c.JSON(http.StatusOK, books)
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

	ctx, cancel := context.WithTimeout(context.Background(), 500000*time.Second)
	defer cancel()

	var image []byte
	var mime string
	if request.Image != nil {
		var err error
		image, mime, err = fromStringToBytes(*request.Image)
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

	book, err := h.service.UpdateBook(ctx, book)

	if err != nil {
		if errors.Is(err, ErrBookNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "book not found"})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.BookResponse{
		Id:          book.Id,
		Title:       book.Title,
		Author:      book.Author,
		Year:        book.Year,
		Cost:        book.Cost,
		Discount:    validate(book.Discount),
		Publisher:   book.Publisher,
		Description: validate(book.Description),
		Amount:      book.Amount,
		Image:       fromBytesToString(book.ImageData, book.ImageMime),
	})
}

// DeleteBook
// @Summary Удалить книгу
// @Tags books
// @Security BearerAuth
// @Param id path string true "ID книги"
// @Success 204
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

func validate[T any](t *T) T {
	if t != nil {
		return *t
	}

	var zero T
	return zero
}

func fromBytesToString(b []byte, mime string) string {
	if b == nil {
		return ""
	}
	if mime == "" {
		mime = "image/png"
	}
	return "data:" + mime + ";base64," +
		base64.StdEncoding.EncodeToString(b)
}

func fromStringToBytes(str string) ([]byte, string, error) {
	if str == "" {
		return nil, "", nil
	}

	const prefix = "data:"
	if !strings.HasPrefix(str, prefix) {
		b, err := base64.StdEncoding.DecodeString(str)
		return b, "", err
	}

	parts := strings.SplitN(str, ",", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("invalid data url")
	}

	meta := parts[0]
	data := parts[1]

	mime := ""
	if strings.Contains(meta, ";") {
		mime = strings.TrimPrefix(strings.Split(meta, ";")[0], "data:")
	}

	b, err := base64.StdEncoding.DecodeString(data)
	return b, mime, err
}
