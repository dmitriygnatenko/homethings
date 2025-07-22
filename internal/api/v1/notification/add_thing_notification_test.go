package notification

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/notification/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
	"github.com/dmitriygnatenko/homethings-v1/internal/repositories"
)

func TestAddThingNotificationHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        *dto.AddThingNotificationRequest
	}

	var (
		thingID          = uint64(gofakeit.Number(1, 1000))
		notificationDate = gofakeit.Date().Truncate(time.Second)
		testError        = gofakeit.Error()
		layout           = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodPost,
			route:  "/v1/things/notifications",
			body: &dto.AddThingNotificationRequest{
				ThingID:          thingID,
				NotificationDate: notificationDate.Format(time.RFC3339),
			},
			contentType: fiber.MIMEApplicationJSON,
		}

		repoRes = models.ThingNotification{
			ThingID:          thingID,
			NotificationDate: notificationDate,
			CreatedAt:        gofakeit.Date(),
			UpdatedAt:        gofakeit.Date(),
		}

		expectedRes = dto.ThingNotificationResponse{
			ThingID:          thingID,
			NotificationDate: notificationDate.Format(layout),
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

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingNotificationRequest) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(nil)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(&repoRes, nil)

				return mock
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPost,
				route:  "/v1/things/notifications",
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name: "negative case - thing id is empty",
			req: req{
				method: fiber.MethodPost,
				route:  "/v1/things/notifications",
				body: &dto.AddThingNotificationRequest{
					NotificationDate: notificationDate.Format(time.RFC3339),
				},
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name: "negative case - incorrect notification date format",
			req: req{
				method: fiber.MethodPost,
				route:  "/v1/things/notifications",
				body: &dto.AddThingNotificationRequest{
					ThingID:          thingID,
					NotificationDate: notificationDate.String(),
				},
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (add)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingNotificationRequest) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(testError)

				return mock
			},
		},
		{
			name:    "negative case - repository error (duplicate)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingNotificationRequest) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(&mysql.MySQLError{Number: repositories.DuplErrCode})

				return mock
			},
		},
		{
			name:    "negative case - repository error (get)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingNotificationRequest) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(nil)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, testError)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Post("/v1/things/notifications", AddThingNotificationHandler(tt.repoMock(mc)))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, test.ConvertDataToIOReader(tt.req.body))
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, test.TestTimeout)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)

			if tt.resBody != nil {
				assert.Equal(t, test.MarshalResponse(tt.resBody), test.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
