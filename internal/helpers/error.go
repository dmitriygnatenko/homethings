package helpers

import (
	"encoding/json"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func FormatError(errors error) error {
	var res []dto.ErrorResponse

	for _, err := range errors.(validator.ValidationErrors) {
		res = append(res, dto.ErrorResponse{
			Field: err.StructNamespace(),
			Tag:   err.Tag(),
			Value: err.Param(),
		})
	}

	errMsg, err := json.Marshal(res)
	if err != nil {
		return err
	}

	return fiber.NewError(fiber.StatusBadRequest, string(errMsg))
}
