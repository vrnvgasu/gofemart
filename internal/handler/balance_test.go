package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/model"
	mockrepository "github.com/vrnvgasu/gofemart/internal/repository/mocks"
)

func TestGetBalance(t *testing.T) {
	t.Parallel()

	var (
		userID  = int64(1)
		balance = model.Balance{Current: 500.5, Withdrawn: 42}
	)

	tests := []struct {
		name           string
		storage        func(*mockrepository.MockStorage) *mockrepository.MockStorage
		body           string
		expectedStatus int
		expectedBody   *app.BalanceResponse
	}{
		{
			name: "Success",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetBalance(gomock.Any(), userID).
					Return(&balance, nil)
				return store
			},
			expectedStatus: http.StatusOK,
			expectedBody: &app.BalanceResponse{
				Current:   balance.Current,
				Withdrawn: balance.Withdrawn,
			},
		},
		{
			name: "Storage error",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetBalance(gomock.Any(), userID).
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
			req := httptest.NewRequest(http.MethodGet, "/api/user/balance", strings.NewReader(tt.body))
			req.Header.Set("Authorization", authHeader(1))
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var resp app.BalanceResponse
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
				assert.Equal(t, *tt.expectedBody, resp)
			}
		})
	}
}
