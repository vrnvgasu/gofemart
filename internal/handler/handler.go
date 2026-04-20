package handler

import (
	"context"

	"github.com/vrnvgasu/gofemart/internal/app"
)

type App interface {
	GetBalance(ctx context.Context, userID int64) (app.BalanceResponse, error)

	CreateOrder(ctx context.Context, userID int64, number string) (bool, error)
	GetUserOrders(ctx context.Context, userID int64) ([]app.OrderResponse, error)

	CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int64) ([]app.WithdrawalResponse, error)

	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, login, password string) (string, error)
}

type Handler struct {
	app App
}

func NewHandler(app App) *Handler {
	return &Handler{
		app: app,
	}
}
