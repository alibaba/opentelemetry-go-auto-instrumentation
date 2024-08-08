package main

import (
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/databasesql"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/goredis"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/mongo"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/test"
)
