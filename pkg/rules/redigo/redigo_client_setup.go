// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package redigo

import (
	"context"
	"net"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/gomodule/redigo/redis"
)

var redigoEnabler = instrumenter.NewDefaultInstrumentEnabler()

func onBeforeDialContext(call api.CallContext, ctx context.Context, network, address string, options ...redis.DialOption) {
	if !redigoEnabler.Enable() {
		return
	}
	data := make(map[string]interface{}, 2)
	data["endpoint"] = address
	data["ctx"] = ctx
	call.SetData(data)
}

func onExitDialContext(call api.CallContext, conn redis.Conn, err error) {
	if !redigoEnabler.Enable() {
		return
	}
	d := call.GetData()
	data, ok := d.(map[string]interface{})
	if !ok {
		return
	}
	e, ok := data["endpoint"]
	if !ok {
		return
	}
	endpoint, ok := e.(string)
	if !ok {
		return
	}
	c, ok := data["ctx"]
	if !ok {
		return
	}
	ctx, ok := c.(context.Context)
	if !ok {
		return
	}
	call.SetReturnVal(0, &armsConn{conn, endpoint, ctx})
}

func onEnterDialURLContext(call api.CallContext, ctx context.Context, rawurl string, options ...redis.DialOption) {
	if !redigoEnabler.Enable() {
		return
	}
	data := make(map[string]interface{}, 2)
	data["endpoint"] = rawurl
	data["ctx"] = ctx
	call.SetData(data)
}

func onExitDialURLContext(call api.CallContext, conn redis.Conn, err error) {
	if !redigoEnabler.Enable() {
		return
	}
	d := call.GetData()
	data, ok := d.(map[string]interface{})
	if !ok {
		return
	}
	e, ok := data["endpoint"]
	if !ok {
		return
	}
	endpoint, ok := e.(string)
	if !ok {
		return
	}
	c, ok := data["ctx"]
	if !ok {
		return
	}
	ctx, ok := c.(context.Context)
	if !ok {
		return
	}
	call.SetReturnVal(0, &armsConn{conn, endpoint, ctx})
}

func onEnterNewConn(call api.CallContext, netConn net.Conn, readTimeout, writeTimeout time.Duration) {
	if !redigoEnabler.Enable() {
		return
	}
	call.SetData(netConn.RemoteAddr().String())
}

func onExitNewConn(call api.CallContext, conn redis.Conn) {
	if !redigoEnabler.Enable() {
		return
	}
	e := call.GetData()
	endpoint, ok := e.(string)
	if !ok {
		return
	}
	call.SetReturnVal(0, &armsConn{conn, endpoint, context.Background()})
}
