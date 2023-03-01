package place

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

func Test_DeletePlaceHandler(t *testing.T) {
	type placeRepoMockFunc func(mc *minimock.Controller) interfaces.PlaceRepository
	type thingRepoMockFunc func(mc *minimock.Controller) interfaces.ThingRepository
	type placeThingRepoMockFunc func(mc *minimock.Controller) interfaces.PlaceThingRepository
	type placeImageRepoMockFunc func(mc *minimock.Controller) interfaces.PlaceImageRepository
	type thingImageRepoMockFunc func(mc *minimock.Controller) interfaces.ThingImageRepository
	type fileRepoMockFunc func(mc *minimock.Controller) interfaces.FileRepository

	type req struct {
		method string
		route  string
	}

	var (
		mc      = minimock.NewController(t)
		placeID = gofakeit.Number(1, 1000)
		thingID = gofakeit.Number(1, 1000)

		placeImageID  = gofakeit.Number(1, 1000)
		placeImageURL = gofakeit.URL()
		thingImageID  = gofakeit.Number(1, 1000)
		thingImageURL = gofakeit.URL()
		testError     = errors.New(gofakeit.Phrase())

		correctReq = req{
			method: fiber.MethodDelete,
			route:  "/v1/places/" + strconv.Itoa(placeID),
		}

		thingRepoRes = []models.Thing{
			{
				ID:      thingID,
				PlaceID: placeID,
			},
		}

		placeImageRepoRes = []models.Image{
			{
				ID:    placeImageID,
				Image: placeImageURL,
			},
		}

		thingImageRepoRes = []models.Image{
			{
				ID:    thingImageID,
				Image: thingImageURL,
			},
		}
	)

	tests := []struct {
		name               string
		req                req
		resCode            int
		resBody            interface{}
		placeRepoMock      placeRepoMockFunc
		thingRepoMock      thingRepoMockFunc
		placeThingRepoMock placeThingRepoMockFunc
		thingImageRepoMock thingImageRepoMockFunc
		placeImageRepoMock placeImageRepoMockFunc
		fileRepoMock       fileRepoMockFunc
	}{
		{
			name: "negative case - bad request",
			req: req{
				method: fiber.MethodDelete,
				route:  "/v1/places/" + gofakeit.Word(),
			},
			resCode: fiber.StatusBadRequest,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				return repoMocks.NewPlaceRepositoryMock(mc)
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				return repoMocks.NewPlaceImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - bad request (place not found)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, sql.ErrNoRows)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				return repoMocks.NewPlaceImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get place)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, testError)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				return repoMocks.NewPlaceImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get nested places)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, testError)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				return repoMocks.NewPlaceImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - bad request (nested places exists)",
			req:     correctReq,
			resCode: fiber.StatusBadRequest,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return([]models.Place{{}}, nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				return repoMocks.NewPlaceImageRepositoryMock(mc)
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get place images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				return repoMocks.NewThingRepositoryMock(mc)
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(nil, testError)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get things)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(nil, testError)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				return repoMocks.NewThingImageRepositoryMock(mc)
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(nil, nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (get things images)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(nil, testError)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(nil, nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (begin tx)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, testError)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(thingImageRepoRes, nil)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(nil, nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete place image)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(thingImageRepoRes, nil)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(placeImageRepoRes, nil)
				mock.DeleteMock.Return(testError)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete thing image)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				return repoMocks.NewPlaceThingRepositoryMock(mc)
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(thingImageRepoRes, nil)
				mock.DeleteMock.Return(testError)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(placeImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete place thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(testError)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(thingImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(placeImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete thing)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
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
				mock.GetByThingIDMock.Return(thingImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(placeImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - repository error (delete place)",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(testError)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(thingImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(placeImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
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
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				mock.CommitTxMock.Return(testError)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(thingImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(placeImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				return repoMocks.NewFileRepositoryMock(mc)
			},
		},
		{
			name:    "negative case - delete place image error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				mock.CommitTxMock.Return(nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(thingImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(placeImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				mock := repoMocks.NewFileRepositoryMock(mc)
				mock.DeleteMock.Expect(placeImageURL).Return(testError)
				return mock
			},
		},
		{
			name:    "negative case - delete thing image error",
			req:     correctReq,
			resCode: fiber.StatusInternalServerError,
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)
				mock.GetMock.Return(nil, nil)
				mock.GetNestedPlacesMock.Return(nil, nil)
				mock.BeginTxMock.Return(nil, nil)
				mock.DeleteMock.Return(nil)
				mock.CommitTxMock.Return(nil)
				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(thingRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)
				mock.DeleteThingMock.Return(nil)
				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)
				mock.GetByThingIDMock.Return(thingImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)
				mock.GetByPlaceIDMock.Return(placeImageRepoRes, nil)
				mock.DeleteMock.Return(nil)
				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				mock := repoMocks.NewFileRepositoryMock(mc)
				mock.DeleteMock.When(placeImageURL).Then(nil)
				mock.DeleteMock.When(thingImageURL).Then(testError)
				return mock
			},
		},
		{
			name:    "positive case",
			req:     correctReq,
			resCode: fiber.StatusOK,
			resBody: &dto.EmptyResponse{},
			placeRepoMock: func(mc *minimock.Controller) interfaces.PlaceRepository {
				mock := repoMocks.NewPlaceRepositoryMock(mc)

				mock.GetMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.GetNestedPlacesMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(nil, nil)

				mock.BeginTxMock.Return(nil, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, placeID, id)
				}).Return(nil)

				mock.CommitTxMock.Return(nil)

				return mock
			},
			thingRepoMock: func(mc *minimock.Controller) interfaces.ThingRepository {
				mock := repoMocks.NewThingRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(thingRepoRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			placeThingRepoMock: func(mc *minimock.Controller) interfaces.PlaceThingRepository {
				mock := repoMocks.NewPlaceThingRepositoryMock(mc)

				mock.DeleteThingMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, thingID, id)
				}).Return(nil)

				return mock
			},
			thingImageRepoMock: func(mc *minimock.Controller) interfaces.ThingImageRepository {
				mock := repoMocks.NewThingImageRepositoryMock(mc)

				mock.GetByThingIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, thingID, id)
				}).Return(thingImageRepoRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, thingImageID, id)
				}).Return(nil)

				return mock
			},
			placeImageRepoMock: func(mc *minimock.Controller) interfaces.PlaceImageRepository {
				mock := repoMocks.NewPlaceImageRepositoryMock(mc)

				mock.GetByPlaceIDMock.Inspect(func(ctx context.Context, id int) {
					assert.Equal(mc, placeID, id)
				}).Return(placeImageRepoRes, nil)

				mock.DeleteMock.Inspect(func(ctx context.Context, id int, tx *sql.Tx) {
					assert.Equal(mc, placeImageID, id)
				}).Return(nil)

				return mock
			},
			fileRepoMock: func(mc *minimock.Controller) interfaces.FileRepository {
				mock := repoMocks.NewFileRepositoryMock(mc)
				mock.DeleteMock.When(placeImageURL).Then(nil)
				mock.DeleteMock.When(thingImageURL).Then(nil)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fiberApp := fiber.New()
			serviceProvider := sp.InitMock(
				tt.placeRepoMock(mc),
				tt.thingRepoMock(mc),
				tt.placeThingRepoMock(mc),
				tt.placeImageRepoMock(mc),
				tt.thingImageRepoMock(mc),
				tt.fileRepoMock(mc),
			)

			fiberApp.Delete("/v1/places/:placeId", DeletePlaceHandler(serviceProvider))

			fiberRes, _ := fiberApp.Test(httptest.NewRequest(tt.req.method, tt.req.route, nil), API.DefaultTestTimeOut)
			assert.Equal(t, tt.resCode, fiberRes.StatusCode)
			if tt.resBody != nil {
				assert.Equal(t, helpers.MarshalResponse(tt.resBody), helpers.ConvertBodyToString(fiberRes.Body))
			}
		})
	}
}
