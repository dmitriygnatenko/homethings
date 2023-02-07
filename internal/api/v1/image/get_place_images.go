package image

import (
	"sort"

	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/images/place/{id} [get]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.ImagesResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Get images by place ID (with child places)
// @Tags  		Images
// @security 	APIKey
// @Accept      json
// @Produce     json
func GetPlaceImagesHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		var res []models.Image
		ctx := fctx.Context()
		id, err := fctx.ParamsInt("id")
		if err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		placesRes, err := sp.GetPlaceImageRepository().GetByPlaceID(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}
		res = append(res, placesRes...)

		thingsRes, err := sp.GetThingImageRepository().GetByPlaceID(ctx, id)
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}
		res = append(res, thingsRes...)

		sort.Slice(res, func(i, j int) bool {
			return res[i].CreatedAt.After(res[j].CreatedAt)
		})

		return fctx.JSON(mappers.ConvertToImagesResponseDTO(res))
	}
}
