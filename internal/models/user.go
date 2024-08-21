package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint64
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUserRequest struct {
	ID       uint64
	Username sql.NullString
	Password sql.NullString
}
