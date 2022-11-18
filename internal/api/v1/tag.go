package v1

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

func GetTagsHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		res, err := sp.GetTagRepository().GetAll(ctx.Context())
		if err != nil {
			return err
		}

		return ctx.JSON(mappers.ConvertToTagsResponseDTO(res))
	}
}
