// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package elasticsearch

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type elasticSearchGetter struct {
}

func (e elasticSearchGetter) GetSystem(request *esRequest) string {
	return "elasticsearch"
}

func (e elasticSearchGetter) GetServerAddress(request *esRequest) string {
	return request.address
}

func (e elasticSearchGetter) GetStatement(request *esRequest) string {
	return request.request.URL.Path
}

func (e elasticSearchGetter) GetOperation(request *esRequest) string {
	return request.op
}

func (e elasticSearchGetter) GetParameters(request *esRequest) []any {
	return request.params
}

func BuildElasticSearchInstrumenter() instrumenter.Instrumenter[*esRequest, interface{}] {
	builder := instrumenter.Builder[*esRequest, any]{}
	getter := elasticSearchGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[*esRequest]{Getter: elasticSearchGetter{}}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[*esRequest]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.ELASTICSEARCH_SCOPE_NAME,
			Version: version.Tag,
		}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[*esRequest, any, db.DbClientAttrsGetter[*esRequest]]{Base: db.DbClientCommonAttrsExtractor[*esRequest, any, db.DbClientAttrsGetter[*esRequest]]{Getter: getter}}).
		BuildInstrumenter()
}
