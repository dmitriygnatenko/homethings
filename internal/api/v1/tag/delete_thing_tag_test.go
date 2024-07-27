package tag

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	API "git.dmitriygnatenko.ru/dima/homethings/internal/api/v1"
	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/tag/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func TestDeleteThingTagHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
	}

	var (
		tagID     = gofakeit.Number(1, 1000)
		thingID   = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method:      fiber.MethodDelete,
			route:       "/v1/tags/" + strconv.Itoa(tagID) + "/thing/" + strconv.Itoa(thingID),
			contentType: fiber.MIMEApplicationJSON,
		}
	)

	tests := []struct {
		name             string
		req              req
		resCode          int
		resBody          interface{}
		tagRepoMock      func(mc *minimock.Controller) TagRepository
		thingRepoMock    func(mc *minimock.Controller) ThingRepository
		thingTagRepoMock func(mc *minimock.Controller) ThingTagRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteMock.Inspect(func(ctx context.Context, req models.DeleteThingTagRequest, tx *sql.Tx) {
					assert.Equal(mc, tagID, req.TagID)
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(nil)

				return mock
			},
		},
		{
			name: "negative case - request without tagID",
			req: req{
				method:      fiber.MethodDelete,
				route:       "/v1/tags/" + gofakeit.Word() + "/thing/" + strconv.Itoa(thingID),
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				return mocks.NewTagRepositoryMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name: "negative case - request without thingID",
			req: req{
				method:      fiber.MethodDelete,
				route:       "/v1/tags/" + strconv.Itoa(tagID) + "/thing/" + gofakeit.Word(),
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				return mocks.NewTagRepositoryMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - tag is not exist",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - tag repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, testError)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				return mocks.NewThingRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing is not exist",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing tag repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteMock.Inspect(func(ctx context.Context, req models.DeleteThingTagRequest, tx *sql.Tx) {
					assert.Equal(mc, tagID, req.TagID)
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(testError)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()

			fiberApp.Delete("/v1/tags/:tagId/thing/:thingId", DeleteThingTagHandler(
				tt.tagRepoMock(mc),
				tt.thingRepoMock(mc),
				tt.thingTagRepoMock(mc),
			))

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
