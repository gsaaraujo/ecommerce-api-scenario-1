package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func NewEchoJWTMiddleware(accessTokenSigningKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(accessTokenSigningKey),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(usecases.JwtAccessTokenClaims)
		},
	})
}
