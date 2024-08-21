package mappers

import (
	"database/sql"
	"strings"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ToUpdateUserRequest(id uint64, req dto.UpdateUserRequest) models.UpdateUserRequest {
	res := models.UpdateUserRequest{
		ID: id,
	}

	var username, password string
	if req.Username != nil {
		username = strings.TrimSpace(*req.Username)
	}
	if req.Password != nil {
		password = strings.TrimSpace(*req.Password)
	}

	if username != "" {
		res.Username = sql.NullString{String: username, Valid: true}
	}

	if password != "" {
		res.Password = sql.NullString{String: password, Valid: true}
	}

	return res
}

func ToUserResponse(req models.User) dto.UserResponse {
	return dto.UserResponse{Username: req.Username}
}

func ToLoginResponse(token string) dto.LoginResponse {
	return dto.LoginResponse{Token: token}
}
