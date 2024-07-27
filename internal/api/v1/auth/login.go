package auth

import (
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
)

// @Router 		/api/v1/auth/login [post]
// @Param       data body dto.LoginRequest true "Request body"
// @Success     200 {object} dto.LoginResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Log in user
// @Tags  		Auth
// @Accept      json
// @Produce     json
func LoginHandler(
	authService AuthService,
	userRepository UserRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.LoginRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		user, err := userRepository.Get(ctx, req.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusForbidden, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if !authService.IsCorrectPassword(req.Password, user.Password) {
			return fiber.NewError(fiber.StatusForbidden, "")
		}

		token, err := authService.GenerateToken(*user)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ToLoginResponse(token))
	}
}
