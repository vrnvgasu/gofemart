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
)

func TestCreateWithdrawal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		number      string
		storage     func(*mockrepository.MockStorage) repository.Storage
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name:   "success",
			number: "12345678903",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().DoInTransaction(gomock.Any(), gomock.Any()).
					Return(nil)
				return store
			},
			expectedErr: require.NoError,
		},
		{
			name:   "db error",
			number: "12345678903",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().DoInTransaction(gomock.Any(), gomock.Any()).
					Return(errors.New("db error"))
				return store
			},
			expectedErr: require.Error,
		},
		{
			name:   "wrong number",
			number: "12345678900",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
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
			err := a.CreateWithdrawal(t.Context(), 1, tt.number, 10)
			tt.expectedErr(t, err)
		})
	}
}

func TestGetWithdrawals(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		storage     func(*mockrepository.MockStorage) repository.Storage
		expected    []WithdrawalResponse
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name: "success",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetUserWithdrawals(gomock.Any(), int64(1)).
					Return([]model.Withdrawal{
						{
							ID:          1,
							UserID:      1,
							OrderNumber: "1",
							Sum:         10,
							ProcessedAt: time.Now().Truncate(24 * time.Hour),
						},
					}, nil)
				return store
			},
			expected: []WithdrawalResponse{
				{
					Order:       "1",
					Sum:         10,
					ProcessedAt: time.Now().Truncate(24 * time.Hour).Format(time.RFC3339),
				},
			},
			expectedErr: require.NoError,
		},
		{
			name: "db error",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetUserWithdrawals(gomock.Any(), int64(1)).
					Return(nil, errors.New("db error"))
				return store
			},
			expectedErr: require.Error,
		},
		{
			name: "no records",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetUserWithdrawals(gomock.Any(), int64(1)).
					Return([]model.Withdrawal{}, nil)
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
			seq, err := a.GetWithdrawals(t.Context(), 1)
			tt.expectedErr(t, err)
			var got []WithdrawalResponse
			if seq != nil {
				got = slices.Collect(seq)
			}
			require.Equal(t, tt.expected, got)
		})
	}
}
