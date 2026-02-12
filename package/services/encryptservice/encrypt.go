package encryptservice

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type EncryptService interface {
	HashPassword(password string) (string, string, error)
	PasswordComparison(hashedPassword, password, salt string) error
}

type encryptService struct {
	SaltLength int
}

func NewEncryptionService(salt int) EncryptService {
	return &encryptService{
		SaltLength: salt,
	}
}

func (e encryptService) HashPassword(password string) (string, string, error) {
	salt, err := e.saltGeneration()
	if err != nil {
		return "", "", ErrGeneratingSalt
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(salt+password), bcrypt.MinCost)
	if err != nil {
		return "", "", errors.Join(ErrHashingPassword, err)
	}
	return string(hashedPassword), salt, nil
}

func (e encryptService) PasswordComparison(hashedPassword, password, salt string) error {
	saltPassword := salt + password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltPassword))
	if err != nil {
		return errors.Join(ErrInvalidPassword, err)
	}
	return nil
}

func (e encryptService) saltGeneration() (string, error) {
	bytes := make([]byte, e.SaltLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.Join(ErrGeneratingSalt, err)
	}
	return base64.StdEncoding.EncodeToString(bytes)[:e.SaltLength], nil
}
