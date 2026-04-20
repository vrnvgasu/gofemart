package app

import (
	"context"
	"fmt"
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

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
