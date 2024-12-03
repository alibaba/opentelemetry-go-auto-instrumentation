// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	pb "kratos/v2.5.2/pkg/api/helloworld/v1"
	"strings"
	"time"
)

func main() {
	go func() {
		startup()
	}()
	time.Sleep(5 * time.Second)
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("localhost:9000"),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "client"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[grpc] SayHello %+v\n", reply)

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		protocolType := verifier.GetAttribute(stubs[0][2].Attributes, "kratos.protocol.type").AsString()
		if protocolType != "grpc" {
			panic("protocol type should be grpc, actually got " + protocolType)
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
