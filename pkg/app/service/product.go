package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"restaurant-assistant/pkg/base"
	"restaurant-assistant/pkg/entity"
	"restaurant-assistant/pkg/storage/database/repository"
)

type ProductService struct {
	log *zap.Logger

	pr *repository.ProductRepository
	ur *repository.UserRepository
}

func NewProductService(
	log *zap.Logger,
	pr *repository.ProductRepository,
	ur *repository.UserRepository,
) (*ProductService, error) {
	if log == nil {
		return nil, fmt.Errorf("uninitialized logger")
	}

	if pr == nil {
		return nil, fmt.Errorf("uninitialized product repository")
	}

	if ur == nil {
		return nil, fmt.Errorf("uninitialized user repository")
	}

	return &ProductService{
		log: log,
		pr:  pr,
		ur:  ur,
	}, nil
}

func (ps *ProductService) CreateProduct(ctx context.Context, product entity.Product) (*entity.Product, error) {
	if err := product.Validate(); err != nil {
		return nil, err
	}

	if p, _ := ps.pr.Context(ctx).FindByName(product.Name); p != nil {
		return nil, base.NewError(
			base.ErrInvalidParameters,
			"product with specified name already exists",
		)
	}

	newProduct, err := ps.pr.Context(ctx).Create(product)
	if err != nil {
		return nil, err
	}

	return newProduct, nil
}
