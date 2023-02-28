package tag

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/tags/{id} [get]
// @Param       id path int true "Tag ID"
// @Success     200 {object} dto.TagResponse
// @Failure     404 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get one tag by ID
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetTagHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		return fctx.JSON(nil)
	}
}
