package postgres

import (
	"context"
	"fmt"

	"github.com/vrnvgasu/gofemart/internal/model"
)

func (s *Storage) GetBalance(ctx context.Context, userID int64) (*model.Balance, error) {
	db := s.dbFromCtx(ctx)

	row := db.QueryRowContext(ctx, `
		SELECT
			COALESCE(o.total, 0) - COALESCE(w.total, 0) AS current,
			COALESCE(w.total, 0)                         AS withdrawn
		FROM
			(SELECT SUM(accrual) AS total FROM orders
			 WHERE user_id = $1 AND status = 'PROCESSED') AS o,
			(SELECT SUM(sum)     AS total FROM withdrawals
			 WHERE user_id = $1) AS w
	`, userID)

	var b model.Balance
	if err := row.Scan(&b.Current, &b.Withdrawn); err != nil {
		return nil, fmt.Errorf("postgres.GetBalance Scan: %w", err)
	}

	return &b, nil
}
