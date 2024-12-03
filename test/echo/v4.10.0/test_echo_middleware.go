// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func loggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			log.Printf("method: %s, path: %s, proto: %s, "+
				"remoteaddr: %s, response_status: %d, cost: %s",
				c.Request().Method, c.Path(), c.Request().Proto, c.Request().RemoteAddr,
				c.Response().Status, time.Since(start))
			return err
		}
	}
}

func setupMiddleware() {
	engine := echo.New()
	engine.Use(middleware.Logger())
	engine.Use(loggingMiddleware())
	engine.GET("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 1,
			"msg":  c.Path(),
		})
	})

	// Start server
	engine.Logger.Fatal(engine.Start(":8080"))
}

func main() {
	go setupMiddleware()
	time.Sleep(5 * time.Second)
	client := &http.Client{}
	resp, err := client.Get("http://127.0.0.1:8080/test")
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://127.0.0.1:8080/test", "http", "1.1", "tcp", "ipv4", "", "127.0.0.1:8080", 200, 0, 8080)
		verifier.VerifyHttpServerAttributes(stubs[0][1], "GET /test", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:8080", "Go-http-client/1.1", "http", "/test", "", "/test", 200)
	}, 1)
}
