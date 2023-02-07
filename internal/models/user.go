package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUserRequest struct {
	ID       int
	Username sql.NullString
	Password sql.NullString
}
