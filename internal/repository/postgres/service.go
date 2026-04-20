package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Storage struct {
	db *sql.DB
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Start(ctx context.Context, dsn string) (err error) {
	classifier := newPostgresErrorClassifier()

	err = withRetry(func() error {
		s.db, err = sql.Open("pgx", dsn)
		return err
	}, classifier)
	if err != nil {
		return fmt.Errorf("postgres.Start Open: %w", err)
	}

	err = withRetry(func() error { return s.db.PingContext(ctx) }, classifier)
	if err != nil {
		return fmt.Errorf("postgres.Start Ping: %w", err)
	}

	err = withRetry(func() error { return s.migrate(ctx) }, classifier)
	if err != nil {
		return fmt.Errorf("postgres.Start Migrate: %w", err)
	}

	return nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}
