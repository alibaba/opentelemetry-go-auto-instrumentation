package verifier

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

func TestNoSqlAttributesPass(t *testing.T) {
	VerifyNoSqlAttributes(tracetest.SpanStub{SpanKind: trace.SpanKindClient, Name: "name", Attributes: []attribute.KeyValue{
		{Key: semconv.DBNameKey, Value: attribute.StringValue("dbname")},
		{Key: semconv.DBSystemKey, Value: attribute.StringValue("system")},
		{Key: semconv.DBUserKey, Value: attribute.StringValue("user")},
		{Key: semconv.DBConnectionStringKey, Value: attribute.StringValue("connString")},
		{Key: semconv.DBStatementKey, Value: attribute.StringValue("statement")},
		{Key: semconv.DBOperationKey, Value: attribute.StringValue("operation")},
	}}, "name", "dbname", "system", "user", "connString", "statement", "operation")
}

func TestNoSqlAttributesFail(t *testing.T) {
	defer func() {
		pass := false
		if r := recover(); r != nil {
			pass = true
		}
		if !pass {
			t.Fatal("Should be recovered from panic")
		}
	}()
	VerifyNoSqlAttributes(tracetest.SpanStub{SpanKind: trace.SpanKindClient, Name: "name", Attributes: []attribute.KeyValue{
		{Key: semconv.DBNameKey, Value: attribute.StringValue("dbname")},
		{Key: semconv.DBSystemKey, Value: attribute.StringValue("system")},
		{Key: semconv.DBUserKey, Value: attribute.StringValue("user")},
		{Key: semconv.DBConnectionStringKey, Value: attribute.StringValue("connString")},
		{Key: semconv.DBStatementKey, Value: attribute.StringValue("wrong statement")},
		{Key: semconv.DBOperationKey, Value: attribute.StringValue("operation")},
	}}, "name", "dbname", "system", "user", "connString", "statement", "operation")
}
