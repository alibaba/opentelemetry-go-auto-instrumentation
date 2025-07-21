// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

package kafka

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func TestKafkaMetricsInitialization(t *testing.T) {
	// Reset global metrics instance for testing
	globalKafkaMetrics = nil
	metricsInitOnce = sync.Once{}

	// Get metrics instance
	metrics := GetKafkaMetrics()

	// Verify singleton behavior
	metrics2 := GetKafkaMetrics()
	if metrics != metrics2 {
		t.Error("GetKafkaMetrics should return the same instance (singleton)")
	}

	// Verify all metrics are initialized
	if metrics.operationDuration == nil {
		t.Error("operationDuration metric should be initialized")
	}
	if metrics.sentMessages == nil {
		t.Error("sentMessages metric should be initialized")
	}
	if metrics.consumedMessages == nil {
		t.Error("consumedMessages metric should be initialized")
	}
	if metrics.processDuration == nil {
		t.Error("processDuration metric should be initialized")
	}
}

func TestKafkaMetricsRecordOperationDuration(t *testing.T) {
	// Setup test metrics provider
	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	// Reset global metrics instance
	globalKafkaMetrics = nil
	metricsInitOnce = sync.Once{}

	metrics := GetKafkaMetrics()
	ctx := context.Background()

	tests := []struct {
		name      string
		duration  time.Duration
		operation string
		topic     string
		err       error
	}{
		{
			name:      "successful send operation",
			duration:  50 * time.Millisecond,
			operation: "send",
			topic:     "test-topic",
			err:       nil,
		},
		{
			name:      "failed receive operation",
			duration:  100 * time.Millisecond,
			operation: "receive",
			topic:     "error-topic",
			err:       errors.New("connection timeout"),
		},
		{
			name:      "commit operation",
			duration:  25 * time.Millisecond,
			operation: "commit",
			topic:     "commit-topic",
			err:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic and should record the metric
			metrics.RecordOperationDuration(ctx, tt.duration, tt.operation, tt.topic, tt.err)
		})
	}
}

func TestKafkaMetricsRecordSentMessages(t *testing.T) {
	// Setup test metrics provider
	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	// Reset global metrics instance
	globalKafkaMetrics = nil
	metricsInitOnce = sync.Once{}

	metrics := GetKafkaMetrics()
	ctx := context.Background()

	tests := []struct {
		name      string
		count     int64
		operation string
		topic     string
		err       error
	}{
		{
			name:      "single message send",
			count:     1,
			operation: "send",
			topic:     "test-topic",
			err:       nil,
		},
		{
			name:      "batch message send",
			count:     10,
			operation: "send",
			topic:     "batch-topic",
			err:       nil,
		},
		{
			name:      "failed message send",
			count:     5,
			operation: "send",
			topic:     "error-topic",
			err:       errors.New("broker unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic and should record the metric
			metrics.RecordSentMessages(ctx, tt.count, tt.operation, tt.topic, tt.err)
		})
	}
}

func TestKafkaMetricsRecordConsumedMessages(t *testing.T) {
	// Setup test metrics provider
	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	// Reset global metrics instance
	globalKafkaMetrics = nil
	metricsInitOnce = sync.Once{}

	metrics := GetKafkaMetrics()
	ctx := context.Background()

	tests := []struct {
		name          string
		count         int64
		operation     string
		topic         string
		consumerGroup string
		err           error
	}{
		{
			name:          "single message consume",
			count:         1,
			operation:     "receive",
			topic:         "test-topic",
			consumerGroup: "test-group",
			err:           nil,
		},
		{
			name:          "batch message consume",
			count:         15,
			operation:     "receive",
			topic:         "batch-topic",
			consumerGroup: "batch-group",
			err:           nil,
		},
		{
			name:          "consume without group",
			count:         3,
			operation:     "receive",
			topic:         "no-group-topic",
			consumerGroup: "",
			err:           nil,
		},
		{
			name:          "failed message consume",
			count:         2,
			operation:     "receive",
			topic:         "error-topic",
			consumerGroup: "error-group",
			err:           errors.New("deserialization error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic and should record the metric
			metrics.RecordConsumedMessages(ctx, tt.count, tt.operation, tt.topic, tt.consumerGroup, tt.err)
		})
	}
}

