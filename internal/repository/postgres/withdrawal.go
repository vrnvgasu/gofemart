package postgres

import (
	"context"
	"fmt"

	"github.com/vrnvgasu/gofemart/internal/model"
)

func (s *Storage) CreateWithdrawal(ctx context.Context, userID int64, orderNumber string, sum float64) error {
	db := s.dbFromCtx(ctx)
	_, err := db.ExecContext(ctx,
		`INSERT INTO withdrawals (user_id, order_number, sum) VALUES ($1, $2, $3)`,
		userID, orderNumber, sum,
	)
	if err != nil {
		return fmt.Errorf("postgres.CreateWithdrawal Insert: %w", err)
	}

	return nil
}

// GetUserWithdrawals returns all withdrawals for the given user, newest first.
func (s *Storage) GetUserWithdrawals(ctx context.Context, userID int64) ([]model.Withdrawal, error) {
	db := s.dbFromCtx(ctx)

	rows, err := db.QueryContext(ctx,
		`SELECT id, user_id, order_number, sum, processed_at
		 FROM withdrawals WHERE user_id = $1
		 ORDER BY processed_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("postgres.GetUserWithdrawals Query: %w", err)
	}
	defer rows.Close()

	var result []model.Withdrawal
	for rows.Next() {
		var w model.Withdrawal
		if err = rows.Scan(&w.ID, &w.UserID, &w.OrderNumber, &w.Sum, &w.ProcessedAt); err != nil {
			return nil, fmt.Errorf("postgres.GetUserWithdrawals Scan: %w", err)
		}
		result = append(result, w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.GetUserWithdrawals rows.Err: %w", err)
	}

	return result, nil
}
