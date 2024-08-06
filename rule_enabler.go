package main

import (
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/http"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/mongo"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/otsdk"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/runtime"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/test"
)
