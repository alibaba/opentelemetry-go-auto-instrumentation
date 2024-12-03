// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func main() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		logger := zap.NewExample()
		logger.Debug("this is debug message")
		logger.Info("this is info message")
		logger.Info("this is info message with fileds",
			zap.Int("age", 37),
			zap.String("agender", "man"),
		)
		logger.Warn("this is warn message")
		logger.Error("this is error message")
	})

	http.HandleFunc("/logwithtrace", func(w http.ResponseWriter, r *http.Request) {
		logger := zap.NewExample()
		// GetTraceAndSpanId will be added while using otel, users must use otel to build the module
		traceId, spanId := trace.GetTraceAndSpanId()
		logger.Info("this is info message with fileds",
			zap.String("traceId", traceId),
			zap.String("spanId", spanId),
		)
	})
	http.ListenAndServe(":9999", nil)
}
