package mappers

import (
	"database/sql"
	"strings"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ConvertToUpdateUserRequestModel(id int, req dto.UpdateUserRequest) models.UpdateUserRequest {
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
		res.Username = sql.NullString{String: password, Valid: true}
	}

	return res
}
