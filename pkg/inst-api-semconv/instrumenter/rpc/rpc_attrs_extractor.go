package rpc

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type RpcAttrsExtractor[REQUEST any, RESPONSE any, GETTER RpcAttrsGetter[REQUEST]] struct {
	getter GETTER
}

func (r *RpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.RPCSystemKey,
		Value: attribute.StringValue(r.getter.GetSystem(request)),
	}, attribute.KeyValue{
		Key:   semconv.RPCServiceKey,
		Value: attribute.StringValue(r.getter.GetService(request)),
	}, attribute.KeyValue{
		Key:   semconv.RPCMethodKey,
		Value: attribute.StringValue(r.getter.GetMethod(request)),
	})
	return attributes
}

func (r *RpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	return attributes
}

type ServerRpcAttrsExtractor[REQUEST any, RESPONSE any, GETTER RpcAttrsGetter[REQUEST]] struct {
	base RpcAttrsExtractor[REQUEST, RESPONSE, GETTER]
}

func (s *ServerRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.RPC_SERVER_KEY
}

func (s *ServerRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	return s.base.OnStart(attributes, parentContext, request)
}

func (s *ServerRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	return s.base.OnEnd(attributes, context, request, response, err)
}

type ClientRpcAttrsExtractor[REQUEST any, RESPONSE any, GETTER RpcAttrsGetter[REQUEST]] struct {
	base RpcAttrsExtractor[REQUEST, RESPONSE, GETTER]
}

func (s *ClientRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.RPC_CLIENT_KEY
}

func (s *ClientRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	return s.base.OnStart(attributes, parentContext, request)
}

func (s *ClientRpcAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attributes []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	return s.base.OnEnd(attributes, context, request, response, err)
}
