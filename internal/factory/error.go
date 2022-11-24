package factory

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateBadRequestResponse(fctx *fiber.Ctx, err error) error {
	if err == nil {
		return fctx.Status(fiber.StatusBadRequest).JSON(CreateEmptyResponse())
	}

	return fctx.Status(fiber.StatusBadRequest).JSON(CreateErrorResponse(err))
}

func CreateInternalErrorResponse(fctx *fiber.Ctx, err error) error {
	if err == nil {
		return fctx.Status(fiber.StatusInternalServerError).JSON(CreateEmptyResponse())
	}

	return fctx.Status(fiber.StatusInternalServerError).JSON(CreateErrorResponse(err))
}

func CreateNotFoundResponse(fctx *fiber.Ctx, err error) error {
	if err == nil {
		return fctx.Status(fiber.StatusNotFound).JSON(CreateEmptyResponse())
	}

	return fctx.Status(fiber.StatusNotFound).JSON(CreateErrorResponse(err))
}

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
