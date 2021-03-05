package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"restaurant-assistant/pkg/base"
	"restaurant-assistant/pkg/entity"
	"restaurant-assistant/pkg/storage/database/repository"
	"restaurant-assistant/pkg/utils"
)

type UserService struct {
	log *zap.Logger

	// TODO: replace to data provider
	ur *repository.UserRepository
}

func NewUserService(
	log *zap.Logger,
	ur *repository.UserRepository,
) (*UserService, error) {
	if log == nil {
		return nil, fmt.Errorf("uninitialized logger")
	}

	if ur == nil {
		return nil, fmt.Errorf("uninitialized user repository")
	}

	return &UserService{
		log: log,
		ur:  ur,
	}, nil
}

func (us *UserService) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// Validate user fields
	err := user.Validate()
	if err != nil {
		return nil, err
	}

	if u, _ := us.ur.Context(ctx).FindByName(user.Name); u != nil {
		return nil, base.NewError(
			base.ErrInvalidParameters,
			"user with specified name already exists",
		)
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	createdUser, err := us.ur.Context(ctx).Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// Validate user fields
	err := user.Validate()
	if err != nil {
		return nil, err
	}

	if u, _ := us.ur.Context(ctx).FindByName(user.Name); u != nil {
		return nil, base.NewError(
			base.ErrInvalidParameters,
			"user with specified name already exists",
		)
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	createdUser, err := us.ur.Context(ctx).Create(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
