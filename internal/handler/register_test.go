package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/config"
	"github.com/vrnvgasu/gofemart/internal/model"
	mockrepository "github.com/vrnvgasu/gofemart/internal/repository/mocks"
	pgstore "github.com/vrnvgasu/gofemart/internal/repository/postgres"
	"github.com/vrnvgasu/gofemart/pkg/jwt"
)

const testSecret = "test-secret"

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret: testSecret,
		LogLevel:  "error",
	}
}

func authHeader(userID int64) string {
	tok, _ := jwt.Generate(userID, testSecret)
	return "Bearer " + tok
}

func TestRegister(t *testing.T) {
	t.Parallel()

	var (
		login    = "login"
		password = "secret"
		userID   = int64(1)
	)

	tests := []struct {
		name           string
		storage        func(*mockrepository.MockStorage) *mockrepository.MockStorage
		body           string
		expectedStatus int
		expectedHeader string
	}{
		{
			name: "Success",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().CreateUser(gomock.Any(), login, gomock.Any()).
					Return(&model.User{
						ID:           userID,
						Login:        login,
						PasswordHash: password,
					}, nil)
				return store
			},
			body:           fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password),
			expectedStatus: http.StatusOK,
			expectedHeader: authHeader(userID),
		},
		{
			name: "BadBody",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				return store
			},
			body:           "dummy",
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "",
		},
		{
			name: "Conflict",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().CreateUser(gomock.Any(), login, gomock.Any()).
					Return(nil, pgstore.ErrLoginConflict)
				return store
			},
			body:           fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password),
			expectedStatus: http.StatusConflict,
			expectedHeader: "",
		},
		{
			name: "Storage error",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().CreateUser(gomock.Any(), login, gomock.Any()).
					Return(nil, errors.New("storage error"))
				return store
			},
			body:           fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password),
			expectedStatus: http.StatusInternalServerError,
			expectedHeader: "",
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
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedHeader, w.Header().Get("Authorization"))
		})
	}

}
