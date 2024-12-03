// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		verifier.NewSpanVerifier().
			HasStringAttribute("kratos.protocol.type", "grpc").
			HasStringAttribute("kratos.service.name", "opentelemetry-kratos-server").
			HasStringAttribute("kratos.service.id", "opentelemetry-id").
			HasStringAttribute("kratos.service.version", "v1").
			HasStringAttribute("kratos.service.meta.agent", "opentelemetry-go").
			HasItemInStringSliceAttribute("kratos.service.endpoint", 0, func(s string) (bool, string) {
				return strings.Contains(s, ":9000"), fmt.Sprintf("First endpoint should be xxx:9000, actual value: %v", s)
			}).
			HasItemInStringSliceAttribute("kratos.service.endpoint", 1, func(s string) (bool, string) {
				return strings.Contains(s, ":8000"), fmt.Sprintf("First endpoint should be xxx:8000, actual value: %v", s)
			}).
			Verify(stubs[0][2])
	}, 1)
}
