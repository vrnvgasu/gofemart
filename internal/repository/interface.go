package repository

import (
	"context"

	"github.com/vrnvgasu/gofemart/internal/model"
)

//go:generate mockgen -destination=./mocks/mock.go . Storage
type Storage interface {
	CreateUser(ctx context.Context, login, passwordHash string) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)

	CreateOrder(ctx context.Context, userID int64, number string) error
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)
	GetUserOrders(ctx context.Context, userID int64) ([]model.Order, error)
	UpdateOrderStatus(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error
	GetPendingOrders(ctx context.Context) ([]model.Order, error)

	GetBalance(ctx context.Context, userID int64) (*model.Balance, error)

	CreateWithdrawal(ctx context.Context, userID int64, orderNumber string, sum float64) error
	GetUserWithdrawals(ctx context.Context, userID int64) ([]model.Withdrawal, error)

	DoInTransaction(ctx context.Context, fn func(context.Context) error) error
	Stop() error
}
