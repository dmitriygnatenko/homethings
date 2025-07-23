package notification

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/dmitriygnatenko/homethings/internal/api/v1/notification/mocks"
	"github.com/dmitriygnatenko/homethings/internal/dto"
	"github.com/dmitriygnatenko/homethings/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings/internal/models"
)

func TestGetThingNotificationHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		thingID   = uint64(gofakeit.Number(1, 1000))
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/things/notifications/" + strconv.FormatUint(thingID, 10),
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
		repoMock func(mc *minimock.Controller) ThingNotificationRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(&repoRes, nil)

				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, testError)

				return mock
			},
		},
		{
			name:    "negative case - not found",
			req:     correctReq,
			resCode: fiber.StatusNotFound,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
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
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Get("/v1/things/notifications/:thingId", GetThingNotificationHandler(tt.repoMock(mc)))

			fiberRes, _ := fiberApp.Test(
				httptest.NewRequest(tt.req.method, tt.req.route, nil),
				test.TestTimeout,
			)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)

			if tt.resBody != nil {
				assert.Equal(t, test.MarshalResponse(tt.resBody), test.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
