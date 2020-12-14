package entity

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	"restaurant-assistant/pkg/base"
)

type OrderItem struct {
	ProductID   int    `json:"product_id"`
	Number      int    `json:"number"`
	Description string `json:"description,omitempty"`
}

type Order struct {
	ID         int         `json:"id,omitempty"`
	Items      []OrderItem `json:"items"`
	UserID     uuid.UUID   `json:"user_id,omitempty"`
	CreateTime time.Time   `json:"create_time"`
	UpdateTime time.Time   `json:"update_time"`
}

func (o *Order) Validate() error {
	var err string
	getInvalidParamError := func(message string) error {
		return base.NewError(
			base.ErrInvalidParameters,
			message,
		)
	}

	if len(o.Items) == 0 {
		return getInvalidParamError("the number of order items must be greater than zero")
	}

	for index, orderItem := range o.Items {
		if orderItem.ProductID <= 0 {
			return getInvalidParamError(
				fmt.Sprintf("incorrect product id = %d at index %d",
					orderItem.ProductID,
					index,
				),
			)
		}

		if orderItem.Number <= 0 {
			return getInvalidParamError(
				fmt.Sprintf("incorrect number (%d) of product with id = %d",
					orderItem.Number,
					orderItem.ProductID,
				),
			)
		}
	}

	if err == "" {
		return nil
	}

	return getInvalidParamError(err)
}

// todo: create refresh token method
func (o *Order) Format() interface{} {
	return &Order{
		ID:         o.ID,
		Items:      o.Items,
		UserID:     o.UserID,
		CreateTime: o.CreateTime,
		UpdateTime: o.UpdateTime,
	}
}
