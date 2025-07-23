package place

import (
	"database/sql"
	"errors"

	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings/internal/helpers/request"
	"github.com/dmitriygnatenko/homethings/internal/mappers"
)

// @Router 		/api/v1/places/{placeId} [get]
// @Param       placeId path int true "Place ID"
// @Success     200 {object} dto.PlaceResponse
// @Failure     404 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get one place
// @Tags  		Places
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetPlaceHandler(placeRepository PlaceRepository) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		id, err := request.ConvertToUint64(fctx, "placeId")

		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		res, err := placeRepository.Get(ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Info(ctx, err.Error())

				return fiber.NewError(fiber.StatusNotFound, "")
			}

			logger.Error(ctx, err.Error())

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToPlaceResponse(*res))
	}
}
