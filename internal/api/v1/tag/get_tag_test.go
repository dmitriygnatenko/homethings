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

	"github.com/dmitriygnatenko/homethings/internal/api/v1/tag/mocks"
	"github.com/dmitriygnatenko/homethings/internal/dto"
	"github.com/dmitriygnatenko/homethings/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings/internal/models"
)

func TestGetTagHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method string
		route  string
	}

	var (
		tagID     = uint64(gofakeit.Number(1, 1000))
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodGet,
			route:  "/v1/tags/" + strconv.FormatUint(tagID, 10),
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

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(&tagRepoRes, nil)

				return mock
			},
		},
		{
			name:    "negative case - not found",
			req:     correctReq,
			resCode: fiber.StatusNotFound,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
		},
		{
			name:    "negative case - repository error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, testError)

				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodGet,
				route:  "/v1/tags/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			resBody: nil,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				return mocks.NewTagRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Get("/v1/tags/:tagId", GetTagHandler(tt.tagRepoMock(mc)))

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
