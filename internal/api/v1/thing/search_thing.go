package thing

import (
	"net/url"
	"regexp"
	"strings"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
)

// @Router 		/api/v1/things/search/{search} [get]
// @Param       search path string true "Search string"
// @Success     200 {object} dto.ThingResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Search things
// @Tags  		Things
// @security 	APIKey
// @Accept      json
// @Produce     json
func SearchThingHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		search, err := url.QueryUnescape(fctx.Params("search", ""))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if match, _ := regexp.MatchString("^[A-Za-zА-Яа-я0-9 ]+$", search); !match {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		search = strings.TrimSpace(search)
		if len([]rune(search)) < 3 {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		res, err := sp.GetThingRepository().Search(ctx, search)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(mappers.ConvertToThingsResponseDTO(res))
	}
}
