// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"os"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	kt "github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

const OTEL_INSTRUMENTATION_KRATOS_EXPERIMENTAL_SPAN_ATTRIBUTES = "OTEL_INSTRUMENTATION_KRATOS_EXPERIMENTAL_SPAN_ATTRIBUTES"

var kratosEnabler = instrumenter.NewDefaultInstrumentEnabler()

var kratosInternalInstrument = BuildKratosInternalInstrumenter()

func kratosNewGRPCServiceOnEnter(call api.CallContext, opts ...grpc.ServerOption) {
	if os.Getenv(OTEL_INSTRUMENTATION_KRATOS_EXPERIMENTAL_SPAN_ATTRIBUTES) != "true" {
		return
	}
	opts = append(opts, AddGRPCMiddleware(ServerTracingMiddleWare()))
	call.SetParam(0, opts)
}

func AddHTTPMiddleware(m middleware.Middleware) http.ServerOption {
	return func(o *http.Server) {
		o.Use("*", m)
	}
}

func AddGRPCMiddleware(m middleware.Middleware) grpc.ServerOption {
	return func(o *grpc.Server) {
		o.Use("*", m)
	}
}

func ServerTracingMiddleWare() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				serviceName, serviceId, serviceVersion := "", "", ""
				serviceEndpoint := make([]string, 0, 0)
				serviceMeta := make(map[string]string)
				app, hasApp := kt.FromContext(ctx)
				if hasApp {
					serviceName, serviceId, serviceVersion, serviceEndpoint = app.Name(), app.ID(), app.Version(), app.Endpoint()
					serviceMeta = app.Metadata()
				}
				var (
					request kratosRequest
					sCtx    context.Context
				)
				request = kratosRequest{
					serviceId:       serviceId,
					serviceName:     serviceName,
					serviceVersion:  serviceVersion,
					serviceEndpoint: serviceEndpoint,
					serviceMeta:     serviceMeta,
				}
				switch tr.Kind() {
				case transport.KindGRPC:
					request.protocolType = "grpc"
					sCtx = kratosInternalInstrument.Start(ctx, request)
				case transport.KindHTTP:
					request.protocolType = "http"
					sCtx = kratosInternalInstrument.Start(ctx, request)
				}
				defer func() {
					if err != nil {
						kratosInternalInstrument.End(sCtx, request, nil, err)
					} else {
						kratosInternalInstrument.End(sCtx, request, nil, err)
					}
				}()

			}
			return handler(ctx, req)
		}
	}
}
