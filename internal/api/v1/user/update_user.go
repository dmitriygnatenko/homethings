package user

import (
	"database/sql"
	"strings"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/users [put]
// @Param       data body dto.UpdateUserRequest true "Request body"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update user
// @Tags  		User
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdateUserHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		var err error
		var username, password string

		ctx := fctx.Context()
		req := dto.UpdateUserRequest{}
		if err = fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if req.Username != nil {
			username = strings.TrimSpace(*req.Username)
		}

		if req.Password != nil {
			password, err = sp.GetAuthService().GeneratePasswordHash(strings.TrimSpace(*req.Password))
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		if username == "" && password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		claims := sp.GetAuthService().GetClaims(fctx)

		user, err := sp.GetUserRepository().Get(ctx, claims["name"].(string))
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		req.Password = &password
		req.Username = &username

		err = sp.GetUserRepository().Update(ctx, mappers.ConvertToUpdateUserRequestModel(user.ID, req))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
