package tag

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"git.dmitriygnatenko.ru/dima/homethings/internal/api/v1/tag/mocks"
	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/helpers/test"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

func TestAddThingTagHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
	}

	var (
		tagID     = uint64(gofakeit.Number(1, 1000))
		thingID   = uint64(gofakeit.Number(1, 1000))
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodPost,
			route: "/v1/tags/" + strconv.FormatUint(tagID, 10) +
				"/thing/" + strconv.FormatUint(thingID, 10),
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
		tagRepoMock      func(mc *minimock.Controller) TagRepository
		thingRepoMock    func(mc *minimock.Controller) ThingRepository
		thingTagRepoMock func(mc *minimock.Controller) ThingTagRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(&tagRepoRes, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingTagRequest) {
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
				route:       "/v1/tags/" + gofakeit.Word() + "/thing/" + strconv.FormatUint(thingID, 10),
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
				method:      fiber.MethodPost,
				route:       "/v1/tags/" + strconv.FormatUint(tagID, 10) + "/thing/" + gofakeit.Word(),
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

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
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

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
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

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
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

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
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

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) ThingRepository {
				mock := mocks.NewThingRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, thingID, id)
				}).Return(nil, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.AddMock.Inspect(func(ctx context.Context, req models.AddThingTagRequest) {
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

			fiberApp.Post("/v1/tags/:tagId/thing/:thingId", AddThingTagHandler(
				tt.tagRepoMock(mc),
				tt.thingRepoMock(mc),
				tt.thingTagRepoMock(mc),
			))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, nil)
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, test.TestTimeout)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, test.MarshalResponse(tt.resBody), test.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
