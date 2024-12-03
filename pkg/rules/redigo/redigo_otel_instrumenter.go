// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package redigo

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"strings"
	"time"
)

type redigoRequest struct {
	args      []interface{}
	endpoint  string
	cmd       string
	ctx       context.Context
	startTime time.Time
}

type redigoAttrsGetter struct {
}

func (m redigoAttrsGetter) GetSystem(request *redigoRequest) string {
	return "redis"
}

func (m redigoAttrsGetter) GetServerAddress(request *redigoRequest) string {
	return request.endpoint
}

func (m redigoAttrsGetter) GetStatement(request *redigoRequest) string {
	builder := strings.Builder{}
	builder.WriteString(request.cmd + " ")
	for _, arg := range request.args {
		builder.WriteString(fmt.Sprintf("%v ", arg))
	}
	return builder.String()
}

func (m redigoAttrsGetter) GetOperation(request *redigoRequest) string {
	return request.cmd
}

func (m redigoAttrsGetter) GetParameters(request *redigoRequest) []any {
	return nil
}

func BuildRedigoInstrumenter() instrumenter.Instrumenter[*redigoRequest, interface{}] {
	builder := instrumenter.Builder[*redigoRequest, any]{}
	getter := redigoAttrsGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[*redigoRequest]{Getter: getter}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[*redigoRequest]{}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.REDIGO_SCOPE_NAME,
			Version: version.Tag,
		}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[*redigoRequest, any, db.DbClientAttrsGetter[*redigoRequest]]{Base: db.DbClientCommonAttrsExtractor[*redigoRequest, any, db.DbClientAttrsGetter[*redigoRequest]]{Getter: getter}}).
		BuildInstrumenter()
}
