package image

import (
	"sort"

	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/request"
	"github.com/dmitriygnatenko/homethings-v1/internal/mappers"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

// @Router 		/api/v1/images/place/{placeId} [get]
// @Param       placeId path int true "Place ID"
// @Success     200 {object} dto.ImagesResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get images by place ID (with child places)
// @Tags  		Images
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetPlaceImagesHandler(
	thingImageRepository ThingImageRepository,
	placeImageRepository PlaceImageRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		var res []models.Image

		ctx := fctx.Context()

		id, err := request.ConvertToUint64(fctx, "placeId")
		if err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		placesRes, err := placeImageRepository.GetByPlaceID(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = append(res, placesRes...)

		thingsRes, err := thingImageRepository.GetByPlaceID(ctx, id)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = append(res, thingsRes...)

		sort.Slice(res, func(i, j int) bool {
			return res[i].CreatedAt.After(res[j].CreatedAt)
		})

		return fctx.JSON(mappers.ToImagesResponse(location.ApplyLocation(fctx, res)))
	}
}
