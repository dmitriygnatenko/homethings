package factory

import (
	"github.com/go-playground/validator/v10"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
)

func CreateValidateErrorResponse(errors error) []dto.ValidateErrorResponse {
	var res []dto.ValidateErrorResponse

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
