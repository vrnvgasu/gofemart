// Package handler содержит HTTP-обработчики приложения.
package handler

import (
	"context"
	"iter"

	"github.com/vrnvgasu/gofemart/internal/app"
)

// App описывает интерфейс бизнес-логики, которую используют обработчики.
type App interface {
	GetBalance(ctx context.Context, userID int64) (app.BalanceResponse, error)

	CreateOrder(ctx context.Context, userID int64, number string) (bool, error)
	GetUserOrders(ctx context.Context, userID int64) (iter.Seq[app.OrderResponse], error)

	CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int64) (iter.Seq[app.WithdrawalResponse], error)

	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, login, password string) (string, error)
}

// Handler содержит зависимости для HTTP-обработчиков.
type Handler struct {
	app App
}

// NewHandler создает новый Handler с переданным сервисом.
func NewHandler(app App) *Handler {
	return &Handler{
		app: app,
	}
}
