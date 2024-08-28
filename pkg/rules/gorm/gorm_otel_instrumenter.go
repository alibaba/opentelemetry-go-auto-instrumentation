// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//go:build ignore

package rule

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
)

type gormAttrsGetter struct {
}

func (g gormAttrsGetter) GetSystem(request gormRequest) string {
	return request.System
}

func (g gormAttrsGetter) GetUser(request gormRequest) string {
	return request.User
}

func (g gormAttrsGetter) GetName(request gormRequest) string {
	return request.DbName
}

func (g gormAttrsGetter) GetConnectionString(request gormRequest) string {
	return request.Endpoint
}

func (g gormAttrsGetter) GetStatement(request gormRequest) string {
	// full sql statement is recorded in database/sql
	return ""
}

func (g gormAttrsGetter) GetOperation(request gormRequest) string {
	return request.Operation
}

func BuildGormInstrumenter() *instrumenter.Instrumenter[gormRequest, interface{}] {
	builder := instrumenter.Builder[gormRequest, interface{}]{}
	getter := gormAttrsGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[gormRequest]{Getter: getter}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[gormRequest]{}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[gormRequest, any, gormAttrsGetter]{Base: db.DbClientCommonAttrsExtractor[gormRequest, any, gormAttrsGetter]{Getter: getter}}).
		BuildInstrumenter()
}
