package user

import (
	"strings"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/users [post]
// @Param       data body dto.AddUserRequest true "Request body"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add user
// @Tags  		User
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddUserHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddUserRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		hash, err := sp.GetAuthService().GeneratePasswordHash(strings.TrimSpace(req.Password))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		_, err = sp.GetUserRepository().Add(ctx, strings.TrimSpace(req.Username), hash)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
