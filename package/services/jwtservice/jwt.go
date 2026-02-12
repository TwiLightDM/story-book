package jwtservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTService interface {
	GenerateRefreshJWT(data map[string]any) (string, *time.Time, error)
	GenerateAccessJWT(data map[string]any) (string, *time.Time, error)
	ParseJWT(tokenString string) (map[string]any, error)
}

type jwtService struct {
	Key             string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

func NewJWTService(key string, accessDuration, refreshDuration time.Duration) JWTService {
	return &jwtService{
		Key:             key,
		AccessDuration:  accessDuration,
		RefreshDuration: refreshDuration,
	}
}

func (s *jwtService) generateJWT(data map[string]any, expiresAt *time.Time) (string, *time.Time, error) {
	claims := jwt.MapClaims{
		"exp": expiresAt.Unix(),
	}

	for key, value := range data {
		if key == "id" {
			claims["sub"] = value
		} else {
			claims[key] = value
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.Key))
	if err != nil {
		return "", nil, errors.Join(ErrSignToken, err)
	}

	return signedToken, expiresAt, nil
}

func (s *jwtService) GenerateRefreshJWT(data map[string]any) (string, *time.Time, error) {
	expiresAt := time.Now().Add(s.RefreshDuration)
	delete(data, "exp")

	return s.generateJWT(data, &expiresAt)
}

func (s *jwtService) GenerateAccessJWT(data map[string]any) (string, *time.Time, error) {
	expiresAt := time.Now().Add(s.AccessDuration)
	delete(data, "exp")

	return s.generateJWT(data, &expiresAt)
}

func (s *jwtService) ParseJWT(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", ErrUnexpectedSigningMethod)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidTokenClaims
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, ErrLifetimeIsOver
		}
	}

	data := make(map[string]any)
	for key, value := range claims {
		data[key] = value
	}

	return data, nil
}
