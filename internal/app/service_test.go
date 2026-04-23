package app

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/gofemart/internal/config"
	mockrepository "github.com/vrnvgasu/gofemart/internal/repository/mocks"
)

func TestNewApp(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := mockrepository.NewMockStorage(ctrl)
	cnf := &config.Config{}

	a := App{
		storage: storage,
		cfg:     cnf,
	}

	require.Equal(t, a, *NewApp(storage, cnf))
}
