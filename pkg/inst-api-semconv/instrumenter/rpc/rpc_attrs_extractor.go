// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type RpcAttrsExtractor[REQUEST any, RESPONSE any, GETTER RpcAttrsGetter[REQUEST]] struct {
	Getter GETTER
}

func (r *RpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) ([]attribute.KeyValue, context.Context) {
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.RPCSystemKey,
		Value: attribute.StringValue(r.Getter.GetSystem(request)),
	}, attribute.KeyValue{
		Key:   semconv.RPCServiceKey,
		Value: attribute.StringValue(r.Getter.GetService(request)),
	}, attribute.KeyValue{
		Key:   semconv.RPCMethodKey,
		Value: attribute.StringValue(r.Getter.GetMethod(request)),
	})
	return attributes, parentContext
}

func (r *RpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) ([]attribute.KeyValue, context.Context) {
	return attributes, context
}

type ServerRpcAttrsExtractor[REQUEST any, RESPONSE any, GETTER RpcAttrsGetter[REQUEST]] struct {
	Base RpcAttrsExtractor[REQUEST, RESPONSE, GETTER]
}

func (s *ServerRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.RPC_SERVER_KEY
}

func (s *ServerRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) ([]attribute.KeyValue, context.Context) {
	return s.Base.OnStart(attributes, parentContext, request)
}

func (s *ServerRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) ([]attribute.KeyValue, context.Context) {
	return s.Base.OnEnd(attributes, context, request, response, err)
}

type ClientRpcAttrsExtractor[REQUEST any, RESPONSE any, GETTER RpcAttrsGetter[REQUEST]] struct {
	Base RpcAttrsExtractor[REQUEST, RESPONSE, GETTER]
}

func (s *ClientRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.RPC_CLIENT_KEY
}

func (s *ClientRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) ([]attribute.KeyValue, context.Context) {
	return s.Base.OnStart(attributes, parentContext, request)
}

func (s *ClientRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) ([]attribute.KeyValue, context.Context) {
	return s.Base.OnEnd(attributes, context, request, response, err)
}
