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
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

func KratosWithMiddlewareOnEnter(call transhttp.CallContext, m ...middleware.Middleware) {
	m = append(m, ClientTracingMiddleWare())
	call.SetParam(0, m)
}

// If user code called NewHttpClient() and middleware option is missing in arg opts, we can inject our tracing middleware by the following logic.
// Otherwise, it will be overwritten by the latter user-given middleware option in the opts.
// The only way user can construct middleware option is through function WithMiddleware(), so we inject additional logics in that.
func KratosNewHTTPClientOnEnter(call transhttp.CallContext, ctx context.Context, opts ...transhttp.ClientOption) {
	nopts := []transhttp.ClientOption{
		transhttp.WithMiddleware(ClientTracingMiddleWare()),
	}
	nopts = append(nopts, opts...)
	call.SetParam(1, nopts)
}

func ClientTracingMiddleWare() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				var (
					sCtx    context.Context
					request kratosRequest
				)
				switch tr.Kind() {
				case transport.KindGRPC:
					remote, _ := otelParseTarget(tr.Endpoint())
					request = kratosRequest{
						method:        tr.Operation(),
						componentName: "kratos-grpc-client",
						header:        tr.RequestHeader(),
						addr:          remote,
					}
					sCtx = kratosClientInstrument.Start(ctx, request)
				case transport.KindHTTP:
					if ht, ok := tr.(transhttp.Transporter); ok {
						remote := ht.Request().Host
						fmt.Println("%%%%%%%%%%%%%")
						fmt.Println(ht.Request().Method)
						fmt.Println("%%%%%%%%%%%%%")
						request = kratosRequest{
							method:        tr.Operation(),
							componentName: "kratos-http-client",
							header:        tr.RequestHeader(),
							addr:          remote,
							httpMethod:    ht.Request().Method,
						}
						sCtx = kratosClientInstrument.Start(ctx, request)
					}

				}
				defer func() {
					if err != nil {
						var code int
						if e := errors.FromError(err); e != nil {
							code = int(e.Code)
						}
						kratosClientInstrument.End(sCtx, request, kratosResponse{
							statusCode: code,
						}, err)
					} else {
						kratosClientInstrument.End(sCtx, request, kratosResponse{
							statusCode: 200,
						}, err)
					}
				}()
			}
			return handler(ctx, req)
		}
	}
}
