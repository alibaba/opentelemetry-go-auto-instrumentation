// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
)

var hertzServerEnabler = instrumenter.NewDefaultInstrumentEnabler()

var hertzInstrumenter = BuildHertzServerInstrumenter()

type hertzOpentelemetryTracer struct{}

func (m *hertzOpentelemetryTracer) Start(ctx context.Context, c *app.RequestContext) context.Context {
	return ctx
}

func (m *hertzOpentelemetryTracer) Finish(ctx context.Context, c *app.RequestContext) {
	if c.GetTraceInfo().Stats().GetEvent(stats.HTTPStart) != nil && c.GetTraceInfo().Stats().GetEvent(stats.HTTPFinish) != nil {
		start := c.GetTraceInfo().Stats().GetEvent(stats.HTTPStart)
		end := c.GetTraceInfo().Stats().GetEvent(stats.HTTPFinish)
		if ctx == nil {
			ctx = context.Background()
		}
		s := start.Time()
		e := end.Time()
		req, resp := &c.Request, &c.Response
		hertzInstrumenter.StartAndEnd(ctx, req, resp, c.GetTraceInfo().Stats().Error(), s, e)
	}
}

func beforeHertzServerBuild(call api.CallContext, opts ...config.Option) {
	if !hertzServerEnabler.Enable() {
		return
	}
	opts = append(opts, server.WithTracer(&hertzOpentelemetryTracer{}))
	call.SetParam(0, opts)
}
