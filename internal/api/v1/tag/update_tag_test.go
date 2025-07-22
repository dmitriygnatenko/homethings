package tag

import (
	"context"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/tag/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
	"github.com/dmitriygnatenko/homethings-v1/internal/models"
)

func TestUpdateTagHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
		body        *dto.UpdateTagRequest
	}

	var (
		tagID     = uint64(gofakeit.Number(1, 1000))
		title     = gofakeit.Phrase()
		style     = gofakeit.Phrase()
		testError = gofakeit.Error()
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodPut,
			route:  "/v1/tags/" + strconv.FormatUint(tagID, 10),
			body: &dto.UpdateTagRequest{
				Title: title,
				Style: style,
			},
			contentType: fiber.MIMEApplicationJSON,
		}

		repoRes = models.Tag{
			ID:        tagID,
			Title:     title,
			Style:     style,
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}

		expectedRes = dto.TagResponse{
			ID:        repoRes.ID,
			Title:     repoRes.Title,
			Style:     repoRes.Style,
			CreatedAt: repoRes.CreatedAt.Format(layout),
			UpdatedAt: repoRes.UpdatedAt.Format(layout),
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

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateTagRequest) {
					assert.Equal(mc, tagID, req.ID)
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, style, req.Style)
				}).Return(nil)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(&repoRes, nil)

				return mock
			},
		},
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/tags/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				return mocks.NewTagRepositoryMock(mc)
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/tags/" + strconv.FormatUint(tagID, 10),
			},
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				return mocks.NewTagRepositoryMock(mc)
			},
		},
		{
			name: "negative case - request without title",
			req: req{
				method:      fiber.MethodPut,
				route:       "/v1/tags/" + strconv.FormatUint(tagID, 10),
				contentType: fiber.MIMEApplicationJSON,
				body: &dto.UpdateTagRequest{
					Style: style,
				},
			},
			resCode: fiber.StatusBadRequest,
			resBody: []*dto.ValidateErrorResponse{
				{
					Field: "UpdateTagRequest.Title",
					Tag:   "required",
				},
			},
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				return mocks.NewTagRepositoryMock(mc)
			},
		},
		{
			name: "negative case - request without style",
			req: req{
				method:      fiber.MethodPut,
				route:       "/v1/tags/" + strconv.FormatUint(tagID, 10),
				contentType: fiber.MIMEApplicationJSON,
				body: &dto.UpdateTagRequest{
					Title: title,
				},
			},
			resCode: fiber.StatusBadRequest,
			resBody: []*dto.ValidateErrorResponse{
				{
					Field: "UpdateTagRequest.Style",
					Tag:   "required",
				},
			},
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				return mocks.NewTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (update)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateTagRequest) {
					assert.Equal(mc, tagID, req.ID)
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, style, req.Style)
				}).Return(testError)

				return mock
			},
		},
		{
			name:    "negative case - repository error (get)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateTagRequest) {
					assert.Equal(mc, tagID, req.ID)
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, style, req.Style)
				}).Return(nil)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, testError)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Put("/v1/tags/:tagId", UpdateTagHandler(tt.tagRepoMock(mc)))

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
