// Package app содержит бизнес-логику приложения.
package app

import (
	"fmt"
	"net/http"
)

// ServiceErrorType — тип ошибки сервисного слоя.
type ServiceErrorType string

const (
	// ErrInternal — внутренняя ошибка сервера.
	ErrInternal ServiceErrorType = "internal server error"
	// ErrUnprocessableEntity — переданные данные не прошли валидацию.
	ErrUnprocessableEntity ServiceErrorType = "unprocessable entity"
	// ErrBadRequest — некорректный запрос.
	ErrBadRequest ServiceErrorType = "bad request"
	// ErrUnauthorized — пользователь не авторизован.
	ErrUnauthorized ServiceErrorType = "unauthorized"
	// ErrConflict — конфликт данных, например дублирование.
	ErrConflict ServiceErrorType = "conflict"
	// ErrNoContent — данные отсутствуют.
	ErrNoContent ServiceErrorType = "no content"
	// ErrPaymentRequired — недостаточно средств на счете.
	ErrPaymentRequired ServiceErrorType = "payment required"
)

// ServiceError — ошибка сервисного слоя с HTTP-кодом.
type ServiceError struct {
	Type     ServiceErrorType
	Message  string
	HTTPCode int
}

// Error возвращает строковое представление ошибки.
func (s *ServiceError) Error() string {
	return fmt.Sprintf("[%s]: %s", s.Type, s.Message)
}

// UnauthorizedError возвращает ошибку 401 Unauthorized.
func UnauthorizedError() error {
	return &ServiceError{
		Type:     ErrUnauthorized,
		Message:  http.StatusText(http.StatusUnauthorized),
		HTTPCode: http.StatusUnauthorized,
	}
}

// PaymentRequiredError возвращает ошибку 402 Payment Required.
func PaymentRequiredError() error {
	return &ServiceError{
		Type:     ErrPaymentRequired,
		Message:  http.StatusText(http.StatusPaymentRequired),
		HTTPCode: http.StatusPaymentRequired,
	}
}

// NoContentError возвращает ошибку 204 No Content.
func NoContentError() error {
	return &ServiceError{
		Type:     ErrNoContent,
		Message:  http.StatusText(http.StatusNoContent),
		HTTPCode: http.StatusNoContent,
	}
}

// ConflictError возвращает ошибку 409 Conflict.
func ConflictError() error {
	return &ServiceError{
		Type:     ErrConflict,
		Message:  http.StatusText(http.StatusConflict),
		HTTPCode: http.StatusConflict,
	}
}

// UnprocessableEntityError возвращает ошибку 422 Unprocessable Entity.
func UnprocessableEntityError() error {
	return &ServiceError{
		Type:     ErrUnprocessableEntity,
		Message:  http.StatusText(http.StatusUnprocessableEntity),
		HTTPCode: http.StatusUnprocessableEntity,
	}
}

// BadRequestError возвращает ошибку 400 Bad Request.
func BadRequestError() error {
	return &ServiceError{
		Type:     ErrBadRequest,
		Message:  http.StatusText(http.StatusBadRequest),
		HTTPCode: http.StatusBadRequest,
	}
}
