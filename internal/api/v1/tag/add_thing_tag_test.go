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
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	repoMocks "git.dmitriygnatenko.ru/dima/homethings/internal/repositories/mocks"
	sp "git.dmitriygnatenko.ru/dima/homethings/internal/service_provider"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func Test_AddThingTagHandler(t *testing.T) {
	type tagRepoMockFunc func(mc *minimock.Controller) interfaces.TagRepository
	type thingRepoMockFunc func(mc *minimock.Controller) interfaces.ThingRepository
	type thingTagRepoMockFunc func(mc *minimock.Controller) interfaces.ThingTagRepository

	type req struct {
		method      string
		route       string
		contentType string
	}

	var (
		mc        = minimock.NewController(t)
		tagID     = gofakeit.Number(1, 1000)
		thingID   = gofakeit.Number(1, 1000)
		testError = errors.New(gofakeit.Phrase())
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method:      fiber.MethodPost,
			route:       "/v1/tags/" + strconv.Itoa(tagID) + "/thing/" + strconv.Itoa(thingID),
			contentType: fiber.MIMEApplicationJSON,
		}

		tagRepoRes = models.Tag{
			ID:        tagID,
			Title:     gofakeit.Phrase(),
			Style:     gofakeit.Phrase(),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}

		expectedRes = dto.TagResponse{
			ID:        tagRepoRes.ID,
			Title:     tagRepoRes.Title,
			Style:     tagRepoRes.Style,
			CreatedAt: tagRepoRes.CreatedAt.Format(layout),
			UpdatedAt: tagRepoRes.UpdatedAt.Format(layout),
		}
	)

	tests := []struct {
		name             string
		req              req
		resCode          int
		resBody          interface{}
		tagRepoMock      tagRepoMockFunc
		thingRepoMock    thingRepoMockFunc
		thingTagRepoMock thingTagRepoMockFunc
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(&tagRepoRes, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingTagRequest, tx *sql.Tx) {
					assert.Equal(mc, tagID, req.TagID)
					assert.Equal(mc, thingID, req.ThingID)
				}).Return(nil)

				return mock
			},
		},
		{
			name: "negative case - request without tagID",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/tags/" + gofakeit.Word() + "/thing/" + strconv.Itoa(thingID),
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				return repoMocks.NewTagRepositoryMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name: "negative case - request without thingID",
			req: req{
				method:      fiber.MethodPost,
				route:       "/v1/tags/" + strconv.Itoa(tagID) + "/thing/" + gofakeit.Word(),
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				return repoMocks.NewTagRepositoryMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
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
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - tag repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, testError)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing is not exist",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				return repoMocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing tag repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) interfaces.ThingTagRepository {
				mock := repoMocks.NewThingTagRepositoryMock(mc)
				mock.AddMock.Return(testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.tagRepoMock(mc), tt.thingRepoMock(mc), tt.thingTagRepoMock(mc))

			fiberApp.Post("/v1/tags/:tagId/thing/:thingId", AddThingTagHandler(serviceProvider))

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
