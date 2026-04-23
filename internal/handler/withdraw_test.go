package handler

import (
	"encoding/json"
	"errors"
	"fmt"
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
)

func TestSaveWithdraw(t *testing.T) {
	t.Parallel()

	var (
		userID = int64(1)
		number = "2377225624"
		sum    = float64(100)
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
				store.EXPECT().DoInTransaction(gomock.Any(), gomock.Any()).
					Return(nil)
				return store
			},
			body:           fmt.Sprintf(`{"order": "%s", "sum": %f}`, number, sum),
			expectedStatus: http.StatusOK,
		},
		{
			name: "Wrong balance",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().DoInTransaction(gomock.Any(), gomock.Any()).
					Return(app.PaymentRequiredError())
				return store
			},
			body:           fmt.Sprintf(`{"order": "%s", "sum": %f}`, number, sum),
			expectedStatus: http.StatusPaymentRequired,
		},
		{
			name: "Wrong order",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				return store
			},
			body:           fmt.Sprintf(`{"order": "%s", "sum": %f}`, "12345678900", sum),
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Bad body",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				return store
			},
			body:           "dummy",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Storage error",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().DoInTransaction(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
				return store
			},
			body:           fmt.Sprintf(`{"order": "%s", "sum": %f}`, number, sum),
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
			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(tt.body))
			req.Header.Set("Authorization", authHeader(userID))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetWithdrawals(t *testing.T) {
	t.Parallel()

	var (
		userID = int64(1)
		number = "2377225624"
		sum    = float64(100)
	)

	tests := []struct {
		name             string
		storage          func(*mockrepository.MockStorage) *mockrepository.MockStorage
		body             string
		expectedStatus   int
		expectedResponse []app.WithdrawalResponse
	}{
		{
			name: "No records",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserWithdrawals(gomock.Any(), userID).
					Return([]model.Withdrawal{}, nil)
				return store
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "Has records",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserWithdrawals(gomock.Any(), userID).
					Return([]model.Withdrawal{
						{
							ID:          1,
							UserID:      userID,
							OrderNumber: number,
							Sum:         sum,
							ProcessedAt: time.Now().Truncate(24 * time.Hour),
						},
					}, nil)
				return store
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []app.WithdrawalResponse{
				{
					Order:       number,
					Sum:         sum,
					ProcessedAt: time.Now().Truncate(24 * time.Hour).Format(time.RFC3339),
				},
			},
		},
		{
			name: "Storage error",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserWithdrawals(gomock.Any(), userID).
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
			req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", http.NoBody)
			req.Header.Set("Authorization", authHeader(userID))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedResponse != nil {
				var resp []app.WithdrawalResponse
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
				assert.Equal(t, tt.expectedResponse, resp)
			}
		})
	}
}
