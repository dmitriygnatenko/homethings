package thing

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings/internal/helpers/location"
	"github.com/dmitriygnatenko/homethings/internal/mappers"
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
func SearchThingHandler(thingRepository ThingRepository) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		search, _ := url.QueryUnescape(fctx.Params("search", ""))

		if match, _ := regexp.MatchString("^[A-Za-zА-Яа-я0-9 ]+$", search); !match {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		search = strings.TrimSpace(search)
		if len([]rune(search)) < 3 {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		res, err := thingRepository.Search(ctx, search)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		res = location.ApplyLocation(fctx, res)

		return fctx.JSON(mappers.ToThingsResponse(res))
	}
}
