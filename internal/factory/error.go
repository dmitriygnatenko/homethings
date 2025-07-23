package factory

import (
	"github.com/go-playground/validator/v10"

	"github.com/dmitriygnatenko/homethings/internal/dto"
)

// nolint:errorlint,prealloc
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
