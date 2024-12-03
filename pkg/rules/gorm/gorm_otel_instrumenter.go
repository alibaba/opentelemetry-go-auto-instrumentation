// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package gorm

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type gormAttrsGetter struct {
}

func (g gormAttrsGetter) GetSystem(gormRequest gormRequest) string {
	return gormRequest.System
}

func (g gormAttrsGetter) GetServerAddress(gormRequest gormRequest) string {
	return gormRequest.Endpoint
}

func (g gormAttrsGetter) GetStatement(gormRequest gormRequest) string {
	return ""
}

func (g gormAttrsGetter) GetOperation(gormRequest gormRequest) string {
	return gormRequest.Operation
}

func (g gormAttrsGetter) GetParameters(gormRequest gormRequest) []any {
	return nil
}

func BuildGormInstrumenter() instrumenter.Instrumenter[gormRequest, interface{}] {
	builder := instrumenter.Builder[gormRequest, interface{}]{}
	getter := gormAttrsGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[gormRequest]{Getter: getter}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[gormRequest]{}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[gormRequest, any, gormAttrsGetter]{Base: db.DbClientCommonAttrsExtractor[gormRequest, any, gormAttrsGetter]{Getter: getter}}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.GORM_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildInstrumenter()
}
