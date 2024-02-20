package pkgs

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func SetupEcho() {
	r := echo.New()

	r.GET("/echo-service1", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Echo")
	})
	_ = r.Start(":9005")
}
