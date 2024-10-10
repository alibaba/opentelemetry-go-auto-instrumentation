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
//go:build ignore

package rule

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

func KratosNewHTTPServiceOnEnter(call transhttp.CallContext, opts ...transhttp.ServerOption) {
	opts = append(opts, AddHTTPMiddleware(ServerTracingMiddleWare()))
	call.SetParam(0, opts)
}

func KratosNewGRPCServiceOnEnter(call grpc.CallContext, opts ...grpc.ServerOption) {
	opts = append(opts, AddGRPCMiddleware(ServerTracingMiddleWare()))
	call.SetParam(0, opts)
}

// AddMiddleware adds service middleware option.
func AddHTTPMiddleware(m middleware.Middleware) transhttp.ServerOption {
	return func(o *transhttp.Server) {
		o.Use("*", m)
	}
}

func AddGRPCMiddleware(m middleware.Middleware) grpc.ServerOption {
	return func(o *grpc.Server) {
		o.Use("*", m)
	}
}

func KratosGRPCWithMiddlewareOnEnter(call grpc.CallContext, m ...middleware.Middleware) {
	m = append(m, ClientTracingMiddleWare())
	call.SetParam(0, m)
}

func ServerTracingMiddleWare() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				var (
					request kratosRequest
					sCtx    context.Context
				)
				switch tr.Kind() {
				case transport.KindGRPC:
					remote, _ := otelParseTarget(tr.Endpoint())
					request = kratosRequest{
						method:        tr.Operation(),
						componentName: "kratos-grpc-server",
						header:        tr.RequestHeader(),
						addr:          remote,
					}
					sCtx = kratosServerInstrument.Start(ctx, request)
				case transport.KindHTTP:
					if ht, ok := tr.(transhttp.Transporter); ok {
						remote := ht.Request().Host
						request = kratosRequest{
							method:        ht.Request().URL.Path,
							componentName: "kratos-http-server",
							header:        tr.RequestHeader(),
							addr:          remote,
							httpMethod:    ht.Request().Method,
						}
						sCtx = kratosServerInstrument.Start(ctx, request)
					}
				}
				defer func() {
					if err != nil {
						var code int
						if e := errors.FromError(err); e != nil {
							code = int(e.Code)
						}
						kratosServerInstrument.End(sCtx, request, kratosResponse{
							statusCode: code,
						}, err)
					} else {
						kratosServerInstrument.End(sCtx, request, kratosResponse{
							statusCode: 200,
						}, err)
					}

				}()

			}
			return handler(ctx, req)
		}
	}
}
