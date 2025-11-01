package app

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sevaho/goforms/src/config"
)

func QueryKeyAuth(fn middleware.KeyAuthValidator) echo.MiddlewareFunc {
	// Check for apiKey in query params to verify
	c := middleware.KeyAuthConfig{
		Skipper:   middleware.DefaultSkipper,
		KeyLookup: "query:apiKey",
	}
	c.Validator = fn
	return middleware.KeyAuthWithConfig(c)
}

func CheckAuthorizationBearerTokenMiddleware(config *config.Config) echo.MiddlewareFunc {
	return middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		fmt.Println(c.QueryParams())
		return config.VerifyApiKey(key), nil
	})
}

func CheckApiTokenQueryParamsMiddleware(config *config.Config) echo.MiddlewareFunc {
	return QueryKeyAuth(func(key string, c echo.Context) (bool, error) {
		fmt.Println(c.QueryParams())
		return config.VerifyApiKey(key), nil
	})
}
