package userservice

import (
	"context"
	"errors"
	"story-book/internal/entities"
	"story-book/package/services/encryptservice"
	"story-book/package/services/jwtservice"
	"story-book/package/services/validateservice"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	ReadByEmail(ctx context.Context, email string) (*entities.User, error)
	ReadById(ctx context.Context, id string) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) (*entities.User, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo     UserRepository
	jwt      jwtservice.JWTService
	encrypt  encryptservice.EncryptService
	validate validateservice.ValidationService
}

func NewUserService(repo UserRepository, jwt jwtservice.JWTService, encrypt encryptservice.EncryptService, validate validateservice.ValidationService) UserService {
	return &userService{repo: repo, jwt: jwt, encrypt: encrypt, validate: validate}
}

func (s *userService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.ReadByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", ErrUserNotFound
	}

	if err = s.encrypt.PasswordComparison(user.Password, password, user.Salt); err != nil {
		return "", "", err
	}

	data := make(map[string]any)
	data["id"] = user.Id
	data["role"] = user.Role

	accessToken, _, err := s.jwt.GenerateAccessJWT(data)
	if err != nil {
		return "", "", err
	}

	refreshToken, _, err := s.jwt.GenerateRefreshJWT(data)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) SignUp(ctx context.Context, user *entities.User) (*entities.User, string, string, error) {
	existing, err := s.repo.ReadByEmail(ctx, user.Email)
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return nil, "", "", err
		}
	}
	if existing != nil {
		return nil, "", "", ErrUserAlreadyExists
	}

	err = s.validate.IsValidEmail(user.Email)
	if err != nil {
		return nil, "", "", err
	}

	err = s.validate.IsStrongPassword(user.Password)
	if err != nil {
		return nil, "", "", err
	}

	hashedPassword, salt, err := s.encrypt.HashPassword(user.Password)
	if err != nil {
		return nil, "", "", err
	}

	user.Id = uuid.NewString()
	user.Password = hashedPassword
	user.Salt = salt
	user.Role = "client"

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	data := make(map[string]any)
	data["id"] = user.Id
	data["role"] = user.Role

	accessToken, _, err := s.jwt.GenerateAccessJWT(data)
	if err != nil {
		return user, "", "", err
	}

	refreshToken, _, err := s.jwt.GenerateRefreshJWT(data)
	if err != nil {
		return user, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *userService) ReedUserById(ctx context.Context, id string) (*entities.User, error) {
	user, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	var err error

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) UpdatePassword(ctx context.Context, user *entities.User) error {
	err := s.validate.IsStrongPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password, user.Salt, err = s.encrypt.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = s.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) RefreshTokens(id, role string) (string, string, error) {
	data := make(map[string]any)
	data["id"] = id
	data["role"] = role
	accessToken, _, err := s.jwt.GenerateAccessJWT(data)
	if err != nil {
		return "", "", err
	}

	refreshToken, _, err := s.jwt.GenerateRefreshJWT(data)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) ResetPassword(ctx context.Context, id, answer string) error {
	user, err := s.repo.ReadById(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if answer == user.Answer {
		return nil
	}

	return ErrWrongAnswer
}
