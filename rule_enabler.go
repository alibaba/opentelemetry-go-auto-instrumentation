package main

import (
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/baggage"
	_ "go.opentelemetry.io/otel/sdk/trace"

	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/databasesql"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/goredis"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/http"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/mongo"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/otsdk"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/runtime"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/test"
)
