package db

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
)

const db_name = attribute.Key("db.name")
const db_system = attribute.Key("db.system")
const db_user = attribute.Key("db.user")
const db_connection_string = attribute.Key("db.connection_string")
const db_statement = attribute.Key("db.statement")
const db_operation = attribute.Key("db.operation")

type DbClientCommonAttrsExtractor[REQUEST any, RESPONSE any, GETTER DbClientCommonAttrsGetter[REQUEST]] struct {
	Getter GETTER
}

func (d *DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.DB_CLIENT_KEY
}

func (d *DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attributes = append(attributes, attribute.KeyValue{
		Key:   db_name,
		Value: attribute.StringValue(d.Getter.GetName(request)),
	}, attribute.KeyValue{
		Key:   db_system,
		Value: attribute.StringValue(d.Getter.GetSystem(request)),
	}, attribute.KeyValue{
		Key:   db_user,
		Value: attribute.StringValue(d.Getter.GetUser(request)),
	}, attribute.KeyValue{
		Key:   db_connection_string,
		Value: attribute.StringValue(d.Getter.GetConnectionString(request)),
	})
	return attributes
}

func (d *DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attrs []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	return attrs
}

type DbClientAttrsExtractor[REQUEST any, RESPONSE any, GETTER DbClientAttrsGetter[REQUEST]] struct {
	Base DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]
}

func (d *DbClientAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attrs []attribute.KeyValue, parentContext context.Context, request REQUEST) []attribute.KeyValue {
	attrs = d.Base.OnStart(attrs, parentContext, request)
	attrs = append(attrs, attribute.KeyValue{
		Key:   db_statement,
		Value: attribute.StringValue(d.Base.Getter.GetStatement(request)),
	}, attribute.KeyValue{
		Key:   db_operation,
		Value: attribute.StringValue(d.Base.Getter.GetOperation(request)),
	})
	return attrs
}

func (d *DbClientAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attrs []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) []attribute.KeyValue {
	return attrs
}

func (d *DbClientAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.DB_CLIENT_KEY
}

// TODO: sanitize sql
// TODO: request size & response size
