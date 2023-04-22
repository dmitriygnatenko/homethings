package notification

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/things/notification [post]
// @Param       data body dto.AddThingNotificationRequest true "Request body"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add thing notification
// @Tags  		Notifications
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddThingNotificationHandler(sp interfaces.ServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		//ctx := fctx.Context()

		req := dto.AddThingNotificationRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		//id, err := sp.GetPlaceRepository().Add(ctx, mappers.ConvertToAddPlaceRequestModel(req), nil)
		//if err != nil {
		//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		//}
		//
		//res, err := sp.GetPlaceRepository().Get(ctx, id)
		//if err != nil {
		//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		//}

		return fctx.JSON(dto.EmptyResponse{})
	}
}
