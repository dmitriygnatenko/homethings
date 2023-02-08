package auth

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/auth/check [get]
// @Success     200 {object} dto.LoginResponse
// @Failure     400 {object} dto.ErrorResponse
// @Summary     Check auth
// @Tags  		Auth
// @Accept      json
// @Produce     json
func CheckAuthHandler(_ interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
