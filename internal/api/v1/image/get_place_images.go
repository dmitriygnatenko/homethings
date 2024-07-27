package image

import (
	"sort"

	"github.com/gofiber/fiber/v2"

	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
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
		id, err := fctx.ParamsInt("placeId")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		placesRes, err := placeImageRepository.GetByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		res = append(res, placesRes...)

		thingsRes, err := thingImageRepository.GetByPlaceID(ctx, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		res = append(res, thingsRes...)

		sort.Slice(res, func(i, j int) bool {
			return res[i].CreatedAt.After(res[j].CreatedAt)
		})

		res = helpers.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToImagesResponse(res))
	}
}
