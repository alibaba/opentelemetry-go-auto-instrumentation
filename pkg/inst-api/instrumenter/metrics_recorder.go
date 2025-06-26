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

package instrumenter

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/exemplar"
)

// MetricsRecorder wraps metric instruments with exemplar support
type MetricsRecorder struct {
	meter      metric.Meter
	contextMgr *exemplar.ContextManager
}

// NewMetricsRecorder creates a new metrics recorder with exemplar support
func NewMetricsRecorder(name string) *MetricsRecorder {
	return &MetricsRecorder{
		meter:      otel.Meter(name),
		contextMgr: exemplar.GetManager(),
	}
}

// RecordHistogramWithExemplar records a histogram value with potential exemplar
func (r *MetricsRecorder) RecordHistogramWithExemplar(
	instrument metric.Float64Histogram,
	value float64,
	opts ...metric.RecordOption,
) {
	gid := exemplar.GetGoroutineID()
	ctx := r.contextMgr.GetContext(gid)
	instrument.Record(ctx, value, opts...)
}

// RecordCounterWithExemplar records a counter value with potential exemplar
func (r *MetricsRecorder) RecordCounterWithExemplar(
	instrument metric.Int64Counter,
	value int64,
	opts ...metric.AddOption,
) {
	gid := exemplar.GetGoroutineID()
	ctx := r.contextMgr.GetContext(gid)
	instrument.Add(ctx, value, opts...)
}

// RecordUpDownCounterWithExemplar records an up-down counter value with potential exemplar
func (r *MetricsRecorder) RecordUpDownCounterWithExemplar(
	instrument metric.Int64UpDownCounter,
	value int64,
	opts ...metric.AddOption,
) {
	gid := exemplar.GetGoroutineID()
	ctx := r.contextMgr.GetContext(gid)
	instrument.Add(ctx, value, opts...)
}

// GetMeter returns the underlying meter for creating instruments
func (r *MetricsRecorder) GetMeter() metric.Meter {
	return r.meter
}
