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

func Test_UpdateThingNotificationHandler(t *testing.T) {
	type repoMockFunc func(mc *minimock.Controller) interfaces.ThingNotificationRepository

	type req struct {
		method      string
		route       string
		body        *dto.UpdateThingNotificationRequest
		contentType string
	}

	var (
		mc               = minimock.NewController(t)
		thingID          = gofakeit.Number(1, 1000)
		notificationDate = gofakeit.Date().Truncate(time.Second)
		testError        = errors.New(gofakeit.Phrase())
		layout           = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodPut,
			route:  "/v1/things/notifications/" + strconv.Itoa(thingID),
			body: &dto.UpdateThingNotificationRequest{
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

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingNotificationRequest, tx *sql.Tx) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(nil)

				mock.GetMock.Set(func(ctx context.Context, id int) (*models.ThingNotification, error) {
					assert.Equal(mc, thingID, id)
					if mock.GetAfterCounter() == 0 {
						return nil, nil
					}
					return &repoRes, nil
				})

				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/things/notifications/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/things/notifications/" + strconv.Itoa(thingID),
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name: "negative case - validate request error",
			req: req{
				method:      fiber.MethodPut,
				route:       "/v1/things/notifications/" + strconv.Itoa(thingID),
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
			name: "negative case - bad request (notification not found)",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/things/notifications/" + strconv.Itoa(thingID),
				body: &dto.UpdateThingNotificationRequest{
					NotificationDate: notificationDate.String(),
				},
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
		},
		{
			name:    "negative case - repository error (update)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingNotificationRequest, tx *sql.Tx) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(testError)

				return mock
			},
		},
		{
			name:    "negative case - repository error (get)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingNotificationRequest, tx *sql.Tx) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(nil)

				mock.GetMock.Set(func(ctx context.Context, id int) (*models.ThingNotification, error) {
					assert.Equal(mc, thingID, id)
					if mock.GetAfterCounter() == 0 {
						return nil, nil
					}
					return nil, testError
				})

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.repoMock(mc))

			fiberApp.Put("/v1/things/notifications/:thingId", UpdateThingNotificationHandler(serviceProvider))

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
