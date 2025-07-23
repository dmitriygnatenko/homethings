package user

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/dmitriygnatenko/go-common/logger"
	"github.com/gofiber/fiber/v2"

	"github.com/dmitriygnatenko/homethings/internal/dto"
	"github.com/dmitriygnatenko/homethings/internal/factory"
	"github.com/dmitriygnatenko/homethings/internal/mappers"
	"github.com/dmitriygnatenko/homethings/internal/services/auth"
)

// @Router 		/api/v1/users [put]
// @Param       data body dto.UpdateUserRequest true "Request body"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update user
// @Tags  		Users
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdateUserHandler(
	authService AuthService,
	userRepository UserRepository,
) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		var err error

		var username, password string

		ctx := fctx.Context()
		req := dto.UpdateUserRequest{}

		if err = fctx.BodyParser(&req); err != nil {
			logger.Info(ctx, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if req.Username != nil {
			username = strings.TrimSpace(*req.Username)
		}

		if req.Password != nil {
			password, err = authService.GeneratePasswordHash(strings.TrimSpace(*req.Password))
			if err != nil {
				logger.Error(ctx, err.Error())
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		if username == "" && password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "")
		}

		claims := authService.GetClaims(fctx)

		user, err := userRepository.Get(ctx, claims[auth.ClaimsKeyName].(string))
		if err != nil {
			logger.Error(ctx, err.Error())

			if errors.Is(err, sql.ErrNoRows) {
				return fiber.NewError(fiber.StatusBadRequest, "")
			}

			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		req.Password = &password
		req.Username = &username

		err = userRepository.Update(ctx, mappers.ToUpdateUserRequest(user.ID, req))
		if err != nil {
			logger.Error(ctx, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
