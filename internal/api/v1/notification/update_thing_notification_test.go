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

	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/notification/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/test"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func TestUpdateThingNotificationHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        *dto.UpdateThingNotificationRequest
	}

	var (
		thingID          = uint64(gofakeit.Number(1, 1000))
		notificationDate = gofakeit.Date().Truncate(time.Second)
		testError        = gofakeit.Error()
		layout           = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodPut,
			route:  "/v1/things/notifications/" + strconv.FormatUint(thingID, 10),
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
		repoMock func(mc *minimock.Controller) ThingNotificationRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingNotificationRequest) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(nil)

				mock.GetMock.Set(func(ctx context.Context, id uint64) (*models.ThingNotification, error) {
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
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/things/notifications/" + strconv.FormatUint(thingID, 10),
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name: "negative case - validate request error",
			req: req{
				method:      fiber.MethodPut,
				route:       "/v1/things/notifications/" + strconv.FormatUint(thingID, 10),
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				return mocks.NewThingNotificationRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get)",
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
			name:    "negative case - bad request (notification not found)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
		},
		{
			name: "negative case - bad request (notification not found)",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/things/notifications/" + strconv.FormatUint(thingID, 10),
				body: &dto.UpdateThingNotificationRequest{
					NotificationDate: notificationDate.String(),
				},
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
		},
		{
			name:    "negative case - repository error (update)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingNotificationRequest) {
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
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateThingNotificationRequest) {
					assert.Equal(mc, thingID, req.ThingID)
					assert.Equal(mc, notificationDate, req.NotificationDate)
				}).Return(nil)

				mock.GetMock.Set(func(ctx context.Context, id uint64) (*models.ThingNotification, error) {
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
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Put("/v1/things/notifications/:thingId", UpdateThingNotificationHandler(tt.repoMock(mc)))

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
