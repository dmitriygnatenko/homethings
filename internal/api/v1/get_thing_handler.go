package v1

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/things/{id} [get]
// @Param       id path int true "Thing ID"
// @Success     200 {object} dto.ThingResponse
// @Failure     404 {object} dto.EmptyResponse
// @Summary     Get one thing by ID
// @Tags  		Things
// @security 	BasicAuth
// @Produce     json
func GetThingHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		res, err := sp.GetThingRepository().Get(ctx, id)
		if err != nil {
			if err == sql.ErrNoRows {
				return factory.CreateNotFoundResponse(fctx, nil)
			}

			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(mappers.ConvertToThingResponseDTO(*res))
	}
}
