package notification

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
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
	"net/http/httptest"
)

func Test_GetThingNotificationHandler(t *testing.T) {
	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		thingID   = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/things/notifications/" + strconv.Itoa(thingID),
		}

		repoRes = models.ThingNotification{
			ThingID:          thingID,
			NotificationDate: gofakeit.Date().Truncate(time.Second),
			CreatedAt:        gofakeit.Date(),
			UpdatedAt:        gofakeit.Date(),
		}

		expectedRes = dto.ThingNotificationResponse{
			ThingID:          thingID,
			NotificationDate: repoRes.NotificationDate.Format(layout),
			CreatedAt:        repoRes.CreatedAt.Format(layout),
			UpdatedAt:        repoRes.UpdatedAt.Format(layout),
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

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(&repoRes, nil)

				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, testError)

				return mock
			},
		},
		{
			name:    "negative case - not found",
			req:     correctReq,
			resCode: fiber.StatusNotFound,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/things/notifications/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.repoMock(mc))

			fiberApp.Get("/v1/things/notifications/:thingId", GetThingNotificationHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
