// Package postgres реализует хранилище данных на основе PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB описывает минимальный интерфейс для работы с базой данных.
// Используется для поддержки как обычного соединения, так и транзакции.
type DB interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// Storage реализует интерфейс repository.Storage для PostgreSQL.
type Storage struct {
	db *sql.DB
}

// NewStorage создает новый экземпляр Storage.
func NewStorage() *Storage {
	return &Storage{}
}

// Start открывает соединение с базой данных и применяет миграции.
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

// Stop закрывает соединение с базой данных.
func (s *Storage) Stop() error {
	return s.db.Close()
}
