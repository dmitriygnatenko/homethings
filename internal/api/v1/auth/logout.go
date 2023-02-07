package auth

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/auth/logout [post]
// @Success     200 {object} dto.EmptyResponse
// @Success     403 {object} dto.EmptyResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Check username/password
// @Tags  		Auth
// @security 	BasicAuth
// @Accept      json
// @Produce     json
func LogoutHandler(_ interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
