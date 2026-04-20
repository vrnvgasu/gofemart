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
	"golang.org/x/crypto/bcrypt"

	"github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/model"
	mockrepository "github.com/vrnvgasu/gofemart/internal/repository/mocks"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	var (
		login    = "login"
		password = "secret"
		userID   = int64(1)
		hash, _  = bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
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
				store.EXPECT().GetUserByLogin(gomock.Any(), login).
					Return(&model.User{
						ID:           userID,
						Login:        login,
						PasswordHash: string(hash),
					}, nil)
				return store
			},
			body:           fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password),
			expectedStatus: http.StatusOK,
			expectedHeader: authHeader(userID),
		},
		{
			name: "Wrong password",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserByLogin(gomock.Any(), login).
					Return(&model.User{
						ID:           userID,
						Login:        login,
						PasswordHash: string(hash),
					}, nil)
				return store
			},
			body:           fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, "wrong-password"),
			expectedStatus: http.StatusUnauthorized,
			expectedHeader: "",
		},
		{
			name: "User not found",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserByLogin(gomock.Any(), login).
					Return(nil, nil)
				return store
			},
			body:           fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password),
			expectedStatus: http.StatusUnauthorized,
			expectedHeader: "",
		},
		{
			name: "Bad body",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				return store
			},
			body:           "dummy",
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "",
		},
		{
			name: "Storage error",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserByLogin(gomock.Any(), login).
					Return(nil, errors.New("db error"))
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
			req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedHeader, w.Header().Get("Authorization"))
		})
	}

}
