// Package repository определяет интерфейс хранилища данных.
package repository

import (
	"context"

	"github.com/vrnvgasu/gofemart/internal/model"
)

//go:generate mockgen -destination=./mocks/mock.go -package=mocks . Storage

// Storage описывает методы для работы с базой данных.
type Storage interface {
	// CreateUser создает нового пользователя с указанным логином и хешем пароля.
	CreateUser(ctx context.Context, login, passwordHash string) (*model.User, error)
	// GetUserByLogin возвращает пользователя по логину. Возвращает nil, если не найден.
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)

	// CreateOrder создает новый заказ для пользователя.
	CreateOrder(ctx context.Context, userID int64, number string) error
	// GetOrderByNumber возвращает заказ по его номеру.
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)
	// GetUserOrders возвращает список заказов пользователя, отсортированных по дате загрузки.
	GetUserOrders(ctx context.Context, userID int64) ([]model.Order, error)
	// UpdateOrderStatus обновляет статус заказа и начисленные баллы.
	UpdateOrderStatus(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error
	// GetPendingOrders возвращает заказы, ожидающие обработки.
	GetPendingOrders(ctx context.Context) ([]model.Order, error)

	// GetBalance возвращает текущий баланс пользователя.
	GetBalance(ctx context.Context, userID int64) (*model.Balance, error)

	// CreateWithdrawal создает запись о списании баллов.
	CreateWithdrawal(ctx context.Context, userID int64, orderNumber string, sum float64) error
	// GetUserWithdrawals возвращает список всех списаний пользователя.
	GetUserWithdrawals(ctx context.Context, userID int64) ([]model.Withdrawal, error)

	// DoInTransaction выполняет переданную функцию в рамках транзакции.
	DoInTransaction(ctx context.Context, fn func(context.Context) error) error
	// Stop закрывает соединение с базой данных.
	Stop() error
}
