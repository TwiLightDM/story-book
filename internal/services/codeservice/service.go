package codeservice

import (
	"context"
	"story-book/internal/entities"

	"github.com/google/uuid"
)

type CodeRepository interface {
	Create(ctx context.Context, code *entities.Code) error
	ReadAll(ctx context.Context, limit, offset int) ([]entities.Code, error)
	Delete(ctx context.Context, code string) error
}

type codeService struct {
	repo CodeRepository
}

func NewCodeService(repo CodeRepository) CodeService {
	return &codeService{repo: repo}
}

func (s *codeService) CreateCode(ctx context.Context, code *entities.Code) error {
	code.Id = uuid.NewString()

	err := s.repo.Create(ctx, code)
	if err != nil {
		return err
	}

	return nil
}

func (s *codeService) ReadCodes(ctx context.Context, limit, offset int) ([]entities.Code, error) {
	codes, err := s.repo.ReadAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return codes, nil
}

func (s *codeService) DeleteCode(ctx context.Context, code string) error {
	err := s.repo.Delete(ctx, code)
	if err != nil {
		return err
	}

	return nil
}
