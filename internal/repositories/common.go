package repositories

import "github.com/lib/pq"

const (
	FKViolationErrorCode  = "23503"
	DuplicateKeyErrorCode = "23505"
)

func IsFKViolationError(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == FKViolationErrorCode
	}

	return false
}

func IsDuplicateKeyError(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == DuplicateKeyErrorCode
	}

	return false
}
