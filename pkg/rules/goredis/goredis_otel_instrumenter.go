//go:build ignore

package goredis

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
)

type goRedisAttrsGetter struct {
}

func (d goRedisAttrsGetter) GetSystem(request goRedisRequest) string {
	return "redis"
}

func (d goRedisAttrsGetter) GetUser(request goRedisRequest) string {
	return ""
}

func (d goRedisAttrsGetter) GetName(request goRedisRequest) string {
	// TODO: parse database name from dsn
	return ""
}

func (d goRedisAttrsGetter) GetConnectionString(request goRedisRequest) string {
	return request.endpoint
}

func (d goRedisAttrsGetter) GetStatement(request goRedisRequest) string {
	return request.cmd.String()
}

func (d goRedisAttrsGetter) GetOperation(request goRedisRequest) string {
	return request.cmd.FullName()
}

func BuildGoRedisOtelInstrumenter() *instrumenter.Instrumenter[goRedisRequest, interface{}] {
	builder := instrumenter.Builder[goRedisRequest, interface{}]{}
	getter := goRedisAttrsGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[goRedisRequest]{Getter: getter}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[goRedisRequest]{}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[goRedisRequest, any, goRedisAttrsGetter]{Base: db.DbClientCommonAttrsExtractor[goRedisRequest, any, goRedisAttrsGetter]{Getter: getter}}).
		BuildInstrumenter()
}
