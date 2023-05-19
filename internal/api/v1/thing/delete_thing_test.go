package thing

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

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

func Test_DeleteThingHandler(t *testing.T) {
	type req struct {
		method string
		route  string
	}

	var (
		mc        = minimock.NewController(t)
		thingID   = gofakeit.Number(1, 1000)
		imageURL  = gofakeit.URL()
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodDelete,
			route:  "/v1/things/" + strconv.Itoa(thingID),
		}

		repoImagesRes = []models.Image{
			{
				Image: imageURL,
			},
		}
	)

	tests := []struct {
		name                      string
		req                       req
		resCode                   int
		resBody                   interface{}
		thingRepoMock             func(mc *minimock.Controller) interfaces.ThingRepository
		placeThingRepoMock        func(mc *minimock.Controller) interfaces.PlaceThingRepository
		thingImageRepoMock        func(mc *minimock.Controller) interfaces.ThingImageRepository
		thingTagRepoMock          func(mc *minimock.Controller) interfaces.ThingTagRepository
		thingNotificationRepoMock func(mc *minimock.Controller) interfaces.ThingNotificationRepository
		fileRepoMock              func(mc *minimock.Controller) interfaces.FileRepository
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodDelete,
				route:  "/v1/things/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - bad request (thing not found)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (begin tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete place thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(testError)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(nil, testError)
				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(testError)
				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete thing tags)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)
				mock.DeleteByThingIDMock.Return(testError)
				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				return repoMocks.NewThingNotificationRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)
				mock.DeleteByThingIDMock.Return(nil)
				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (commit tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				mock.CommitTxMock.Return(testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)
				mock.DeleteByThingIDMock.Return(nil)
				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - file delete error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				mock.CommitTxMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)
				mock.DeleteByThingIDMock.Return(nil)
				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				mock := repoMocks.NewFileRepositoryMock(mc)
				mock.DeleteMock.Return(testError)
				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				mock.CommitTxMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(repoImagesRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)
				mock.DeleteByThingIDMock.Return(nil)
				return mock
			},
			thingNotificationRepoMock: func(mc *minimock.Controller) interfaces.ThingNotificationRepository {
				mock := repoMocks.NewThingNotificationRepositoryMock(mc)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				mock := repoMocks.NewFileRepositoryMock(mc)
				mock.DeleteMock.Return(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(
				tt.thingRepoMock(mc),
				tt.placeThingRepoMock(mc),
				tt.thingImageRepoMock(mc),
				tt.thingTagRepoMock(mc),
				tt.thingNotificationRepoMock(mc),
				tt.fileRepoMock(mc),
			)

			fiberApp.Delete("/v1/things/:thingId", DeleteThingHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
