package middlewares

import (
	"net/http"
	"story-book/package/services/jwtservice"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(jwtService jwtservice.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrMissingAuthorizationHeader)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, ErrInvalidAuthorizationHeader)
			}

			claims, err := jwtService.ParseJWT(parts[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			c.Set("id", claims["sub"])
			c.Set("role", claims["role"])

			return next(c)
		}
	}
}
