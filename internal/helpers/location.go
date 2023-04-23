package helpers

import (
	"time"

	"git.dmitriygnatenko.ru/dima/homethings/internal/middleware/timezone"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ApplyLocation[T interface{}](fctx *fiber.Ctx, req T) T {
	tz, ok := fctx.Locals(timezone.CtxTimezoneKey).(string)
	if !ok {
		return req
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return req
	}

	switch v := any(req).(type) {
	case []models.Place:
		for i := range v {
			v[i].CreatedAt = v[i].CreatedAt.In(loc)
			v[i].UpdatedAt = v[i].UpdatedAt.In(loc)
		}
	case []models.Thing:
		for i := range v {
			v[i].CreatedAt = v[i].CreatedAt.In(loc)
			v[i].UpdatedAt = v[i].UpdatedAt.In(loc)
		}
	case []models.Tag:
		for i := range v {
			v[i].CreatedAt = v[i].CreatedAt.In(loc)
			v[i].UpdatedAt = v[i].UpdatedAt.In(loc)
		}
	case []models.Image:
		for i := range v {
			v[i].CreatedAt = v[i].CreatedAt.In(loc)
		}
	case []models.ExtThingNotification:
		for i := range v {
			v[i].CreatedAt = v[i].CreatedAt.In(loc)
			v[i].UpdatedAt = v[i].UpdatedAt.In(loc)
			v[i].NotificationDate = v[i].NotificationDate.In(loc)
		}

	case *models.Thing:
		v.CreatedAt = v.CreatedAt.In(loc)
		v.UpdatedAt = v.UpdatedAt.In(loc)
	case *models.Place:
		v.CreatedAt = v.CreatedAt.In(loc)
		v.UpdatedAt = v.UpdatedAt.In(loc)
	case *models.ThingTag:
		v.CreatedAt = v.CreatedAt.In(loc)
		v.UpdatedAt = v.UpdatedAt.In(loc)
	case *models.Tag:
		v.CreatedAt = v.CreatedAt.In(loc)
		v.UpdatedAt = v.UpdatedAt.In(loc)
	case *models.Image:
		v.CreatedAt = v.CreatedAt.In(loc)
	case *models.ThingNotification:
		v.CreatedAt = v.CreatedAt.In(loc)
		v.UpdatedAt = v.UpdatedAt.In(loc)
		v.NotificationDate = v.NotificationDate.In(loc)
	}

	return req
}
