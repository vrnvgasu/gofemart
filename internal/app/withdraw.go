package app

import (
	"context"
	"fmt"
	"iter"
	"time"

	"github.com/vrnvgasu/gofemart/pkg/luhn"
)

// WithdrawalResponse содержит данные о списании для HTTP-ответа.
type WithdrawalResponse struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

// CreateWithdrawal списывает баллы со счета пользователя в счет оплаты заказа.
// Возвращает PaymentRequiredError, если на счете недостаточно средств.
// Возвращает UnprocessableEntityError, если номер заказа не прошел проверку Луна.
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

		if err = a.storage.CreateWithdrawal(txCtx, userID, order, sum); err != nil {
			return fmt.Errorf("app.CreateWithdrawal CreateWithdrawal: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("app.CreateWithdrawal DoInTransaction: %w", err)
	}

	return nil
}

// GetWithdrawals возвращает итератор по всем списаниям пользователя.
// Возвращает NoContentError, если списаний не было.
func (a *App) GetWithdrawals(ctx context.Context, userID int64) (iter.Seq[WithdrawalResponse], error) {
	withdrawals, err := a.storage.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("app.GetWithdrawals: %w", err)
	}

	if len(withdrawals) == 0 {
		return nil, NoContentError()
	}

	return func(yield func(WithdrawalResponse) bool) {
		for _, w := range withdrawals {
			if !yield(WithdrawalResponse{
				Order:       w.OrderNumber,
				Sum:         w.Sum,
				ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
			}) {
				return
			}
		}
	}, nil
}
