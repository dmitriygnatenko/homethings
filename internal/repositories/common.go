package repositories

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
)

var placeholder = sq.Question

const (
	DuplErrCode   = 1062
	FKViolErrCode = 1452
)

type DB interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func IsFKViolationError(err error) bool {
	var me *mysql.MySQLError

	ok := errors.As(err, &me)
	if !ok {
		return false
	}

	return me.Number == FKViolErrCode
}

func IsDuplicateKeyError(err error) bool {
	var me *mysql.MySQLError

	ok := errors.As(err, &me)
	if !ok {
		return false
	}

	return me.Number == DuplErrCode
}
