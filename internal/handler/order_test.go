package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/model"
	mockrepository "github.com/vrnvgasu/gofemart/internal/repository/mocks"
	pgstore "github.com/vrnvgasu/gofemart/internal/repository/postgres"
)

func TestUploadOrder(t *testing.T) {
	t.Parallel()

	var (
		userID = int64(1)
		number = "12345678903"
	)

	tests := []struct {
		name           string
		storage        func(*mockrepository.MockStorage) *mockrepository.MockStorage
		body           string
		expectedStatus int
	}{
		{
			name: "Success",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().CreateOrder(gomock.Any(), userID, number).
					Return(nil)
				return store
			},
			body:           number,
			expectedStatus: http.StatusAccepted,
		},
		{
			name: "User has order",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().CreateOrder(gomock.Any(), userID, number).
					Return(pgstore.ErrOrderAlreadyExists)
				return store
			},
			body:           number,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Another ser has order",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().CreateOrder(gomock.Any(), userID, number).
					Return(pgstore.ErrOrderConflict)
				return store
			},
			body:           number,
			expectedStatus: http.StatusConflict,
		},
		{
			name: "Invalid luhn",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				return store
			},
			body:           "12345678900",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Storage error",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().CreateOrder(gomock.Any(), userID, number).
					Return(errors.New("db error"))
				return store
			},
			body:           number,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Empty body",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				return store
			},
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			controller := gomock.NewController(t)
			defer controller.Finish()
			store := mockrepository.NewMockStorage(controller)

			appService := app.NewApp(tt.storage(store), testConfig())
			r := NewRouter(NewHandler(appService), testConfig())
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader(tt.body))
			req.Header.Set("Authorization", authHeader(1))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetOrders(t *testing.T) {
	t.Parallel()

	var (
		userID = int64(1)
		number = "12345678903"
	)

	tests := []struct {
		name           string
		storage        func(*mockrepository.MockStorage) *mockrepository.MockStorage
		expectedStatus int
		expectedBody   []app.OrderResponse
	}{
		{
			name: "No orders",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserOrders(gomock.Any(), userID).
					Return(nil, nil)
				return store
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "Has orders",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserOrders(gomock.Any(), userID).
					Return([]model.Order{
						{
							ID:         1,
							UserID:     userID,
							Number:     number,
							Status:     model.OrderStatusNew,
							Accrual:    nil,
							UploadedAt: time.Now().Truncate(24 * time.Hour),
						},
					}, nil)
				return store
			},
			expectedStatus: http.StatusOK,
			expectedBody: []app.OrderResponse{
				{
					Number:     number,
					Status:     string(model.OrderStatusNew),
					Accrual:    nil,
					UploadedAt: time.Now().Truncate(24 * time.Hour).Format(time.RFC3339),
				},
			},
		},
		{
			name: "Storage error",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserOrders(gomock.Any(), userID).
					Return(nil, errors.New("db error"))
				return store
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			controller := gomock.NewController(t)
			defer controller.Finish()
			store := mockrepository.NewMockStorage(controller)

			appService := app.NewApp(tt.storage(store), testConfig())
			r := NewRouter(NewHandler(appService), testConfig())
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", http.NoBody)
			req.Header.Set("Authorization", authHeader(1))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var resp []app.OrderResponse
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
				assert.Equal(t, tt.expectedBody, resp)
			}
		})
	}
}
