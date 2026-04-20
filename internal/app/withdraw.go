package app

import (
	"context"
	"fmt"
	"time"

	"github.com/vrnvgasu/gofemart/pkg/luhn"
)

type WithdrawalResponse struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (a *App) CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) error {
	if !luhn.Valid(order) {
		return UnprocessableEntityError()
	}

	err := a.storage.DoInTransaction(ctx, func(txCtx context.Context) error {
		balance, err := a.storage.GetBalance(txCtx, userID)
		if err != nil {
			return fmt.Errorf("app.CreateWithdrawal GetBalance: %w", err)
		}

		if balance.Current < sum {
			return PaymentRequiredError()
		}

		if err = a.storage.CreateWithdrawal(ctx, userID, order, sum); err != nil {
			return fmt.Errorf("app.CreateWithdrawal CreateWithdrawal: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("app.CreateWithdrawal DoInTransaction: %w", err)
	}

	return nil
}

func (a *App) GetWithdrawals(ctx context.Context, userID int64) ([]WithdrawalResponse, error) {
	withdrawals, err := a.storage.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("app.GetWithdrawals: %w", err)
	}

	if len(withdrawals) == 0 {
		return nil, NoContentError()
	}

	resp := make([]WithdrawalResponse, 0, len(withdrawals))
	for _, w := range withdrawals {
		resp = append(resp, WithdrawalResponse{
			Order:       w.OrderNumber,
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
		})
	}

	return resp, nil
}