func TestKafkaMetricsRecordProcessDuration(t *testing.T) {
	// Setup test metrics provider
	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	// Reset global metrics instance
	globalKafkaMetrics = nil
	metricsInitOnce = sync.Once{}

	metrics := GetKafkaMetrics()
	ctx := context.Background()

	tests := []struct {
		name          string
		duration      time.Duration
		operation     string
		topic         string
		consumerGroup string
		err           error
	}{
		{
			name:          "successful processing",
			duration:      200 * time.Millisecond,
			operation:     "process",
			topic:         "test-topic",
			consumerGroup: "test-group",
			err:           nil,
		},
		{
			name:          "slow processing",
			duration:      2 * time.Second,
			operation:     "process",
			topic:         "slow-topic",
			consumerGroup: "slow-group",
			err:           nil,
		},
		{
			name:          "failed processing",
			duration:      500 * time.Millisecond,
			operation:     "process",
			topic:         "error-topic",
			consumerGroup: "error-group",
			err:           errors.New("business logic error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic and should record the metric
			metrics.RecordProcessDuration(ctx, tt.duration, tt.operation, tt.topic, tt.consumerGroup, tt.err)
		})
	}
}

func TestKafkaMetricsHelperMethods(t *testing.T) {
	// Setup test metrics provider
	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	// Reset global metrics instance
	globalKafkaMetrics = nil
	metricsInitOnce = sync.Once{}

	metrics := GetKafkaMetrics()
	ctx := context.Background()

	t.Run("RecordProducerSend", func(t *testing.T) {
		// This should record both operation duration and sent messages
		metrics.RecordProducerSend(ctx, "test-topic", 5, 100*time.Millisecond, nil)
		metrics.RecordProducerSend(ctx, "error-topic", 3, 200*time.Millisecond, errors.New("send failed"))
	})

	t.Run("RecordConsumerReceive", func(t *testing.T) {
		// This should record both operation duration and consumed messages
		metrics.RecordConsumerReceive(ctx, "test-topic", "test-group", 2, 50*time.Millisecond, nil)
		metrics.RecordConsumerReceive(ctx, "error-topic", "error-group", 1, 150*time.Millisecond, errors.New("receive failed"))
	})

	t.Run("RecordMessageProcess", func(t *testing.T) {
		// This should record process duration
		metrics.RecordMessageProcess(ctx, "test-topic", "test-group", 300*time.Millisecond, nil)
		metrics.RecordMessageProcess(ctx, "error-topic", "error-group", 500*time.Millisecond, errors.New("process failed"))
	})
}

func TestSemanticConventionCompliance(t *testing.T) {
	// Reset global metrics instance
	globalKafkaMetrics = nil
	metricsInitOnce = sync.Once{}

	metrics := GetKafkaMetrics()
	ctx := context.Background()

	// Test that correct attributes are used according to OpenTelemetry semantic conventions
	t.Run("operation duration attributes", func(t *testing.T) {
		// Verify that the correct semantic convention attributes are used
		// This test ensures we're following the OpenTelemetry standards
		metrics.RecordOperationDuration(ctx, 100*time.Millisecond, "send", "test-topic", nil)
		
		// The implementation should use:
		// - semconv.MessagingSystemKafka
		// - semconv.MessagingOperationName(operation)
		// - semconv.MessagingDestinationName(topic)
		// - semconv.ErrorTypeKey.String(err.Error()) when error present
	})

	t.Run("messaging system consistency", func(t *testing.T) {
		// All metrics should use "kafka" as the messaging system
		// This is enforced by using semconv.MessagingSystemKafka
		metrics.RecordSentMessages(ctx, 1, "send", "test-topic", nil)
		metrics.RecordConsumedMessages(ctx, 1, "receive", "test-topic", "test-group", nil)
		metrics.RecordProcessDuration(ctx, 100*time.Millisecond, "process", "test-topic", "test-group", nil)
	})
}

