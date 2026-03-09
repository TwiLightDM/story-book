package shopservice

import (
	"context"
	"errors"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"

	"github.com/labstack/echo/v4"
)

type ShopService interface {
	CreateShop(ctx context.Context, shop *entities.Shop) (*entities.Shop, error)
	ReadShops(ctx context.Context) ([]entities.Shop, error)
	ReedShopById(ctx context.Context, id string) (*entities.Shop, error)
	UpdateShop(ctx context.Context, shop *entities.Shop) (*entities.Shop, error)
	DeleteShop(ctx context.Context, id string) error
}

type ShopHandler struct {
	service ShopService
}

func NewShopHandler(service ShopService) *ShopHandler {
	return &ShopHandler{service: service}
}

// CreateShop
// @Summary Создать магазин
// @Tags shops
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ShopRequest true "Данные книги"
// @Success 201 {object} dto.ShopResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shops [post]
func (h *ShopHandler) CreateShop(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	var request dto.ShopRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	shop := &entities.Shop{
		Name:    request.Name,
		Address: request.Address,
	}

	shop, err := h.service.CreateShop(context.Background(), shop)
	if err != nil {
		if errors.Is(err, ErrShopAlreadyExists) {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.ShopResponse{
		Id:      shop.Id,
		Name:    shop.Name,
		Address: shop.Address,
	})
}

// ReadShop
// @Summary Получить магазин по ID
// @Tags shops
// @Param id path string true "ID магазина"
// @Produce json
// @Success 200 {object} dto.ShopResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shops/{id} [get]
func (h *ShopHandler) ReadShop(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	shop, err := h.service.ReedShopById(context.Background(), id)
	if err != nil {
		if errors.Is(err, ErrShopNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ShopResponse{
		Id:      shop.Id,
		Name:    shop.Name,
		Address: shop.Address,
	})
}

// ReadShops
// @Summary Получить магазины
// @Tags shops
// @Produce json
// @Success 200 {object} dto.ShopListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shops [get]
func (h *ShopHandler) ReadShops(c echo.Context) error {
	response, err := h.service.ReadShops(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	shops := make([]dto.ShopResponse, 0, len(response))

	for _, shop := range response {
		shops = append(shops, dto.ShopResponse{
			Id:      shop.Id,
			Name:    shop.Name,
			Address: shop.Address,
		})
	}

	return c.JSON(http.StatusOK, dto.ShopListResponse{
		Shops: shops,
	})
}

// UpdateShop
// @Summary Обновить магазин
// @Tags shops
// @Security BearerAuth
// @Param id path string true "ID магазина"
// @Accept json
// @Produce json
// @Param request body dto.ShopRequest true "Данные магазина"
// @Success 200 {object} dto.ShopResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shops/{id} [put]
func (h *ShopHandler) UpdateShop(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	var request dto.ShopRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	shop := &entities.Shop{
		Id:      id,
		Name:    request.Name,
		Address: request.Address,
	}

	shop, err := h.service.UpdateShop(context.Background(), shop)

	if err != nil {
		if errors.Is(err, ErrShopNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ShopResponse{
		Id:      shop.Id,
		Name:    shop.Name,
		Address: shop.Address,
	})
}

// DeleteShop
// @Summary Удалить магазин
// @Tags shops
// @Security BearerAuth
// @Param id path string true "ID магазина"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shops/{id} [delete]
func (h *ShopHandler) DeleteShop(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	err := h.service.DeleteShop(context.Background(), id)
	if err != nil {
		if errors.Is(err, ErrShopNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
