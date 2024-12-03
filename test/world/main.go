// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

// This test matches rules as much as possible and check if compilation works
import (
	_ "database/sql"
	_ "log"
	_ "log/slog"
	_ "net/http"
	_ "runtime"

	_ "github.com/cloudwego/hertz/pkg/app/server"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-kratos/kratos/v2/transport/http"
	_ "github.com/go-redis/redis/v8"
	_ "github.com/gofiber/fiber/v2"
	_ "github.com/gorilla/mux"
	_ "github.com/labstack/echo/v4"
	_ "github.com/sirupsen/logrus"
	_ "github.com/valyala/fasthttp"
	_ "go.mongodb.org/mongo-driver/mongo"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/baggage"
	_ "go.opentelemetry.io/otel/trace"
	_ "go.uber.org/zap/zapcore"
	_ "google.golang.org/grpc"
	_ "gorm.io/driver/mysql"
	_ "gorm.io/gorm"
)

func main() {

}
