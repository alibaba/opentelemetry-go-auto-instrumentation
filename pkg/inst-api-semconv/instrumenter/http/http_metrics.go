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

package http

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"time"
)

const http_server_request_duration = "http.server.request.duration"

const http_client_request_duration = "http.client.request.duration"

type HttpServerMetric struct {
	key                   attribute.Key
	serverRequestDuration metric.Float64Histogram
}

type HttpClientMetric struct {
	key                   attribute.Key
	clientRequestDuration metric.Float64Histogram
}

func NewHttpServerMetric(key string, meter metric.Meter) (*HttpServerMetric, error) {
	m := &HttpServerMetric{
		key: attribute.Key(key),
	}
	d, err := newHttpServerRequestDurationMeasures(meter)
	if err != nil {
		return nil, err
	}
	m.serverRequestDuration = d
	return m, nil
}

func newHttpServerRequestDurationMeasures(meter metric.Meter) (metric.Float64Histogram, error) {
	d, err := meter.Float64Histogram(http_server_request_duration,
		metric.WithUnit("ms"),
		metric.WithDescription("Duration of HTTP server requests."))
	if err == nil {
		return d, nil
	} else {
		return d, errors.New(fmt.Sprintf("failed to create http.server.request.duratio histogram, %v", err))
	}
}

func NewHttpClientMetric(key string, meter metric.Meter) (*HttpClientMetric, error) {
	m := &HttpClientMetric{
		key: attribute.Key(key),
	}
	d, err := newHttpClientRequestDurationMeasures(meter)
	if err != nil {
		return nil, err
	}
	m.clientRequestDuration = d
	return m, nil
}

func newHttpClientRequestDurationMeasures(meter metric.Meter) (metric.Float64Histogram, error) {
	d, err := meter.Float64Histogram(http_client_request_duration,
		metric.WithUnit("ms"),
		metric.WithDescription("Duration of HTTP client requests."))
	if err == nil {
		return d, nil
	} else {
		return d, errors.New(fmt.Sprintf("failed to create http.client.request.duratio histogram, %v", err))
	}
}

func (h *HttpServerMetric) OnBeforeStart(parentContext context.Context, startTimestamp time.Time) context.Context {
	return context.WithValue(parentContext, h.key, startTimestamp)
}

func (h *HttpServerMetric) OnBeforeEnd(context context.Context, startAttributes []attribute.KeyValue, startTimestamp time.Time) context.Context {
	return context
}

func (h *HttpServerMetric) OnAfterStart(context context.Context, endTimestamp time.Time) {
	return
}

func (h *HttpServerMetric) OnAfterEnd(context context.Context, endAttributes []attribute.KeyValue, endTimestamp time.Time) {
	startTime := context.Value(h.key).(time.Time)
	// end attributes should be shadowed by AttrsShadower
	if h.serverRequestDuration != nil {
		h.serverRequestDuration.Record(context, float64(endTimestamp.Sub(startTime)), metric.WithAttributeSet(attribute.NewSet(endAttributes...)))
	}
}

func (h HttpClientMetric) OnBeforeStart(parentContext context.Context, startTimestamp time.Time) context.Context {
	return context.WithValue(parentContext, h.key, startTimestamp)
}

func (h HttpClientMetric) OnBeforeEnd(context context.Context, startAttributes []attribute.KeyValue, startTimestamp time.Time) context.Context {
	return context
}

func (h HttpClientMetric) OnAfterStart(context context.Context, endTimestamp time.Time) {
	return
}

func (h HttpClientMetric) OnAfterEnd(context context.Context, endAttributes []attribute.KeyValue, endTimestamp time.Time) {
	startTime := context.Value(h.key).(time.Time)
	// end attributes should be shadowed by AttrsShadower
	if h.clientRequestDuration != nil {
		h.clientRequestDuration.Record(context, float64(endTimestamp.Sub(startTime)), metric.WithAttributeSet(attribute.NewSet(endAttributes...)))
	}
}
