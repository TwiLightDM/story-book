package userservice

import (
	"context"
	"errors"
	"story-book/internal/entities"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) ReadByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.
		WithContext(ctx).
		Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ReadById(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	err := r.db.
		WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", user.Id).
		Updates(user).
		Scan(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entities.User{Id: id}).Error
}
