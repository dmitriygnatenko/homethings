package user

import (
	"database/sql"
	"strings"

	"git.dmitriygnatenko.ru/dima/homethings/internal/dto"
	"git.dmitriygnatenko.ru/dima/homethings/internal/factory"
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// @Router 		/api/v1/users [put]
// @Param       data body dto.UpdateUserRequest true "Request body"
// @Success     200 {object} dto.EmptyResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     500 {object} dto.ErrorResponse
// @Summary     Update user
// @Tags  		User
// @security 	APIKey
// @Accept      json
// @Produce     json
func UpdateUserHandler(sp interfaces.IServiceProvider) fiber.Handler {
	return func(fctx *fiber.Ctx) error {
		ctx := fctx.Context()
		req := dto.UpdateUserRequest{}
		if err := fctx.BodyParser(&req); err != nil {
			return factory.CreateBadRequestResponse(fctx, err)
		}

		var username, password string
		if req.Username != nil {
			username = strings.TrimSpace(*req.Username)
		}

		if req.Password != nil {
			hash, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(*req.Password)), bcrypt.DefaultCost)
			if err != nil {
				return factory.CreateInternalErrorResponse(fctx, err)
			}
			password = string(hash)
		}

		if username == "" && password == "" {
			return factory.CreateBadRequestResponse(fctx, nil)
		}

		jwtUser := fctx.Locals("user").(*jwt.Token)
		claims := jwtUser.Claims.(jwt.MapClaims)

		user, err := sp.GetUserRepository().Get(ctx, claims["name"].(string))
		if err != nil {
			if err == sql.ErrNoRows {
				return factory.CreateBadRequestResponse(fctx, nil)
			}

			return factory.CreateInternalErrorResponse(fctx, err)
		}

		req.Password = &password
		req.Username = &username

		err = sp.GetUserRepository().Update(ctx, mappers.ConvertToUpdateUserRequestModel(user.ID, req))
		if err != nil {
			return factory.CreateInternalErrorResponse(fctx, err)
		}

		return fctx.JSON(factory.CreateEmptyResponse())
	}
}
