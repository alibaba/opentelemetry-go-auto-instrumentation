// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

var grpcServerInstrument = BuildGrpcServerInstrumenter()

func grpcServerOnEnter(call api.CallContext, opts ...grpc.ServerOption) {
	if !grpcEnabler.Enable() {
		return
	}
	h := grpc.StatsHandler(NewServerHandler())
	var opt []grpc.ServerOption
	opt = append(opt, h)
	opt = append(opt, opts...)
	call.SetParam(0, opt)
}

func grpcServerOnExit(call api.CallContext, s *grpc.Server) {
	if !grpcEnabler.Enable() {
		return
	}
	return
}

func NewServerHandler(opts ...Option) stats.Handler {
	h := &serverHandler{
		grpcOtelConfig: newConfig(opts, "server"),
	}

	return h
}

type serverHandler struct {
	*grpcOtelConfig
}

// TagConn can attach some information to the given context.
func (h *serverHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn processes the Conn stats.
func (h *serverHandler) HandleConn(ctx context.Context, info stats.ConnStats) {
}

// TagRPC can attach some information to the given context.
func (h *serverHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	var md metadata.MD
	ctx, md = extract(ctx, h.grpcOtelConfig.Propagators)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	nCtx := grpcServerInstrument.Start(ctx, grpcRequest{
		methodName: info.FullMethodName,
		propagators: &grpcMetadataSupplier{
			metadata: &md,
		},
	})

	gctx := gRPCContext{
		methodName: info.FullMethodName,
	}

	return context.WithValue(nCtx, gRPCContextKey{}, &gctx)
}

// HandleRPC processes the RPC stats.
func (h *serverHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	isServer := true
	h.handleRPC(ctx, rs, isServer)
}
