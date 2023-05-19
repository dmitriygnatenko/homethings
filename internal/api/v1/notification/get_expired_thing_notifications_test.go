package notification

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	repoMocks "git.dmitriygnatenko.ru/dima/homethings/internal/repositories/mocks"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func Test_GetExpiredThingNotificationsHandler(t *testing.T) {
	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		testError = errors.New(gofakeit.Phrase())
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/things/notifications/expired",
		}

		repoRes = []models.ExtThingNotification{
			{
				ThingID:          gofakeit.Number(1, 1000),
				PlaceID:          gofakeit.Number(1, 1000),
				ThingTitle:       gofakeit.Phrase(),
				PlaceTitle:       gofakeit.Phrase(),
				NotificationDate: gofakeit.Date().Truncate(time.Second),
				CreatedAt:        gofakeit.Date(),
				UpdatedAt:        gofakeit.Date(),
			},
		}

		expectedRes = dto.ThingNotificationsExtResponse{
			Notifications: []dto.ThingNotificationExtResponse{
				{
					ThingID:          repoRes[0].ThingID,
					PlaceID:          repoRes[0].PlaceID,
					ThingTitle:       repoRes[0].ThingTitle,
					PlaceTitle:       repoRes[0].PlaceTitle,
					NotificationDate: repoRes[0].NotificationDate.Format(layout),
					CreatedAt:        repoRes[0].CreatedAt.Format(layout),
					UpdatedAt:        repoRes[0].UpdatedAt.Format(layout),
				},
			},
		}
	)

	tests := []struct {
		name     string
		req      req
		resCode  int
		resBody  interface{}
		repoMock func(mc *minimock.Controller) interfaces.ThingNotificationRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)
				mock.GetExpiredMock.Return(repoRes, nil)
				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)
				mock.GetExpiredMock.Return(nil, testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.repoMock(mc))

			fiberApp.Get("/v1/things/notifications/expired", GetExpiredThingNotificationsHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
