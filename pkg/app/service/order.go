package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"restaurant-assistant/pkg/base"
	"restaurant-assistant/pkg/entity"
	"restaurant-assistant/pkg/storage/database/repository"
)

type OrderService struct {
	log *zap.Logger

	// TODO: replace to data provider
	ur *repository.UserRepository
	or *repository.OrderRepository
	pr *repository.ProductRepository
}

func NewOrderService(
	log *zap.Logger,
	ur *repository.UserRepository,
	or *repository.OrderRepository,
	pr *repository.ProductRepository,
) (*OrderService, error) {
	if log == nil {
		return nil, fmt.Errorf("uninitialized logger")
	}

	if ur == nil {
		return nil, fmt.Errorf("uninitialized user repository")
	}

	if or == nil {
		return nil, fmt.Errorf("uninitialized order repository")
	}

	if pr == nil {
		return nil, fmt.Errorf("uninitialized product repository")
	}

	return &OrderService{
		log: log,
		ur:  ur,
		or:  or,
		pr:  pr,
	}, nil
}

func (os *OrderService) CreateOrder(ctx context.Context, order entity.Order) (*entity.Order, error) {
	if err := order.Validate(); err != nil {
		return nil, err
	}

	orderProductsIDs := make([]int, len(order.Items))
	for i, orderItem := range order.Items {
		orderProductsIDs[i] = orderItem.ProductID
	}

	products, err := os.pr.FindByIDs(orderProductsIDs)
	if err != nil {
		return nil, err
	}

	if len(orderProductsIDs) > len(products) {
		for _, orderProductID := range orderProductsIDs {
			productExists := false

			for _, product := range products {
				if product.ID == orderProductID {
					productExists = true
					break
				}
			}

			if !productExists {
				return nil, base.NewError(
					base.ErrInvalidParameters,
					fmt.Sprintf("product with id = %d doesn't exists", orderProductID),
				)
			}
		}
	}

	newOrder, err := os.or.Create(order)

	return newOrder, nil
}
