package db

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

type DbClientCommonAttrsExtractor[REQUEST any, RESPONSE any, GETTER DbClientCommonAttrsGetter[REQUEST]] struct {
	getter GETTER
}

func (d *DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.DB_CLIENT_KEY
}

func (d *DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   semconv.DBNameKey,
		Value: attribute.StringValue(d.getter.GetName(request)),
	}, attribute.KeyValue{
		Key:   semconv.DBSystemKey,
		Value: attribute.StringValue(d.getter.GetSystem(request)),
	}, attribute.KeyValue{
		Key:   semconv.DBUserKey,
		Value: attribute.StringValue(d.getter.GetUser(request)),
	}, attribute.KeyValue{
		Key:   semconv.DBConnectionStringKey,
		Value: attribute.StringValue(d.getter.GetConnectionString(request)),
	})
	return attributes
}

func (d *DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attrs []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	return attrs
}

type DbClientAttrsExtractor[REQUEST any, RESPONSE any, GETTER DbClientAttrsGetter[REQUEST]] struct {
	base DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]
}

func (d *DbClientAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attrs []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attrs = d.base.OnStart(attrs, parentContext, request)
	attrs = append(attrs, attribute.KeyValue{
		Key:   semconv.DBStatementKey,
		Value: attribute.StringValue(d.base.getter.GetStatement(request)),
	}, attribute.KeyValue{
		Key:   semconv.DBOperationKey,
		Value: attribute.StringValue(d.base.getter.GetOperation(request)),
	})
	return attrs
}

func (d *DbClientAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attrs []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	attrs = d.base.OnEnd(attrs, context, request, response, err)
	return attrs
}

func (d *DbClientAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.DB_CLIENT_KEY
}

// TODO: sanitize sql
// TODO: request size & response size
