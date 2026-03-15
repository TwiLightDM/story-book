package codeservice

import (
	"context"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"
	"story-book/package/services/helperservice"
	"time"

	"github.com/labstack/echo/v4"
)

type CodeService interface {
	CreateCode(ctx context.Context, code *entities.Code) error
	ReadCodes(ctx context.Context, limit, offset int) ([]entities.Code, error)
	DeleteCode(ctx context.Context, code string) error
}

type CodeHandler struct {
	service CodeService
}

func NewCodeHandler(service CodeService) *CodeHandler {
	return &CodeHandler{service: service}
}

// CreateCode
// @Summary Создать промо-код
// @Tags codes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CodeRequest true "Данные промо-кода"
// @Success 201 "Created"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /codes [post]
func (h *CodeHandler) CreateCode(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	var request dto.CodeRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	var expiredAt *time.Time
	if request.ExpiredAt != "" {
		t, err := time.Parse("2006-01-02", request.ExpiredAt)
		if err != nil {
			return err
		}
		expiredAt = &t
	}

	code := &entities.Code{
		Code:          request.Code,
		Percent:       request.Percent,
		AmountOfUsage: &request.AmountOfUsage,
		ExpiredAt:     expiredAt,
	}

	err := h.service.CreateCode(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

// ReadCodes
// @Summary Получить промо-коды
// @Tags codes
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Количество записей на странице (по умолчанию 10)"
// @Param offset query int false "Количество пропускаемых записей (по умолчанию 0)"
// @Success 200 {object} dto.CodeListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /codes [get]
func (h *CodeHandler) ReadCodes(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	limit, offset, err := helperservice.GetLimitAndOffset(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	}

	response, err := h.service.ReadCodes(context.Background(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	codes := make([]dto.CodeResponse, 0, len(response))

	for _, code := range response {
		promo := dto.CodeResponse{
			Code:          code.Code,
			Percent:       code.Percent,
			AmountOfUsage: helperservice.Validate(code.AmountOfUsage),
		}

		if code.ExpiredAt != nil && !code.ExpiredAt.IsZero() {
			promo.ExpiredAt = code.ExpiredAt.Format("2006-01-02")
		}

		codes = append(codes, promo)
	}

	return c.JSON(http.StatusOK, dto.CodeListResponse{
		Codes: codes,
	})
}

// DeleteCode
// @Summary Удалить промо-код
// @Tags codes
// @Security BearerAuth
// @Accept json
// @Param code path string true "Название промо-кода"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /codes/{code} [delete]
func (h *CodeHandler) DeleteCode(c echo.Context) error {
	role := c.Get("role").(string)
	if role == "client" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
	}

	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	err := h.service.DeleteCode(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
