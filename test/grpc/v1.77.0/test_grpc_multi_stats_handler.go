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
	"log"
	"net"
	"time"

	"github.com/alibaba/loongsuite-go/test/verifier"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/stats"
)

// extraStatsHandler mimics otelgrpc: it starts a span in TagRPC and ends
// trace.SpanFromContext(ctx) in HandleRPC. Registering it next to the agent's
// handler means ctx carries this span, not the agent's, during HandleRPC.
type extraStatsHandler struct{}

func (h *extraStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	ctx, _ = otel.Tracer("extra-stats-handler").Start(ctx, "extra-stats-handler")
	return ctx
}

func (h *extraStatsHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	if _, ok := rs.(*stats.End); ok {
		trace.SpanFromContext(ctx).End()
	}
}

func (h *extraStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

func (h *extraStatsHandler) HandleConn(ctx context.Context, info stats.ConnStats) {}

func setupGRPCWithExtraStatsHandler() {
	lis, err := net.Listen("tcp", "0.0.0.0:9003")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StatsHandler(&extraStatsHandler{}))
	RegisterHelloGrpcServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	go setupGRPCWithExtraStatsHandler()
	time.Sleep(3 * time.Second)
	conn, err := grpc.NewClient("localhost:9003",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(&extraStatsHandler{}))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := NewHelloGrpcClient(conn)
	if _, err := c.Hello(context.Background(), &Req{}); err != nil {
		panic(err)
	}
	// the agent's own client and server spans must still be ended and exported
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		var client, server *tracetest.SpanStub
		for i := range stubs {
			for j := range stubs[i] {
				switch stubs[i][j].SpanKind {
				case trace.SpanKindClient:
					client = &stubs[i][j]
				case trace.SpanKindServer:
					server = &stubs[i][j]
				}
			}
		}
		verifier.Assert(client != nil, "Expect the agent's client span to be ended")
		verifier.Assert(server != nil, "Expect the agent's server span to be ended")
		verifier.VerifyRpcClientAttributes(*client, "/HelloGrpc/Hello", "grpc", "/HelloGrpc", "Hello")
		verifier.VerifyRpcServerAttributes(*server, "/HelloGrpc/Hello", "grpc", "/HelloGrpc", "Hello")
	}, 1)
	// leaked client spans pile up in this goroutine's TraceContext; beyond its
	// limit every new span prints "Failed to add span to TraceContext"
	for i := 0; i < 400; i++ {
		if _, err := c.Hello(context.Background(), &Req{}); err != nil {
			panic(err)
		}
	}
}
