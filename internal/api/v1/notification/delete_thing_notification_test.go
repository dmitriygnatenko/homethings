package notification

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
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

func Test_DeleteThingNotificationHandler(t *testing.T) {
	type repoMockFunc func(mc *minimock.Controller) interfaces.ThingNotificationRepository

	type req struct {
		method      string
		route       string
		contentType string
	}

	var (
		mc        = minimock.NewController(t)
		thingID   = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method:      fiber.MethodDelete,
			route:       "/v1/things/notifications/" + strconv.Itoa(thingID),
			contentType: fiber.MIMEApplicationJSON,
		}

		repoRes = models.ThingNotification{
			ThingID:          thingID,
			NotificationDate: gofakeit.Date().Truncate(time.Second),
			CreatedAt:        gofakeit.Date(),
			UpdatedAt:        gofakeit.Date(),
		}
	)

	tests := []struct {
		name     string
		req      req
		resCode  int
		resBody  interface{}
		repoMock repoMockFunc
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(&repoRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method:      fiber.MethodDelete,
				route:       "/v1/things/notifications/" + gofakeit.Word(),
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get)",
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
			name:    "negative case - bad request (notification not found)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
		},
		{
			name:    "negative case - repository error (delete)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(&repoRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, thingID, id)
				}).Return(testError)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.repoMock(mc))

			fiberApp.Delete("/v1/things/notifications/:thingId", DeleteThingNotificationHandler(serviceProvider))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, nil)
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, API.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
