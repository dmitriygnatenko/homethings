package dto

type AddUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}
