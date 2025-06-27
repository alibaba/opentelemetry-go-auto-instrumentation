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

package test

import (
	"context"
	"testing"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/exemplar"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TestExemplarContextManager(t *testing.T) {
	ctx := context.Background()
	mgr := exemplar.GetManager()

	// Create a traced context
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(ctx, "test-operation")
	defer span.End()

	// Store context
	gid := exemplar.GetGoroutineID()
	mgr.StoreContext(gid, ctx)

	// Retrieve context
	retrievedCtx := mgr.GetContext(gid)
	retrievedSpan := trace.SpanFromContext(retrievedCtx)

	if !retrievedSpan.SpanContext().IsValid() {
		t.Error("Expected valid span context")
	}

	// Cleanup
	mgr.CleanupContext(gid)
	retrievedCtx = mgr.GetContext(gid)
	if retrievedCtx != context.Background() {
		t.Error("Expected context to be cleaned up")
	}
}

func TestMetricsRecorderWithExemplar(t *testing.T) {
	ctx := context.Background()
	
	// Create a traced context
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(ctx, "test-operation")
	defer span.End()

	// Store context for current goroutine
	gid := exemplar.GetGoroutineID()
	exemplar.GetManager().StoreContext(gid, ctx)
	defer exemplar.GetManager().CleanupContext(gid)

	// Create metrics recorder
	recorder := instrumenter.NewMetricsRecorder("test")

	// Create a histogram
	histogram, err := recorder.GetMeter().Float64Histogram("test.duration")
	if err != nil {
		t.Fatalf("Failed to create histogram: %v", err)
	}

	// Record value with exemplar
	recorder.RecordHistogramWithExemplar(histogram, 0.123)

	// In a real test, we would verify that the exemplar was created
	// by checking the metric export data
}

func TestExemplarPerformance(t *testing.T) {
	ctx := context.Background()
	tracer := otel.Tracer("test")
	recorder := instrumenter.NewMetricsRecorder("test")

	histogram, err := recorder.GetMeter().Float64Histogram("test.duration")
	if err != nil {
		t.Fatalf("Failed to create histogram: %v", err)
	}

	// Measure performance of recording with exemplars
	start := time.Now()
	for i := 0; i < 10000; i++ {
		ctx, span := tracer.Start(ctx, "test-op")
		gid := exemplar.GetGoroutineID()
		exemplar.GetManager().StoreContext(gid, ctx)
		
		recorder.RecordHistogramWithExemplar(histogram, float64(i))
		
		exemplar.GetManager().CleanupContext(gid)
		span.End()
	}
	elapsed := time.Since(start)

	t.Logf("Recorded 10000 measurements with exemplars in %v", elapsed)
	if elapsed > 1*time.Second {
		t.Errorf("Performance test took too long: %v", elapsed)
	}
}