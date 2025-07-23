package notification

import (
	"net/http/httptest"
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

func TestGetExpiredThingNotificationsHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/things/notifications/expired",
		}

		repoRes = []models.ExtThingNotification{
			{
				ThingID:          uint64(gofakeit.Number(1, 1000)),
				PlaceID:          uint64(gofakeit.Number(1, 1000)),
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
		repoMock func(mc *minimock.Controller) ThingNotificationRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)
				mock.GetExpiredMock.Return(repoRes, nil)
				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			repoMock: func(mc *minimock.Controller) ThingNotificationRepository {
				mock := mocks.NewThingNotificationRepositoryMock(mc)
				mock.GetExpiredMock.Return(nil, testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Get("/v1/things/notifications/expired", GetExpiredThingNotificationsHandler(tt.repoMock(mc)))

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
