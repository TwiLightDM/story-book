package cardservice

import (
	"context"
	"errors"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"

	"github.com/labstack/echo/v4"
)

type CardService interface {
	CreateCard(ctx context.Context, card *entities.Card) (*entities.Card, error)
	ReadCards(ctx context.Context) ([]entities.Card, error)
	DeleteCard(ctx context.Context, numberOfCard string) error
}

type CardHandler struct {
	service CardService
}

func NewCardHandler(service CardService) *CardHandler {
	return &CardHandler{service: service}
}

// CreateCard
// @Summary Создать карту
// @Tags cards
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CardRequest true "Данные карты"
// @Success 201 {object} dto.CardResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cards [post]
func (h *CardHandler) CreateCard(c echo.Context) error {
	userId := c.Get("id").(string)
	var request dto.CardRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	card := &entities.Card{
		NumberOfCard:   request.NumberOfCard,
		ExpirationDate: request.ExpirationDate,
		Cvv:            request.Cvv,
		UserId:         userId,
	}

	card, err := h.service.CreateCard(context.Background(), card)
	if err != nil {
		if errors.Is(err, ErrCardAlreadyExists) {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.CardResponse{
		NumberOfCard: card.NumberOfCard,
	})
}

// ReadCards
// @Summary Получить карты
// @Tags cards
// @Produce json
// @Success 200 {object} dto.CardListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cards [get]
func (h *CardHandler) ReadCards(c echo.Context) error {
	response, err := h.service.ReadCards(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	cards := make([]dto.CardResponse, 0, len(response))

	for _, card := range response {
		cards = append(cards, dto.CardResponse{
			NumberOfCard: card.NumberOfCard,
		})
	}

	return c.JSON(http.StatusOK, dto.CardListResponse{
		Cards: cards,
	})
}

// DeleteCard
// @Summary Удалить карту
// @Tags cards
// @Security BearerAuth
// @Param number_of_card query int true "Номер карты для удаления"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cards [delete]
func (h *CardHandler) DeleteCard(c echo.Context) error {
	numberOfCard := c.QueryParams().Get("number_of_card")
	if numberOfCard == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	err := h.service.DeleteCard(context.Background(), numberOfCard)
	if err != nil {
		if errors.Is(err, ErrCardNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
