package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const (
	ClaimsKeyName = "name"
	claimsKeyExp  = "exp"
	claimsKeyUser = "user"
	defaultCost   = bcrypt.DefaultCost
)

type Env interface {
	JWTSecretKey() string
	JWTLifetime() int
}

type Service struct {
	env Env
}

func Init(env Env) (*Service, error) {
	return &Service{env: env}, nil
}

func (a Service) GeneratePasswordHash(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (a Service) IsCorrectPassword(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (a Service) GetClaims(fctx *fiber.Ctx) jwt.MapClaims {
	jwtUser := fctx.Locals(claimsKeyUser).(*jwt.Token)
	return jwtUser.Claims.(jwt.MapClaims)
}

func (a Service) GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		ClaimsKeyName: user.Username,
		claimsKeyExp:  time.Now().Add(time.Duration(a.env.JWTLifetime()) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(a.env.JWTSecretKey()))
}
