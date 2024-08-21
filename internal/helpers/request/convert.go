package request

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ConvertToUint64(fctx *fiber.Ctx, key string) (uint64, error) {
	val, err := fctx.ParamsInt(key)
	if err != nil {
		return 0, err
	}

	if val < 0 {
		return 0, errors.New("value must be positive")
	}

	return uint64(val), nil
}
