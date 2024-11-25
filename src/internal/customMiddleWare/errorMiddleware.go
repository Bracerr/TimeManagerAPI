package customMiddleWare

import (
	"TimeManagerAuth/src/pkg/customErrors"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		if err != nil {
			var customErr *customErrors.AppError
			if errors.As(err, &customErr) {
				return c.JSON(customErr.Code, customErr)
			}
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": err.Error()})
		}
		return nil
	}
}
