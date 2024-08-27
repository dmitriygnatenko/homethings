package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

const (
	FKViolationErrorCode  = "23503"
	DuplicateKeyErrorCode = "23505"
)

type DB interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func IsFKViolationError(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == FKViolationErrorCode
	}

	return false
}

func IsDuplicateKeyError(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == DuplicateKeyErrorCode
	}

	return false
}
