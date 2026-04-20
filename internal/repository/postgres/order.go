package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/vrnvgasu/gofemart/internal/model"
)

func (s *Storage) CreateOrder(ctx context.Context, userID int64, number string) error {
	db := s.dbFromCtx(ctx)

	_, err := db.ExecContext(ctx,
		`INSERT INTO orders (user_id, number) VALUES ($1, $2)`,
		userID, number,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			existing, qErr := s.GetOrderByNumber(ctx, number)
			if qErr != nil {
				return fmt.Errorf("postgres.CreateOrder lookup conflict: %w", qErr)
			}
			if existing != nil && existing.UserID == userID {
				return ErrOrderAlreadyExists
			}
			return ErrOrderConflict
		}
		return fmt.Errorf("postgres.CreateOrder Exec: %w", err)
	}

	return nil
}

func (s *Storage) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	db := s.dbFromCtx(ctx)

	row := db.QueryRowContext(ctx,
		`SELECT id, user_id, number, status, accrual, uploaded_at FROM orders WHERE number = $1`,
		number,
	)

	return scanOrder(row)
}

func (s *Storage) GetUserOrders(ctx context.Context, userID int64) ([]model.Order, error) {
	db := s.dbFromCtx(ctx)

	rows, err := db.QueryContext(ctx,
		`SELECT id, user_id, number, status, accrual, uploaded_at
		 FROM orders WHERE user_id = $1
		 ORDER BY uploaded_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("postgres.GetUserOrders Query: %w", err)
	}
	defer rows.Close()

	var result []model.Order
	for rows.Next() {
		var o model.Order
		if err = rows.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			return nil, fmt.Errorf("postgres.GetUserOrders Scan: %w", err)
		}
		result = append(result, o)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.GetUserOrders rows.Err: %w", err)
	}

	return result, nil
}

func (s *Storage) UpdateOrderStatus(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error {
	db := s.dbFromCtx(ctx)

	_, err := db.ExecContext(ctx,
		`UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`,
		status, accrual, number,
	)
	if err != nil {
		return fmt.Errorf("postgres.UpdateOrderStatus Exec: %w", err)
	}

	return nil
}

func (s *Storage) GetPendingOrders(ctx context.Context) ([]model.Order, error) {
	db := s.dbFromCtx(ctx)

	rows, err := db.QueryContext(ctx,
		`SELECT id, user_id, number, status, accrual, uploaded_at
		 FROM orders WHERE status IN ('NEW', 'PROCESSING')
		 ORDER BY uploaded_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("postgres.GetPendingOrders Query: %w", err)
	}
	defer rows.Close()

	var result []model.Order
	for rows.Next() {
		var o model.Order
		if err = rows.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			return nil, fmt.Errorf("postgres.GetPendingOrders Scan: %w", err)
		}
		result = append(result, o)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.GetPendingOrders rows.Err: %w", err)
	}

	return result, nil
}

func scanOrder(row *sql.Row) (*model.Order, error) {
	var o model.Order
	err := row.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("postgres.scanOrder Scan: %w", err)
	}
	return &o, nil
}
