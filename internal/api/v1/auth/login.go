package auth

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/auth/login [post]
// @Param       data body dto.LoginRequest true "Request body"
// @Success     200 {object} dto.LoginResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Login user
// @Tags  		Auth
// @Accept      json
// @Produce     json
func LoginHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.LoginRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		user, err := sp.GetUserRepository().Get(ctx, req.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return factory.CreateBadRequestResponse(fctx, nil)
			}

			return factory.CreateInternalErrorResponse(fctx, err)
		}

		if !sp.GetAuthService().IsCorrectPassword(req.Password, user.Password) {
			return factory.CreateBadRequestResponse(fctx, nil)
		}

		token, err := sp.GetAuthService().GenerateToken(*user)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(dto.LoginResponse{Token: token})
	}
}
