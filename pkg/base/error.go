package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// Server
	ErrInternal = 1

	// Received data
	ErrInvalidParameters = 50

	// Auth
	ErrAuthentication       = 51
	ErrAuthTokenExpired     = 52
	ErrInvalidAuthTokenType = 53
	ErrNoPermissions        = 54
)

type Error struct {
	code    int
	message string
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.message
}

func NewError(code int, message string) Error {
	return Error{
		code:    code,
		message: message,
	}
}

func NewAuthError(message string) Error {
	return Error{
		code:    ErrAuthentication,
		message: message,
	}
}

func NewInternalError(message string) Error {
	return Error{
		code:    ErrInternal,
		message: message,
	}
}

func SendError(c *gin.Context, err error) {
	status := http.StatusInternalServerError

	var resp ErrorResponse
	if e, ok := err.(Error); ok {
		status = mapErrorToStatus(e.code)

		resp = ErrorResponse{
			Code:    e.code,
			Message: e.message,
		}
	} else {
		resp = ErrorResponse{
			Code:    ErrInternal,
			Message: err.Error(),
		}
	}

	c.JSON(status, resp)
}

func SendNewError(c *gin.Context, code int, message string) {
	SendError(c, NewError(code, message))
}

func mapErrorToStatus(code int) int {
	switch code {
	case ErrInternal:
		return http.StatusInternalServerError

	case ErrInvalidParameters:
		return http.StatusUnprocessableEntity

	case ErrAuthentication, ErrAuthTokenExpired, ErrInvalidAuthTokenType:
		return http.StatusUnauthorized

	case ErrNoPermissions:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
