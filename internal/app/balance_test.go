package app

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/gofemart/internal/model"
	mockrepository "github.com/vrnvgasu/gofemart/internal/repository/mocks"
)

func TestGetBalance(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockrepository.NewMockStorage(ctrl)
	store.EXPECT().GetBalance(gomock.Any(), int64(1)).
		Return(&model.Balance{
			Current:   10,
			Withdrawn: 10,
		}, nil)
	a := App{
		storage: store,
	}
	res, err := a.GetBalance(t.Context(), int64(1))
	require.NoError(t, err)
	require.Equal(t, BalanceResponse{
		Current:   10,
		Withdrawn: 10,
	}, res)
}
