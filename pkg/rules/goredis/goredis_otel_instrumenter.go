// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
