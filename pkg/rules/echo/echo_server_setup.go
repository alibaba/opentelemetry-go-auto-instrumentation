//go:build ignore

package rule

import (
	"strconv"

	echo "github.com/labstack/echo/v4"
)

func otelTraceMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			request := echoRequest{
				method:  c.Request().Method,
				path:    c.Path(),
				url:     *c.Request().URL,
				header:  c.Request().Header,
				version: strconv.Itoa(c.Request().ProtoMajor) + "." + strconv.Itoa(c.Request().ProtoMinor),
				host:    c.Request().Host,
				isTls:   c.Request().TLS != nil,
			}
			ctx := netEchoServerInstrument.Start(c.Request().Context(), request)
			if err = next(c); err != nil {
				c.Error(err)
			}

			netEchoServerInstrument.End(ctx, request, echoResponse{
				statusCode: c.Response().Status,
				header:     c.Response().Header(),
			}, err)
			return
		}
	}
}

func afterNewEcho(call echo.CallContext, e *echo.Echo) {
	if e == nil {
		return
	}

	e.Use(otelTraceMiddleware())
}
