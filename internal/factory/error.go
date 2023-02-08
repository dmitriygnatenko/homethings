package factory

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"github.com/go-playground/validator/v10"
)

func CreateValidateErrorResponse(errors error) []dto.ValidateErrorResponse {
	var res []dto.ValidateErrorResponse //nolint

	for _, err := range errors.(validator.ValidationErrors) {
		res = append(res, dto.ValidateErrorResponse{
			Field: err.StructNamespace(),
			Tag:   err.Tag(),
			Value: err.Param(),
		})
	}

	return res
}

func CreateErrorResponse(err error) dto.ErrorResponse {
	return dto.ErrorResponse{Error: err.Error()}
}
