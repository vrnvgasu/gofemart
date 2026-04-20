package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrLoginConflict       = errors.New("login already taken")
	ErrOrderAlreadyExists  = errors.New("order already submitted by this user")
	ErrOrderConflict       = errors.New("order submitted by another user")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type pgErrorClassification int

const (
	nonRetriable pgErrorClassification = iota
	retriable
)

type postgresErrorClassifier struct{}

func newPostgresErrorClassifier() *postgresErrorClassifier {
	return &postgresErrorClassifier{}
}

func (c *postgresErrorClassifier) classify(err error) pgErrorClassification {
	if err == nil {
		return nonRetriable
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return classifyPgError(pgErr)
	}

	return nonRetriable
}

func classifyPgError(pgErr *pgconn.PgError) pgErrorClassification {
	switch pgErr.Code {
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure,
		pgerrcode.TransactionRollback,
		pgerrcode.SerializationFailure,
		pgerrcode.DeadlockDetected,
		pgerrcode.TooManyConnections,
		pgerrcode.CannotConnectNow,
		pgerrcode.AdminShutdown:
		return retriable
	default:
		return nonRetriable
	}
}
