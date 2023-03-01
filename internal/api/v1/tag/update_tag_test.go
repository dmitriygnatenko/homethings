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

func Test_UpdateTagHandler(t *testing.T) {
	type tagRepoMockFunc func(mc *minimock.Controller) interfaces.TagRepository

	type req struct {
		method      string
		route       string
		body        *dto.UpdateTagRequest
		contentType string
	}

	var (
		mc        = minimock.NewController(t)
		tagID     = gofakeit.Number(1, 1000)
		title     = gofakeit.Phrase()
		style     = gofakeit.Phrase()
		testError = errors.New(gofakeit.Phrase())
		layout    = "2006-01-02 15:04:05"

		correctReq = req{
			method: fiber.MethodPut,
			route:  "/v1/tags/" + strconv.Itoa(tagID),
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
		tagRepoMock tagRepoMockFunc
	}{
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: expectedRes,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)

				mock.UpdateMock.Inspect(func(ctx context.Context, req models.UpdateTagRequest, tx *sql.Tx) {
					assert.Equal(mc, tagID, req.ID)
					assert.Equal(mc, title, req.Title)
					assert.Equal(mc, style, req.Style)
				}).Return(nil)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
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
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				return repoMocks.NewTagRepositoryMock(mc)
			},
		},
		{
			name: "negative case - body parse error",
			req: req{
				method: fiber.MethodPut,
				route:  "/v1/tags/" + strconv.Itoa(tagID),
			},
			resCode: fiber.StatusBadRequest,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				return repoMocks.NewTagRepositoryMock(mc)
			},
		},
		{
			name: "negative case - request without title",
			req: req{
				method:      fiber.MethodPut,
				route:       "/v1/tags/" + strconv.Itoa(tagID),
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
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				return repoMocks.NewTagRepositoryMock(mc)
			},
		},
		{
			name: "negative case - request without style",
			req: req{
				method:      fiber.MethodPut,
				route:       "/v1/tags/" + strconv.Itoa(tagID),
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
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				return repoMocks.NewTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (update)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)
				mock.UpdateMock.Return(testError)
				return mock
			},
		},
		{
			name:    "negative case - repository error (get)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tagRepoMock: func(mc *minimock.Controller) interfaces.TagRepository {
				mock := repoMocks.NewTagRepositoryMock(mc)
				mock.UpdateMock.Return(nil)
				mock.GetMock.Return(nil, testError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(tt.tagRepoMock(mc))

			fiberApp.Put("/v1/tags/:tagId", UpdateTagHandler(serviceProvider))

			fiberReq := httptest.NewRequest(tt.req.method, tt.req.route, helpers.ConvertDataToIOReader(tt.req.body))
			fiberReq.Header.Add(fiber.HeaderContentType, tt.req.contentType)
			fiberRes, _ := fiberApp.Test(fiberReq, API.DefaultTestTimeOut)

			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
