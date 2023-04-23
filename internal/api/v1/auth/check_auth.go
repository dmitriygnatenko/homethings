package auth

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/auth/check [get]
// @Success     200 {object} dto.UserResponse
// @Failure     403 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Check auth
// @Tags  		Auth
// @Accept      json
// @Produce     json
func CheckAuthHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		claims := sp.GetAuthService().GetClaims(fctx)

		user, err := sp.GetUserRepository().Get(ctx, claims["name"].(string))
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.NewError(fiber.StatusForbidden, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ConvertToUserResponseDTO(*user))
	}
}
