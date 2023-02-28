package tag

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/tags/{id} [put]
// @Param       id path int true "Tag ID"
// @Param       data body dto.UpdateTagRequest true "Request body"
// @Success     200 {object} dto.TagResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update tag
// @Tags  		Tags
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdateTagHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		return fctx.JSON(nil)
	}
}
