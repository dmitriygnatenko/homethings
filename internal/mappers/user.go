package mappers

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func ConvertToAddUserRequestModel(username string, password string) models.AddUserRequest {
	return models.AddUserRequest{
		Username: username,
		Password: password,
	}
}

func ConvertToUpdateUserRequestModel(username string, password string) models.UpdateUserRequest {
	return models.UpdateUserRequest{
		Username: username,
		Password: password,
	}
}
