package app

import (
	"context"
	"fmt"
)

// BalanceResponse содержит данные о балансе пользователя для HTTP-ответа.
type BalanceResponse struct {
	// Current — текущий остаток баллов.
	Current float64 `json:"current"`
	// Withdrawn — общая сумма списанных баллов.
	Withdrawn float64 `json:"withdrawn"`
}

// GetBalance возвращает текущий баланс пользователя.
func (a *App) GetBalance(ctx context.Context, userID int64) (BalanceResponse, error) {
	balance, err := a.storage.GetBalance(ctx, userID)
	if err != nil {
		return BalanceResponse{}, fmt.Errorf("app.GetBalance: %w", err)
	}

	return BalanceResponse{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}, nil
}