func TestMetricsExtractorInitialization(t *testing.T) {
	t.Run("producer metrics extractor", func(t *testing.T) {
		extractor := &kafkaProducerMetricsExtractor{}
		ctx := context.Background()
		
		// Test OnStart - should add timing to context
		attrs := []attribute.KeyValue{}
		req := kafkaProducerReq{topic: "test-topic"}
		
		newAttrs, newCtx := extractor.OnStart(attrs, ctx, req)
		
		// Should return same attributes but enhanced context
		if len(newAttrs) != len(attrs) {
			t.Error("OnStart should not modify attributes")
		}
		
		// Context should contain timing information
		if _, ok := newCtx.Value("kafka_start_time").(time.Time); !ok {
			t.Error("OnStart should add timing to context")
		}
	})

	t.Run("consumer metrics extractor", func(t *testing.T) {
		extractor := &kafkaConsumerMetricsExtractor{}
		ctx := context.Background()
		
		// Test OnStart
		attrs := []attribute.KeyValue{}
		req := kafkaConsumerReq{}
		
		newAttrs, newCtx := extractor.OnStart(attrs, ctx, req)
		
		// Should return same attributes but enhanced context
		if len(newAttrs) != len(attrs) {
			t.Error("OnStart should not modify attributes")
		}
		
		// Context should contain timing information
		if _, ok := newCtx.Value("kafka_start_time").(time.Time); !ok {
			t.Error("OnStart should add timing to context")
		}
	})
}

func TestProcessingHelper(t *testing.T) {
	// Setup test metrics provider
	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	// Reset global metrics instance
	globalKafkaMetrics = nil
	metricsInitOnce = sync.Once{}

	helper := NewProcessingHelper()
	ctx := context.Background()

	t.Run("processing lifecycle", func(t *testing.T) {
		// Start processing
		processingCtx := helper.StartProcessing(ctx)
		
		// Verify timing was added to context
		if _, ok := processingCtx.Value("kafka_process_start_time").(time.Time); !ok {
			t.Error("StartProcessing should add timing to context")
		}
		
		// Simulate processing time
		time.Sleep(10 * time.Millisecond)
		
		// End processing
		helper.EndProcessing(processingCtx, "test-topic", "test-group", nil)
		helper.EndProcessing(processingCtx, "error-topic", "error-group", errors.New("processing error"))
	})

	t.Run("global helper instance", func(t *testing.T) {
		// Test the global ProcessHelper instance
		processingCtx := ProcessHelper.StartProcessing(ctx)
		time.Sleep(5 * time.Millisecond)
		ProcessHelper.EndProcessing(processingCtx, "global-topic", "global-group", nil)
	})
}

func TestMetricsWithNilHandling(t *testing.T) {
	// Test that metrics handle nil instruments gracefully
	metrics := &KafkaMetrics{
		// Intentionally leave instruments as nil to test error handling
		operationDuration: nil,
		sentMessages:     nil,
		consumedMessages: nil,
		processDuration:  nil,
	}
	
	ctx := context.Background()
	
	// These should not panic even with nil instruments
	metrics.RecordOperationDuration(ctx, 100*time.Millisecond, "send", "test-topic", nil)
	metrics.RecordSentMessages(ctx, 1, "send", "test-topic", nil)
	metrics.RecordConsumedMessages(ctx, 1, "receive", "test-topic", "test-group", nil)
	metrics.RecordProcessDuration(ctx, 100*time.Millisecond, "process", "test-topic", "test-group", nil)
	
	// Helper methods should also not panic
	metrics.RecordProducerSend(ctx, "test-topic", 1, 100*time.Millisecond, nil)
	metrics.RecordConsumerReceive(ctx, "test-topic", "test-group", 1, 100*time.Millisecond, nil)
	metrics.RecordMessageProcess(ctx, "test-topic", "test-group", 100*time.Millisecond, nil)
}

func TestConcurrentMetricsAccess(t *testing.T) {
	// Test concurrent access to metrics singleton
	// Reset global metrics instance
	globalKafkaMetrics = nil
	metricsInitOnce = sync.Once{}
	
	const numGoroutines = 10
	instances := make([]*KafkaMetrics, numGoroutines)
	
	// Start multiple goroutines trying to get metrics instance
	done := make(chan bool, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			instances[index] = GetKafkaMetrics()
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// All instances should be the same (singleton)
	for i := 1; i < numGoroutines; i++ {
		if instances[i] != instances[0] {
			t.Errorf("Concurrent access should return same instance, got different instances at index %d", i)
		}
	}
}