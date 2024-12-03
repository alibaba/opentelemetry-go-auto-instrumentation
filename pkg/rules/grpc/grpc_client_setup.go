// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

func grpcClientOnEnter(call api.CallContext, ctx context.Context, target string, opts ...grpc.DialOption) {
	if !grpcEnabler.Enable() {
		return
	}
	h := grpc.WithStatsHandler(NewClientHandler())
	var opt []grpc.DialOption
	opt = append(opt, h)
	opt = append(opt, opts...)
	call.SetParam(2, opt)
}

func grpcClientOnExit(call api.CallContext, cc *grpc.ClientConn, err error) {
	if !grpcEnabler.Enable() {
		return
	}
	return
}

type clientHandler struct {
	*grpcOtelConfig
}

func NewClientHandler(opts ...Option) stats.Handler {
	h := &clientHandler{
		grpcOtelConfig: newConfig(opts, "client"),
	}

	return h
}

// TagRPC can attach some information to the given context.
func (h *clientHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	nCtx := grpcClientInstrument.Start(ctx, grpcRequest{
		methodName: info.FullMethodName,
	})
	gctx := gRPCContext{
		methodName: info.FullMethodName,
	}

	return inject(context.WithValue(nCtx, gRPCContextKey{}, &gctx), h.grpcOtelConfig.Propagators, info.FullMethodName)
}

// HandleRPC processes the RPC stats.
func (h *clientHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	isServer := false
	h.handleRPC(ctx, rs, isServer)
}

// TagConn can attach some information to the given context.
func (h *clientHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn processes the Conn stats.
func (h *clientHandler) HandleConn(context.Context, stats.ConnStats) {
	// no-op
}
