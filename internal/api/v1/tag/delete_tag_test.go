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

	"github.com/dmitriygnatenko/homethings-v1/internal/api/v1/tag/mocks"
	"github.com/dmitriygnatenko/homethings-v1/internal/dto"
	"github.com/dmitriygnatenko/homethings-v1/internal/helpers/test"
)

func TestDeleteTagHandler(t *testing.T) {
	t.Parallel()

	type req struct {
		method      string
		route       string
		contentType string
	}

	var (
		tagID     = uint64(gofakeit.Number(1, 1000))
		testError = gofakeit.Error()

		txMockFunc = func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		}

		correctReq = req{
			method:      fiber.MethodDelete,
			route:       "/v1/tags/" + strconv.FormatUint(tagID, 10),
			contentType: fiber.MIMEApplicationJSON,
		}
	)

	tests := []struct {
		name             string
		req              req
		resCode          int
		resBody          interface{}
		tmMock           func(mc *minimock.Controller) TransactionManager
		tagRepoMock      func(mc *minimock.Controller) TagRepository
		thingTagRepoMock func(mc *minimock.Controller) ThingTagRepository
	}{
		{
			name: "negative case - bad request",
			req: req{
				method:      fiber.MethodDelete,
				route:       "/v1/tags/" + gofakeit.Word(),
				contentType: fiber.MIMEApplicationJSON,
			},
			resCode: fiber.StatusBadRequest,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				return mocks.NewTagRepositoryMock(mc)
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - tag is not exist",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, sql.ErrNoRows)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - tag repository error (get)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				return mocks.NewTransactionManagerMock(mc)
			},
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				return mocks.NewThingTagRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - thing tag repository error (delete)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteByTagIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(testError)

				return mock
			},
		},
		{
			name:    "negative case - tag repository error (delete)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(testError)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteByTagIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil)

				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: dto.EmptyResponse{},
			tmMock: func(mc *minimock.Controller) TransactionManager {
				mock := mocks.NewTransactionManagerMock(mc)
				mock.ReadCommittedMock.Set(txMockFunc)
				return mock
			},
			tagRepoMock: func(mc *minimock.Controller) TagRepository {
				mock := mocks.NewTagRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil)

				return mock
			},
			thingTagRepoMock: func(mc *minimock.Controller) ThingTagRepository {
				mock := mocks.NewThingTagRepositoryMock(mc)

				mock.DeleteByTagIDMock.Inspect(func(ctx context.Context, id uint64) {
					assert.Equal(mc, tagID, id)
				}).Return(nil)

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			fiberApp := fiber.New()
			fiberApp.Delete("/v1/tags/:tagId", DeleteTagHandler(
				tt.tmMock(mc),
				tt.tagRepoMock(mc),
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
