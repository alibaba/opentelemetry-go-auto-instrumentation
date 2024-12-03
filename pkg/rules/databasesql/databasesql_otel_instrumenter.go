// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package databasesql

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/sdk/instrumentation"
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

func (d databaseSqlAttrsGetter) GetServerAddress(request databaseSqlRequest) string {
	return request.endpoint
}

func (d databaseSqlAttrsGetter) GetStatement(request databaseSqlRequest) string {
	return request.sql
}

func (d databaseSqlAttrsGetter) GetOperation(request databaseSqlRequest) string {
	return request.opType
}

func (d databaseSqlAttrsGetter) GetParameters(request databaseSqlRequest) []any {
	return request.params
}

func BuildDatabaseSqlOtelInstrumenter() instrumenter.Instrumenter[databaseSqlRequest, any] {
	builder := instrumenter.Builder[databaseSqlRequest, any]{}
	getter := databaseSqlAttrsGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[databaseSqlRequest]{Getter: getter}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[databaseSqlRequest]{}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[databaseSqlRequest, any, db.DbClientAttrsGetter[databaseSqlRequest]]{Base: db.DbClientCommonAttrsExtractor[databaseSqlRequest, any, db.DbClientAttrsGetter[databaseSqlRequest]]{Getter: getter}}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.DATABASE_SQL_SCOPE_NAME,
			Version: version.Tag,
		}).BuildInstrumenter()
}
