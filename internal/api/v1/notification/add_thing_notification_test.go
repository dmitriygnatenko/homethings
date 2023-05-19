package notification

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	"git.dmitriygnatenko.ru/dima/homethings/internal/repositories"
	repoMocks "git.dmitriygnatenko.ru/dima/homethings/internal/repositories/mocks"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_AddThingNotificationHandler(t *testing.T) {
	type repoMockFunc func(mc *minimock.Controller) interfaces.ThingNotificationRepository

	type req struct {
		method      string
		route       string
		body        *dto.AddThingNotificationRequest
		contentType string
	}

	var (
		mc               = minimock.NewController(t)
		thingID          = gofakeit.Number(1, 1000)
		notificationDate = gofakeit.Date().Truncate(time.Second)
		testError        = errors.New(gofakeit.Phrase())
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
		repoMock repoMockFunc
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingNotificationRequest, tx *sql.Tx) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(nil)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
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
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
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
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
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
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (add)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingNotificationRequest, tx *sql.Tx) {
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
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingNotificationRequest, tx *sql.Tx) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(&pq.Error{Code: repositories.DuplicateKeyErrorCode})

				return mock
			},
		},
		{
			name:    "negative case - repository error (get)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingNotificationRequest, tx *sql.Tx) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(nil)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, testError)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.repoMock(mc))

			fiberApp.Post("/v1/things/notifications", AddThingNotificationHandler(serviceProvider))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, helpers.ConvertDataToIOReader(tt.req.body))
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, API.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
