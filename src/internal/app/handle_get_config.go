package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sevaho/goforms/src/internal/models"
)

func handleGetConfig(
	fc models.FormsConfig,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, fc, " ")
	}
}
