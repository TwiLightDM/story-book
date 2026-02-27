package userservice

import (
	"context"
	"errors"
	"net/http"
	"story-book/internal/dto"
	"story-book/internal/entities"
	"time"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	SignUp(ctx context.Context, user *entities.User) (*entities.User, string, string, error)
	ReedUserById(ctx context.Context, id string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	UpdatePassword(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id string) error
	RefreshTokens(id, role string) (string, string, error)
	ResetPassword(ctx context.Context, id, answer string) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Login
// @Summary Вход пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserRequest true "Данные входа"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accessToken, refreshToken, err := h.service.Login(ctx, request.Email, request.Password)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// SignUp
// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserRequest true "Данные регистрации"
// @Success 201 {object} dto.SignUpResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/signup [post]
func (h *UserHandler) SignUp(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, accessToken, refreshToken, err := h.service.SignUp(ctx, &entities.User{
		Name:     request.Name,
		Surname:  request.Surname,
		Email:    request.Email,
		Phone:    request.Phone,
		Password: request.Password,
		Question: request.Question,
		Answer:   request.Answer,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.SignUpResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			Id:       user.Id,
			Name:     user.Name,
			Surname:  user.Surname,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			Question: user.Question,
			Points:   user.Points,
		},
	})
}

// Refresh
// @Summary Обновление токенов
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *UserHandler) Refresh(c echo.Context) error {
	id := c.Get("id").(string)
	role := c.Get("role").(string)
	access, refresh, err := h.service.RefreshTokens(id, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// ResetPassword
// @Summary Сброс пароля
// @Tags auth
// @Security BearerAuth
// @Param answer query string true "Ответ на секретный вопрос"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/reset-password [post]
func (h *UserHandler) ResetPassword(c echo.Context) error {
	id := c.Get("id").(string)
	answer := c.QueryParam("answer")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.service.ResetPassword(ctx, id, answer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// ReadSelf
// @Summary Получить текущего пользователя
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) ReadSelf(c echo.Context) error {
	id := c.Get("id").(string)
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.service.ReedUserById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Question: user.Question,
		Points:   user.Points,
	})
}

// ReadUser
// @Summary Получить пользователя по ID
// @Tags users
// @Security BearerAuth
// @Param id path string true "ID пользователя"
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) ReadUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	role := c.Get("role").(string)
	if role == "client" {
		tokenId := c.Get("id").(string)
		if !(tokenId == id) {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.service.ReedUserById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:      user.Id,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		Phone:   user.Phone,
		Role:    user.Role,
		Points:  user.Points,
	})
}

// UpdateUser
// @Summary Обновить профиль
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UserRequest true "Данные пользователя"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Get("id").(string)
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.service.UpdateUser(ctx, &entities.User{
		Id:       id,
		Name:     request.Name,
		Surname:  request.Surname,
		Email:    request.Email,
		Phone:    request.Phone,
		Question: request.Question,
		Answer:   request.Answer,
		Points:   0,
	})
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       user.Id,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Question: user.Question,
		Points:   user.Points,
	})
}

// ChangePassword
// @Summary Изменить пароль
// @Tags users
// @Security BearerAuth
// @Accept json
// @Param request body dto.UserRequest true "Новый пароль"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me/password [patch]
func (h *UserHandler) ChangePassword(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Get("id").(string)
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.service.UpdatePassword(ctx, &entities.User{
		Id:       id,
		Password: request.Password,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteUser
// @Summary Удалить аккаунт
// @Tags users
// @Security BearerAuth
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Get("id").(string)
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.service.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
