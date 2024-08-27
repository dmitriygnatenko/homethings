package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint64    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UpdateUserRequest struct {
	ID       uint64
	Username sql.NullString
	Password sql.NullString
}
