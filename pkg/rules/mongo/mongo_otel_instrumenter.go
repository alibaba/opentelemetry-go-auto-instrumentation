// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
	return request.DatabaseName
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
