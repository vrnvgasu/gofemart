package accrual_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/gofemart/internal/accrual"
	"github.com/vrnvgasu/gofemart/pkg/helper"
)

func TestGetOrder(t *testing.T) {
	t.Parallel()
	var (
		order = "12345678903"
		url   = "/api/orders/" + order
	)

	tests := []struct {
		name         string
		handler      http.Handler
		expectedErr  require.ErrorAssertionFunc
		expectedInfo *accrual.OrderInfo
	}{
		{
			name: "processed",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, url, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(accrual.OrderInfo{
					Order:   order,
					Status:  accrual.StatusProcessed,
					Accrual: helper.NewRefFloat64(500.0),
				})
			}),
			expectedErr: require.NoError,
			expectedInfo: &accrual.OrderInfo{
				Order:   order,
				Status:  accrual.StatusProcessed,
				Accrual: helper.NewRefFloat64(500.0),
			},
		},
		{
			name: "processing",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, url, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(accrual.OrderInfo{
					Order:  order,
					Status: accrual.StatusProcessing,
				})
			}),
			expectedErr: require.NoError,
			expectedInfo: &accrual.OrderInfo{
				Order:  order,
				Status: accrual.StatusProcessing,
			},
		},
		{
			name: "invalid",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, url, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(accrual.OrderInfo{
					Order:  order,
					Status: accrual.StatusInvalid,
				})
			}),
			expectedErr: require.NoError,
			expectedInfo: &accrual.OrderInfo{
				Order:  order,
				Status: accrual.StatusInvalid,
			},
		},
		{
			name: "not registered",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, url, r.URL.Path)
				w.WriteHeader(http.StatusNoContent)
			}),
			expectedErr: require.NoError,
		},
		{
			name: "server error",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, url, r.URL.Path)
				w.WriteHeader(http.StatusInternalServerError)
			}),
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			client := accrual.NewClient(srv.URL)
			info, err := client.GetOrder(context.Background(), "12345678903")
			tt.expectedErr(t, err)
			if tt.expectedInfo != nil {
				require.NotNil(t, info)
				assert.Equal(t, tt.expectedInfo.Status, info.Status)

				if tt.expectedInfo.Accrual != nil {
					assert.Equal(t, *tt.expectedInfo.Accrual, *info.Accrual)
				}
			} else {
				require.Nil(t, info)
			}
		})
	}
}

func TestErrTooManyRequests_Error(t *testing.T) {
	t.Parallel()

	err := &accrual.ErrTooManyRequests{RetryAfter: 30 * time.Second}
	assert.Contains(t, err.Error(), "30s")
}

func TestGetOrder_TooManyRequests(t *testing.T) {
	t.Parallel()
	url := "/api/orders/12345678903"

	tests := []struct {
		name                       string
		handler                    http.Handler
		expecterRetryAfterResponse time.Duration
	}{
		{
			name: "30s",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, url, r.URL.Path)
				w.Header().Set("Retry-After", "30")
				w.WriteHeader(http.StatusTooManyRequests)
			}),
			expecterRetryAfterResponse: 30 * time.Second,
		},
		{
			name: "60s",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, url, r.URL.Path)
				w.WriteHeader(http.StatusTooManyRequests)
			}),
			expecterRetryAfterResponse: 60 * time.Second,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			client := accrual.NewClient(srv.URL)
			info, err := client.GetOrder(context.Background(), "12345678903")
			require.Error(t, err)
			require.Nil(t, info)

			var tooMany *accrual.ErrTooManyRequests
			assert.ErrorAs(t, err, &tooMany)
			assert.Equal(t, tt.expecterRetryAfterResponse, tooMany.RetryAfter)
		})
	}
}
