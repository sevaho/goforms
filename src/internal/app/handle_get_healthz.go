package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func handleGetHealthz() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	}
}
