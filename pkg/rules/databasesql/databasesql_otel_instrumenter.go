//go:build ignore

package databasesql

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
)

type databaseSqlAttrsGetter struct {
}

func (d databaseSqlAttrsGetter) GetSystem(request databaseSqlRequest) string {
	switch request.driverName {
	case "mysql":
		return "mysql"
	case "postgres":
		fallthrough
	case "postgresql":
		return "postgresql"
	}
	return "database"
}

func (d databaseSqlAttrsGetter) GetUser(request databaseSqlRequest) string {
	return ""
}

func (d databaseSqlAttrsGetter) GetName(request databaseSqlRequest) string {
	// TODO: parse database name from dsn
	return ""
}

func (d databaseSqlAttrsGetter) GetConnectionString(request databaseSqlRequest) string {
	return request.dsn
}

func (d databaseSqlAttrsGetter) GetStatement(request databaseSqlRequest) string {
	return request.sql
}

func (d databaseSqlAttrsGetter) GetOperation(request databaseSqlRequest) string {
	return request.opType
}

func BuildDatabaseSqlOtelInstrumenter() *instrumenter.Instrumenter[databaseSqlRequest, interface{}] {
	builder := instrumenter.Builder[databaseSqlRequest, interface{}]{}
	getter := databaseSqlAttrsGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[databaseSqlRequest]{Getter: getter}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[databaseSqlRequest]{}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[databaseSqlRequest, any, databaseSqlAttrsGetter]{Base: db.DbClientCommonAttrsExtractor[databaseSqlRequest, any, databaseSqlAttrsGetter]{Getter: getter}}).
		BuildInstrumenter()
}
