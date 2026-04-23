package app

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/gofemart/internal/model"
	"github.com/vrnvgasu/gofemart/internal/repository"
	mockrepository "github.com/vrnvgasu/gofemart/internal/repository/mocks"
	"github.com/vrnvgasu/gofemart/internal/repository/postgres"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		number      string
		storage     func(*mockrepository.MockStorage) repository.Storage
		expected    bool
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name:   "success",
			number: "12345678903",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().CreateOrder(gomock.Any(), int64(1), "12345678903").
					Return(nil)
				return store
			},
			expected:    true,
			expectedErr: require.NoError,
		},
		{
			name:   "exist",
			number: "12345678903",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().CreateOrder(gomock.Any(), int64(1), "12345678903").
					Return(postgres.ErrOrderAlreadyExists)
				return store
			},
			expected:    false,
			expectedErr: require.NoError,
		},
		{
			name:   "conflict",
			number: "12345678903",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().CreateOrder(gomock.Any(), int64(1), "12345678903").
					Return(postgres.ErrOrderConflict)
				return store
			},
			expected:    false,
			expectedErr: require.Error,
		},
		{
			name:   "db error",
			number: "12345678903",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().CreateOrder(gomock.Any(), int64(1), "12345678903").
					Return(errors.New("db error"))
				return store
			},
			expected:    false,
			expectedErr: require.Error,
		},
		{
			name:   "wrong number",
			number: "12345678900",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				return store
			},
			expected:    false,
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			a := App{
				storage: tt.storage(mockrepository.NewMockStorage(ctrl)),
			}
			ok, err := a.CreateOrder(t.Context(), 1, tt.number)
			require.Equal(t, tt.expected, ok)
			tt.expectedErr(t, err)
		})
	}
}

func TestGetUserOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		storage     func(*mockrepository.MockStorage) repository.Storage
		expected    []OrderResponse
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetUserOrders(gomock.Any(), int64(1)).
					Return([]model.Order{{
						Number:     "1",
						Status:     model.OrderStatusNew,
						UploadedAt: time.Now().Truncate(24 * time.Hour),
					}}, nil)
				return store
			},
			expected: []OrderResponse{
				{
					Number:     "1",
					Status:     string(model.OrderStatusNew),
					Accrual:    nil,
					UploadedAt: time.Now().Truncate(24 * time.Hour).Format(time.RFC3339),
				},
			},
			expectedErr: require.NoError,
		},
		{
			name: "no orders",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetUserOrders(gomock.Any(), int64(1)).
					Return([]model.Order{}, nil)
				return store
			},
			expectedErr: require.Error,
		},
		{
			name: "db error",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetUserOrders(gomock.Any(), int64(1)).
					Return(nil, errors.New("db error"))
				return store
			},
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			a := App{
				storage: tt.storage(mockrepository.NewMockStorage(ctrl)),
			}
			seq, err := a.GetUserOrders(t.Context(), 1)
			tt.expectedErr(t, err)
			var got []OrderResponse
			if seq != nil {
				got = slices.Collect(seq)
			}
			require.Equal(t, tt.expected, got)
		})
	}
}
