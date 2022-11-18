package helpers

import (
	"encoding/json"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func FormatValidateErrors(errors error) []dto.ErrorResponse {
	var res []dto.ErrorResponse

	for _, err := range errors.(validator.ValidationErrors) {
		res = append(res, dto.ErrorResponse{
			Field: err.StructNamespace(),
			Tag:   err.Tag(),
			Value: err.Param(),
		})
	}

	return res
}

func BadRequestJsonResponse(body interface{}) error {
	errMsg, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return fiber.NewError(fiber.StatusBadRequest, string(errMsg))
}

func NotFoundJsonResponse(body interface{}) error {
	errMsg, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return fiber.NewError(fiber.StatusNotFound, string(errMsg))
}
