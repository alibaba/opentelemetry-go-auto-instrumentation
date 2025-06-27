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

package rules

import (
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

// MetricsExemplarRule adds exemplar support to metric recording
type MetricsExemplarRule struct {
	api.SkipRuleBase
}

func NewMetricsExemplarRule() *MetricsExemplarRule {
	return &MetricsExemplarRule{}
}

func (r *MetricsExemplarRule) ID() string {
	return "metrics_exemplar_recording"
}

func (r *MetricsExemplarRule) Version() string {
	return "v0.1.0"
}

func (r *MetricsExemplarRule) Filter(call *api.CallContext) bool {
	// Apply to metric recording functions
	return call.IsMetricRecording()
}

func (r *MetricsExemplarRule) Apply(call *api.CallContext) {
	call.OnBefore(func(ctx *api.CallContext) {
		// Inject exemplar context capture
		ctx.InjectCode(fmt.Sprintf(`
			// Capture trace context for exemplar
			if _span := trace.SpanFromContext(%s); _span.SpanContext().IsValid() {
				_gid := exemplar.GetGoroutineID()
				exemplar.GetManager().StoreContext(_gid, %s)
				defer exemplar.GetManager().CleanupContext(_gid)
			}
		`, ctx.Params[0].Name, ctx.Params[0].Name))
	})

	call.OnAfter(func(ctx *api.CallContext) {
		// Replace standard metric recording with exemplar-aware recording
		if ctx.IsHistogram() {
			ctx.ReplaceCode(`
				_recorder := instrumenter.NewMetricsRecorder("%s")
				_recorder.RecordHistogramWithExemplar(%s, %s, %s)
			`, ctx.InstrumentName, ctx.InstrumentVar, ctx.Value, ctx.Options)
		} else if ctx.IsCounter() {
			ctx.ReplaceCode(`
				_recorder := instrumenter.NewMetricsRecorder("%s")
				_recorder.RecordCounterWithExemplar(%s, %s, %s)
			`, ctx.InstrumentName, ctx.InstrumentVar, ctx.Value, ctx.Options)
		}
	})
}