package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"restaurant-assistant/pkg/base"
)

const (
	UserTypeDefault = "default"
	UserTypeAdmin   = "admin"
)

type User struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"username"`
	Password    string    `json:"password,omitempty"`
	Email       string    `json:"email,omitempty"`
	Type        string    `json:"type,omitempty"`
	Description string    `json:"description,omitempty"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func (u *User) Format() interface{} {
	return &User{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Description: u.Description,
		Type:        u.Type,
		CreateTime:  u.CreateTime,
		UpdateTime:  u.UpdateTime,
	}
}

func (u *User) Validate() error {
	var err string
	if u.Name == "" {
		err = "incorrect user name"
	}

	if u.Password == "" {
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

// TODO: Add update user method
func (u *User) IsEditableFields() error {
	var err string
	if u.Name == "" {
		err = "incorrect user name"
	}

	if u.Password == "" {
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

func (u *User) IsAdmin() bool {
	return u.Type == UserTypeAdmin
}
