package app

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func() error
		expErr ServiceError
	}{
		{
			name: "UnauthorizedError",
			fn:   UnauthorizedError,
			expErr: ServiceError{
				Type:     ErrUnauthorized,
				Message:  http.StatusText(http.StatusUnauthorized),
				HTTPCode: http.StatusUnauthorized,
			},
		},
		{
			name: "PaymentRequiredError",
			fn:   PaymentRequiredError,
			expErr: ServiceError{
				Type:     ErrPaymentRequired,
				Message:  http.StatusText(http.StatusPaymentRequired),
				HTTPCode: http.StatusPaymentRequired,
			},
		},
		{
			name: "NoContentError",
			fn:   NoContentError,
			expErr: ServiceError{
				Type:     ErrNoContent,
				Message:  http.StatusText(http.StatusNoContent),
				HTTPCode: http.StatusNoContent,
			},
		},
		{
			name: "ConflictError",
			fn:   ConflictError,
			expErr: ServiceError{
				Type:     ErrConflict,
				Message:  http.StatusText(http.StatusConflict),
				HTTPCode: http.StatusConflict,
			},
		},
		{
			name: "UnprocessableEntityError",
			fn:   UnprocessableEntityError,
			expErr: ServiceError{
				Type:     ErrUnprocessableEntity,
				Message:  http.StatusText(http.StatusUnprocessableEntity),
				HTTPCode: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "BadRequestError",
			fn:   BadRequestError,
			expErr: ServiceError{
				Type:     ErrBadRequest,
				Message:  http.StatusText(http.StatusBadRequest),
				HTTPCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expErr, *(tt.fn().(*ServiceError)))
		})
	}
}
