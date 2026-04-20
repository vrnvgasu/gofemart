package app

import (
	"fmt"
	"net/http"
)

type ServiceErrorType string

const (
	ErrNotFound            ServiceErrorType = "not found"
	ErrInternal            ServiceErrorType = "internal server error"
	ErrUnprocessableEntity ServiceErrorType = "unprocessable entity"
	ErrBadRequest          ServiceErrorType = "bad request"
	ErrUnauthorized        ServiceErrorType = "unauthorized"
	ErrConflict            ServiceErrorType = "conflict"
	ErrNoContent           ServiceErrorType = "no content"
	ErrPaymentRequired     ServiceErrorType = "payment required"
)

type ServiceError struct {
	Type     ServiceErrorType
	Message  string
	HTTPCode int
}

func (s *ServiceError) Error() string {
	return fmt.Sprintf("[%s]: %s", s.Type, s.Message)
}

func UnauthorizedError() error {
	return &ServiceError{
		Type:     ErrUnauthorized,
		Message:  http.StatusText(http.StatusUnauthorized),
		HTTPCode: http.StatusUnauthorized,
	}
}

func PaymentRequiredError() error {
	return &ServiceError{
		Type:     ErrPaymentRequired,
		Message:  http.StatusText(http.StatusPaymentRequired),
		HTTPCode: http.StatusPaymentRequired,
	}
}

func NotFoundError() error {
	return &ServiceError{
		Type:     ErrNotFound,
		Message:  http.StatusText(http.StatusNotFound),
		HTTPCode: http.StatusNotFound,
	}
}

func NoContentError() error {
	return &ServiceError{
		Type:     ErrNoContent,
		Message:  http.StatusText(http.StatusNoContent),
		HTTPCode: http.StatusNoContent,
	}
}

func ConflictError() error {
	return &ServiceError{
		Type:     ErrConflict,
		Message:  http.StatusText(http.StatusConflict),
		HTTPCode: http.StatusConflict,
	}
}

func UnprocessableEntityError() error {
	return &ServiceError{
		Type:     ErrUnprocessableEntity,
		Message:  http.StatusText(http.StatusUnprocessableEntity),
		HTTPCode: http.StatusUnprocessableEntity,
	}
}

func BadRequestError() error {
	return &ServiceError{
		Type:     ErrBadRequest,
		Message:  http.StatusText(http.StatusBadRequest),
		HTTPCode: http.StatusBadRequest,
	}
}

func InternalError() error {
	return &ServiceError{
		Type:     ErrInternal,
		Message:  http.StatusText(http.StatusInternalServerError),
		HTTPCode: http.StatusInternalServerError,
	}
}
