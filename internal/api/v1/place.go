package v1

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/v1/places/{id} [get]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.PlaceResponse
// @Failure     404 {object} dto.EmptyResponse
// @Summary     Get one place by ID
// @Tags  		Places
// @security 	BasicAuth
// @Produce     json
func GetPlaceHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		return ctx.JSON(dto.PlaceResponse{})
	}
}

// @Router 		/v1/places [post]
// @Param       data body dto.AddPlaceRequest true "Request body"
// @Success     200 {object} dto.PlaceResponse
// @Summary     Add place
// @Tags  		Places
// @security 	BasicAuth
// @Produce     json
func AddPlaceHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := dto.UpdatePlaceRequest{}
		if err := ctx.BodyParser(&req); err != nil {
			return err
		}

		return ctx.JSON(dto.PlaceResponse{})
	}
}

// @Router 		/v1/places/{id} [put]
// @Param       id path int true "Place ID"
// @Param       data body dto.UpdatePlaceRequest true "Request body"
// @Success     200 {object} dto.PlaceResponse
// @Summary     Update place
// @Tags  		Places
// @security 	BasicAuth
// @Produce     json
func UpdatePlaceHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			return err
		}

		req := dto.UpdatePlaceRequest{}
		if err = ctx.BodyParser(&req); err != nil {
			return err
		}

		_ = id

		return ctx.JSON(dto.PlaceResponse{})
	}
}

// @Router 		/v1/places/{id} [delete]
// @Param       id path int true "Place ID"
// @Success     200 {object} dto.EmptyResponse
// @Summary     Delete place
// @Tags  		Places
// @security 	BasicAuth
// @Produce     json
func DeletePlaceHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id, err := ctx.ParamsInt("id")
		if err != nil {
			return err
		}

		_ = id

		return ctx.JSON(dto.EmptyResponse{})
	}
}
