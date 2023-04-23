package timezone

import (
	"github.com/gofiber/fiber/v2"
)

func New(config ...Config) fiber.Handler {
	cfg := defaultConfig(config...)

	return func(c *fiber.Ctx) (err error) {
		headers := c.GetReqHeaders()

		if tz, ok := headers[cfg.HeaderName]; ok {
			c.Locals(CtxTimezoneKey, tz)
		}

		return c.Next()
	}
}
