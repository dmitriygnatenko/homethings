package models

import "time"

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AddUserRequest struct {
	Username string
	Password string
}

type UpdateUserRequest struct {
	Username string
	Password string
}
