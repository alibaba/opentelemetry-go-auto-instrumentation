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

package verifier

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestWaitAndAssertTracesOneTrace(t *testing.T) {
	err := spanExporter.ExportSpans(context.Background(), []trace.ReadOnlySpan{
		testSpan{ID: "1", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "2", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "3", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "4", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer ResetTestSpans()
	WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		if len(stubs) != 1 {
			t.Fatalf("expecting 1 traces but got %d", len(stubs))
		}
	}, 1)
}

func TestWaitAndAssertTracesMultipleTrace(t *testing.T) {
	err := spanExporter.ExportSpans(context.Background(), []trace.ReadOnlySpan{
		testSpan{ID: "1", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "2", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x020},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "3", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x030},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "4", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x040},
			SpanID:  oteltrace.SpanID{0x01},
		})},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer ResetTestSpans()
	WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		if len(stubs) != 4 {
			t.Fatalf("expecting 4 traces but got %d", len(stubs))
		}
	}, 4)
}

func TestWaitAndAssertTraceLink(t *testing.T) {
	err := spanExporter.ExportSpans(context.Background(), []trace.ReadOnlySpan{
		testSpan{ID: "1", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
		testSpan{ID: "2", spanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x02},
		}), parent: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
			TraceID: oteltrace.TraceID{0x010},
			SpanID:  oteltrace.SpanID{0x01},
		})},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer ResetTestSpans()
	WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		if len(stubs) != 1 {
			t.Fatalf("expecting 1 traces but got %d", len(stubs))
		}
		if stubs[0].Snapshots()[1].Parent().SpanID() != stubs[0].Snapshots()[0].SpanContext().SpanID() {
			t.Fatalf("expecting parent span id to be equal")
		}
	}, 1)
}
