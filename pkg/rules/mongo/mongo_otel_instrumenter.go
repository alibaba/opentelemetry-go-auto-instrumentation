// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package mongo

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/trace"
)

type mongoAttrsGetter struct {
}

func (m mongoAttrsGetter) GetSystem(request mongoRequest) string {
	return "mongodb"
}

func (m mongoAttrsGetter) GetServerAddress(request mongoRequest) string {
	return request.Host
}

func (m mongoAttrsGetter) GetStatement(request mongoRequest) string {
	return request.CommandName
}

func (m mongoAttrsGetter) GetOperation(request mongoRequest) string {
	return request.CommandName
}

func (m mongoAttrsGetter) GetParameters(request mongoRequest) []any {
	return nil
}

type mongoSpanNameExtractor struct {
}

type mongoSpanKindExtractor struct {
}

func (m *mongoSpanKindExtractor) Extract(request mongoRequest) trace.SpanKind {
	return trace.SpanKindClient
}

func (m *mongoSpanNameExtractor) Extract(request mongoRequest) string {
	return request.CommandName
}

func BuildMongoOtelInstrumenter() instrumenter.Instrumenter[mongoRequest, interface{}] {
	builder := instrumenter.Builder[mongoRequest, interface{}]{}
	return builder.Init().SetSpanNameExtractor(&mongoSpanNameExtractor{}).
		SetSpanKindExtractor(&mongoSpanKindExtractor{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.MONGO_SCOPE_NAME,
			Version: version.Tag,
		}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[mongoRequest, any, db.DbClientAttrsGetter[mongoRequest]]{Base: db.DbClientCommonAttrsExtractor[mongoRequest, any, db.DbClientAttrsGetter[mongoRequest]]{Getter: mongoAttrsGetter{}}}).
		BuildInstrumenter()
}
