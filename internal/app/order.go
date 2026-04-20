package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vrnvgasu/gofemart/internal/repository/postgres"
	"github.com/vrnvgasu/gofemart/pkg/luhn"
)

// OrderResponse содержит данные о заказе для возврата в HTTP-ответе.
type OrderResponse struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

// CreateOrder принимает новый заказ от пользователя.
// Возвращает true, если заказ новый, false — если уже был загружен этим пользователем.
// Возвращает ConflictError, если номер уже загружен другим пользователем.
func (a *App) CreateOrder(ctx context.Context, userID int64, number string) (bool, error) {
	if !luhn.Valid(number) {
		return false, UnprocessableEntityError()
	}

	err := a.storage.CreateOrder(ctx, userID, number)
	if err != nil {
		if errors.Is(err, postgres.ErrOrderConflict) {
			return false, ConflictError()
		}
		if errors.Is(err, postgres.ErrOrderAlreadyExists) {
			return false, nil
		}

		return false, fmt.Errorf("app.CreateOrder: %w", err)
	}

	return true, nil
}

// GetUserOrders возвращает список заказов пользователя.
// Возвращает NoContentError, если заказов нет.
func (a *App) GetUserOrders(ctx context.Context, userID int64) ([]OrderResponse, error) {
	orders, err := a.storage.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("app.GetUserOrders: %w", err)
	}

	if len(orders) == 0 {
		return nil, NoContentError()
	}

	resp := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		resp = append(resp, OrderResponse{
			Number:     o.Number,
			Status:     string(o.Status),
			Accrual:    o.Accrual,
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
		})
	}

	return resp, nil
}
