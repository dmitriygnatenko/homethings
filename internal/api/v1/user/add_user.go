package user

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i AuthService,UserRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"strings"

	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/dmitriygnatenko/homethings/internal/dto"
	"github.com/dmitriygnatenko/homethings/internal/factory"
	"github.com/dmitriygnatenko/homethings/internal/models"
)

type (
	AuthService interface {
		GeneratePasswordHash(password string) (string, error)
		GetClaims(fctx *fiber.Ctx) jwt.MapClaims
	}

	UserRepository interface {
		Get(ctx context.Context, username string) (*models.User, error)
		Add(ctx context.Context, username string, password string) (uint64, error)
		Update(ctx context.Context, req models.UpdateUserRequest) error
	}
)

// @Router 		/api/v1/users [post]
// @Param       data body dto.AddUserRequest true "Request body"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Add user
// @Tags  		Users
// @security 	APIKey
// @Accept      json
// @Produce     json
func AddUserHandler(
	authService AuthService,
	userRepository UserRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.AddUserRequest{}

		if err := fctx.BodyParser(&req); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		var validate = validator.New()
		if err := validate.Struct(req); err != nil {
			logger.Info(ctx, err.Error())
			return fctx.Status(fiber.StatusBadRequest).JSON(factory.CreateValidateErrorResponse(err))
		}

		hash, err := authService.GeneratePasswordHash(strings.TrimSpace(req.Password))
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		_, err = userRepository.Add(ctx, strings.TrimSpace(req.Username), hash)
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
