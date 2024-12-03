// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
)

var hertzClientEnabler = instrumenter.NewDefaultInstrumentEnabler()

var hertzClientInstrumenter = BuildHertzClientInstrumenter()

func otelClientMiddleware(next client.Endpoint) client.Endpoint {
	return func(ctx context.Context, req *protocol.Request, resp *protocol.Response) (err error) {
		ctx = hertzClientInstrumenter.Start(ctx, req)
		err = next(ctx, req, resp)
		if err != nil {
			hertzClientInstrumenter.End(ctx, req, resp, err)
			return err
		}
		hertzClientInstrumenter.End(ctx, req, resp, nil)
		return nil
	}
}

func afterHertzClientBuild(call api.CallContext, c *client.Client, err error) {
	if !hertzClientEnabler.Enable() {
		return
	}
	if err != nil {
		return
	}
	c.Use(otelClientMiddleware)
}
