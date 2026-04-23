package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey struct{}

func (s *Storage) DoInTransaction(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return fn(ctx)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("postgres.DoInTransaction BeginTx: %w", err)
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()

			panic(p)
		}
	}()

	// Выполняем функцию с транзакцией в контексте
	if err = fn(txCtx); err != nil {
		tx.Rollback()

		return fmt.Errorf("postgres.DoInTransaction do fn: %w", err)
	}

	return tx.Commit()
}

func (s *Storage) dbFromCtx(ctx context.Context) DB {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok && tx != nil {
		return tx
	}

	return s.db
}
