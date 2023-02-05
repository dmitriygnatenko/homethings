package models

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt string
	UpdatedAt string
}

type AddUserRequest struct {
	Username string
	Password string
}

type UpdateUserRequest struct {
	Username string
	Password string
}
