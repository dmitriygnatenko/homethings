package tag

import (
	"net/http/httptest"
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

func TestGetTagsHandler(t *testing.T) {
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
			route:  "/v1/tags/",
		}

		tagRepoRes = []models.Tag{
			{
				ID:        uint64(gofakeit.Number(1, 1000)),
				Title:     "A" + gofakeit.Phrase(),
				Style:     gofakeit.Phrase(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
			{
				ID:        uint64(gofakeit.Number(1, 1000)),
				Title:     "B" + gofakeit.Phrase(),
				Style:     gofakeit.Phrase(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}

		expectedRes = dto.TagsResponse{
			Tags: []dto.TagResponse{
				{
					ID:        tagRepoRes[0].ID,
					Title:     tagRepoRes[0].Title,
					Style:     tagRepoRes[0].Style,
					CreatedAt: tagRepoRes[0].CreatedAt.Format(layout),
					UpdatedAt: tagRepoRes[0].UpdatedAt.Format(layout),
				},
				{
					ID:        tagRepoRes[1].ID,
					Title:     tagRepoRes[1].Title,
					Style:     tagRepoRes[1].Style,
					CreatedAt: tagRepoRes[1].CreatedAt.Format(layout),
					UpdatedAt: tagRepoRes[1].UpdatedAt.Format(layout),
				},
			},
		}
	)

	tests := []struct {
		name        string
		req         req
		resCode     int
		resBody     interface{}
		tagRepoMock func(mc *minimock.Controller) TagRepository
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)
				mock.GetAllMock.Return(tagRepoRes, nil)
				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)
				mock.GetAllMock.Return(nil, testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Get("/v1/tags", GetTagsHandler(tt.tagRepoMock(mc)))

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
