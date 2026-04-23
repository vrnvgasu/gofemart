package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"github.com/vrnvgasu/gofemart/internal/config"
	"github.com/vrnvgasu/gofemart/internal/model"
	mockrepository "github.com/vrnvgasu/gofemart/internal/repository/mocks"
	"github.com/vrnvgasu/gofemart/pkg/jwt"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	var (
		userID     = int64(1)
		testSecret = "test-secret"
	)

	tests := []struct {
		name        string
		login       string
		password    string
		storage     func(*mockrepository.MockStorage) *mockrepository.MockStorage
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name:     "Success",
			login:    "login",
			password: "secret",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
				store.EXPECT().GetUserByLogin(gomock.Any(), "login").
					Return(&model.User{
						ID:           userID,
						Login:        "login",
						PasswordHash: string(hash),
					}, nil)
				return store
			},
			expectedErr: require.NoError,
		},
		{
			name:     "Wrong password",
			login:    "login",
			password: "wrong-password",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
				store.EXPECT().GetUserByLogin(gomock.Any(), "login").
					Return(&model.User{
						ID:           userID,
						Login:        "login",
						PasswordHash: string(hash),
					}, nil)
				return store
			},
			expectedErr: require.Error,
		},
		{
			name:     "User not found",
			login:    "login",
			password: "wrong-password",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().GetUserByLogin(gomock.Any(), "login").
					Return(nil, nil)
				return store
			},
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			controller := gomock.NewController(t)
			defer controller.Finish()

			a := App{
				storage: tt.storage(mockrepository.NewMockStorage(controller)),
				cfg: &config.Config{
					JWTSecret: testSecret,
				},
			}
			token, err := a.Login(t.Context(), tt.login, tt.password)
			tt.expectedErr(t, err)

			if err == nil {
				tok, _ := jwt.Generate(userID, testSecret)
				require.Equal(t, tok, token)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	t.Parallel()

	var (
		userID     = int64(1)
		testSecret = "test-secret"
	)

	tests := []struct {
		name        string
		login       string
		password    string
		storage     func(*mockrepository.MockStorage) *mockrepository.MockStorage
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name:     "Success",
			login:    "login",
			password: "secret",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().CreateUser(gomock.Any(), "login", gomock.Any()).
					Return(&model.User{
						ID: userID,
					}, nil)
				return store
			},
			expectedErr: require.NoError,
		},
		{
			name:     "DB error",
			login:    "login",
			password: "secret",
			storage: func(store *mockrepository.MockStorage) *mockrepository.MockStorage {
				store.EXPECT().CreateUser(gomock.Any(), "login", gomock.Any()).
					Return(nil, errors.New("database error"))
				return store
			},
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			controller := gomock.NewController(t)
			defer controller.Finish()

			a := App{
				storage: tt.storage(mockrepository.NewMockStorage(controller)),
				cfg: &config.Config{
					JWTSecret: testSecret,
				},
			}
			token, err := a.Register(t.Context(), tt.login, tt.password)
			tt.expectedErr(t, err)

			if err == nil {
				tok, _ := jwt.Generate(userID, testSecret)
				require.Equal(t, tok, token)
			}
		})
	}
}
