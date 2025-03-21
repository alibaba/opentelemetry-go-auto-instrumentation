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

package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	// `db.collection.name` is only available in semconv/v1.29.0 and above.
	// Considering the compatibility, we cannot upgrade to that version at this time,
	// so define a transitional constant here instead.
	// ref: https://github.com/open-telemetry/semantic-conventions/releases/tag/v1.29.0
	DBCollectionNameKey = attribute.Key("db.collection.name")
)

type DbExperimentalAttributesEnabler interface {
	Enable() bool
}

type defaultDbExperimentalAttributesEnabler struct {
}

func (d defaultDbExperimentalAttributesEnabler) Enable() bool {
	return false
}

var experimentalAttributesEnabler = instrumenter.NewDefaultInstrumentEnabler()

type DbClientCommonAttrsExtractor[REQUEST any, RESPONSE any, GETTER DbClientCommonAttrsGetter[REQUEST]] struct {
	Getter           GETTER
	AttributesFilter func(attrs []attribute.KeyValue) []attribute.KeyValue
}

func (d *DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.DB_CLIENT_KEY
}

func (d *DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attributes []attribute.KeyValue, parentContext context.Context, request REQUEST) ([]attribute.KeyValue, context.Context) {
	return attributes, parentContext
}

func (d *DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attrs []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) ([]attribute.KeyValue, context.Context) {
	attrs = append(attrs, attribute.KeyValue{
		Key:   semconv.DBSystemKey,
		Value: attribute.StringValue(d.Getter.GetSystem(request)),
	})
	if d.AttributesFilter != nil {
		attrs = d.AttributesFilter(attrs)
	}
	return attrs, context
}

type DbClientAttrsExtractor[REQUEST any, RESPONSE any, GETTER DbClientAttrsGetter[REQUEST]] struct {
	Base DbClientCommonAttrsExtractor[REQUEST, RESPONSE, GETTER]
}

func (d *DbClientAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnStart(attrs []attribute.KeyValue, parentContext context.Context, request REQUEST) ([]attribute.KeyValue, context.Context) {
	attrs, parentContext = d.Base.OnStart(attrs, parentContext, request)
	if d.Base.AttributesFilter != nil {
		attrs = d.Base.AttributesFilter(attrs)
	}
	return attrs, parentContext
}

func (d *DbClientAttrsExtractor[REQUEST, RESPONSE, GETTER]) OnEnd(attrs []attribute.KeyValue, context context.Context, request REQUEST, response RESPONSE, err error) ([]attribute.KeyValue, context.Context) {
	attrs, context = d.Base.OnEnd(attrs, context, request, response, err)
	attrs = append(attrs, attribute.KeyValue{
		Key:   semconv.DBQueryTextKey,
		Value: attribute.StringValue(d.Base.Getter.GetStatement(request)),
	}, attribute.KeyValue{
		Key:   semconv.DBOperationNameKey,
		Value: attribute.StringValue(d.Base.Getter.GetOperation(request)),
	}, attribute.KeyValue{
		Key:   semconv.ServerAddressKey,
		Value: attribute.StringValue(d.Base.Getter.GetServerAddress(request)),
	}, attribute.KeyValue{
		Key:   DBCollectionNameKey,
		Value: attribute.StringValue(d.Base.Getter.GetCollection(request)),
	})
	if d.Base.AttributesFilter != nil {
		attrs = d.Base.AttributesFilter(attrs)
	}
	if experimentalAttributesEnabler.Enable() {
		params := d.Base.Getter.GetParameters(request)
		if len(params) > 0 {
			for i, param := range params {
				attrs = append(attrs, attribute.String("db.query.parameter."+strconv.Itoa(i), fmt.Sprintf("%v", param)))
			}
		}
	}
	return attrs, context
}

func (d *DbClientAttrsExtractor[REQUEST, RESPONSE, GETTER]) GetSpanKey() attribute.Key {
	return utils.DB_CLIENT_KEY
}

// TODO: sanitize sql
// TODO: request size & response size
