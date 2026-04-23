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

func (s *Storage) CreateUser(ctx context.Context, login, passwordHash string) (*model.User, error) {
	db := s.dbFromCtx(ctx)

	row := db.QueryRowContext(ctx,
		`INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id, login, password_hash`,
		login, passwordHash,
	)

	var u model.User
	err := row.Scan(&u.ID, &u.Login, &u.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, ErrLoginConflict
		}
		return nil, fmt.Errorf("postgres.CreateUser Scan: %w", err)
	}

	return &u, nil
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	db := s.dbFromCtx(ctx)

	row := db.QueryRowContext(ctx,
		`SELECT id, login, password_hash FROM users WHERE login = $1`,
		login,
	)

	var u model.User
	err := row.Scan(&u.ID, &u.Login, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("postgres.GetUserByLogin Scan: %w", err)
	}

	return &u, nil
}
