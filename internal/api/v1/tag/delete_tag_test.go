package tag

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
	repoMocks "git.dmitriygnatenko.ru/dima/homethings/internal/repositories/mocks"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func Test_DeleteTagHandler(t *testing.T) {
	type tagRepoMockFunc func(mc *minimock.Controller) interfaces.TagRepository
	type thingTagRepoMockFunc func(mc *minimock.Controller) interfaces.ThingTagRepository

	type req struct {
		method      string
		route       string
		contentType string
	}

	var (
		mc        = minimock.NewController(t)
		tagID     = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())

		correctReq = req{
			method:      fiber.MethodDelete,
			route:       "/v1/tags/" + strconv.Itoa(tagID),
			contentType: fiber.MIMEApplicationJSON,
		}
	)

	tests := []struct {
		name             string
		req              req
		resCode          int
		resBody          interface{}
		tagRepoMock      tagRepoMockFunc
		thingTagRepoMock thingTagRepoMockFunc
	}{
		{
			name: "negative case - bad request",
			req: req{
				method:      fiber.MethodDelete,
				route:       "/v1/tags/" + gofakeit.Word(),
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				return repoMocks.NewTagRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - tag is not exist",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - tag repository error (get)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - tag repository error (begin tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				mock.BeginTxMock.Return(nil, testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing tag repository error (delete)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				mock.BeginTxMock.Return(nil, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)
				mock.DeleteByTagIDMock.Return(testError)
				return mock
			},
		},
		{
			name:    "negative case - tag repository error (delete)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				mock.BeginTxMock.Return(nil, nil)

				mock.DeleteMock.Return(testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)
				mock.DeleteByTagIDMock.Return(nil)
				return mock
			},
		},
		{
			name:    "negative case - tag repository error (commit tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				mock.BeginTxMock.Return(nil, nil)

				mock.DeleteMock.Return(nil)

				mock.CommitTxMock.Return(testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)
				mock.DeleteByTagIDMock.Return(nil)
				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				mock.BeginTxMock.Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, tagID, id)
				}).Return(nil)

				mock.CommitTxMock.Return(nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)

				mock.DeleteByTagIDMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, tagID, id)
				}).Return(nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.tagRepoMock(mc), tt.thingTagRepoMock(mc))

			fiberApp.Delete("/v1/tags/:tagId", DeleteTagHandler(serviceProvider))

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
