package customMiddleWare

import (
	"TimeManagerAuth/src/internal/repository"
	"TimeManagerAuth/src/pkg/auth"
	"TimeManagerAuth/src/pkg/customErrors"
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

func JWTMiddleware(jwtManager *auth.Manager, userRepo *repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return customErrors.NewAppError(http.StatusUnauthorized, "missing token")
			}

			tokenString = strings.TrimPrefix(tokenString, "Bearer ")

			validateToken, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				return customErrors.NewAppError(http.StatusUnauthorized, "invalid token")
			}

			user, err := userRepo.FindUserByLogin(validateToken.Login)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					return customErrors.NewAppError(http.StatusUnauthorized, "user not found")
				}
				return err
			}

			c.Set("userID", user.ID.Hex())

			return next(c)
		}
	}
}
