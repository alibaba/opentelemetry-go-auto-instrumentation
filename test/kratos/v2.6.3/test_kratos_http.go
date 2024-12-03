// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"net/http"
	"strings"
	"time"
)

func main() {
	go func() {
		startup()
	}()
	time.Sleep(5 * time.Second)
	client := http.Client{}
	client.Get("http://localhost:8000/helloworld/kratos")
	fmt.Printf("Send http request to kratos")

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		for i, _ := range stubs[0] {
			span := stubs[0][i]
			println(span.Name)
			for _, attr := range span.Attributes {
				fmt.Printf("%v %v\n", attr.Key, attr.Value)
			}
			println()
		}
		protocolType := verifier.GetAttribute(stubs[0][2].Attributes, "kratos.protocol.type").AsString()
		if protocolType != "http" {
			panic("protocol type should be http, actually got " + protocolType)
		}
		serviceName := verifier.GetAttribute(stubs[0][2].Attributes, "kratos.service.name").AsString()
		if serviceName != "opentelemetry-kratos-server" {
			panic("service name should be opentelemetry-kratos-server, actually got " + serviceName)
		}
		serviceId := verifier.GetAttribute(stubs[0][2].Attributes, "kratos.service.id").AsString()
		if serviceId != "opentelemetry-id" {
			panic("service id should be opentelemetry-id, actually got " + serviceId)
		}
		serviceVersion := verifier.GetAttribute(stubs[0][2].Attributes, "kratos.service.version").AsString()
		if serviceVersion != "v1" {
			panic("service version should be v1, actually got " + serviceVersion)
		}
		serviceMetaAgent := verifier.GetAttribute(stubs[0][2].Attributes, "kratos.service.meta.agent").AsString()
		if serviceMetaAgent != "opentelemetry-go" {
			panic("service meta agent should be opentelemetry-go, actually got " + serviceMetaAgent)
		}
		serviceEndpoint := verifier.GetAttribute(stubs[0][2].Attributes, "kratos.service.endpoint").AsStringSlice()
		if !strings.Contains(serviceEndpoint[0], ":9000") || !strings.Contains(serviceEndpoint[1], ":8000") {
			panic("service endpoint should be grpc://30.221.144.142:9000 http://30.221.144.142:8000, actually got " + fmt.Sprintf("%v", serviceEndpoint))
		}
	}, 1)
}
