package verifier

import (
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

func VerifyNoSqlAttributes(span tracetest.SpanStub, name, dbName, system, user, connString, statement, operation string) {
	Assert(span.SpanKind == trace.SpanKindClient, "Expect to be client span, got %d", span.SpanKind)
	Assert(span.Name == name, "Except client span name to be %s, got %s", name, span.Name)
	actualDbName := GetAttribute(span.Attributes, "db.name").AsString()
	Assert(actualDbName == dbName, "Except client db name to be %s, got %s", dbName, actualDbName)
	actualSystem := GetAttribute(span.Attributes, "db.system").AsString()
	Assert(actualSystem == system, "Except client db system to be %s, got %s", system, actualSystem)
	actualUser := GetAttribute(span.Attributes, "db.user").AsString()
	if actualUser != "" {
		Assert(actualUser == user, "Except client db user to be %s, got %s", user, actualUser)
	}
	actualConnStr := GetAttribute(span.Attributes, "db.connection_string").AsString()
	Assert(strings.Contains(actualConnStr, connString), "Except client db conn str to be %s, got %s", connString, actualConnStr)
	actualStatement := GetAttribute(span.Attributes, "db.statement").AsString()
	Assert(actualStatement == statement, "Except client db statement to be %s, got %s", statement, actualStatement)
	actualOperation := GetAttribute(span.Attributes, "db.operation").AsString()
	Assert(actualOperation == operation, "Except client db name to be %s, got %s", operation, actualOperation)
}
