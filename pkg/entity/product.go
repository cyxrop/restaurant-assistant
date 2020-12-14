package entity

import (
	"time"

	"restaurant-assistant/pkg/base"
)

type Product struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func (p *Product) Validate() error {
	var err string
	if p.Name == "" {
		err = "incorrect user name"
	}

	if err == "" {
		return nil
	}

	return base.NewError(
		base.ErrInvalidParameters,
		err,
	)
}

func (p *Product) Format() interface{} {
	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreateTime:  p.CreateTime,
		UpdateTime:  p.UpdateTime,
	}
}
