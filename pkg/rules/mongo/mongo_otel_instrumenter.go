//go:build ignore

package mongo

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"go.opentelemetry.io/otel/trace"
)

type mongoAttrsGetter struct {
}

func (m mongoAttrsGetter) GetSystem(request mongoRequest) string {
	return "mongodb"
}

func (m mongoAttrsGetter) GetUser(request mongoRequest) string {
	return ""
}

func (m mongoAttrsGetter) GetName(request mongoRequest) string {
	return request.CommandName
}

func (m mongoAttrsGetter) GetConnectionString(request mongoRequest) string {
	return request.ConnectionID
}

func (m mongoAttrsGetter) GetStatement(request mongoRequest) string {
	return request.CommandName
}

func (m mongoAttrsGetter) GetOperation(request mongoRequest) string {
	return request.CommandName
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

func BuildMongoOtelInstrumenter() *instrumenter.Instrumenter[mongoRequest, interface{}] {
	builder := instrumenter.Builder[mongoRequest, interface{}]{}
	return builder.Init().SetSpanNameExtractor(&mongoSpanNameExtractor{}).SetSpanKindExtractor(&mongoSpanKindExtractor{}).AddAttributesExtractor(&db.DbClientAttrsExtractor[mongoRequest, any, mongoAttrsGetter]{}).BuildInstrumenter()
}
