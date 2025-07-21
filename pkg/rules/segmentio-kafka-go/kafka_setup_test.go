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
	"testing"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Mock CallContext for testing
type mockCallContext struct {
	data  interface{}
	param interface{}
}

func (m *mockCallContext) SetData(data interface{}) {
	m.data = data
}

func (m *mockCallContext) GetData() interface{} {
	return m.data
}

func (m *mockCallContext) SetParam(index int, param interface{}) {
	m.param = param
}

func TestConsumerSetupFunctions(t *testing.T) {
	// Setup test metrics provider
	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	// Enable Kafka instrumentation for testing
	kafkaEnabler = kafkaInnerEnabler{enabled: true}

	t.Run("consumerReadMessageOnEnter", func(t *testing.T) {
		call := &mockCallContext{}
		ctx := context.Background()

		// Should not panic and should set data
		consumerReadMessageOnEnter(call, nil, ctx)

		// Verify data was set
		if call.data == nil {
			t.Error("consumerReadMessageOnEnter should set instrumentation data")
		}

		data, ok := call.data.(map[string]interface{})
		if !ok {
			t.Error("Expected data to be map[string]interface{}")
		}

		if data["parentContext"] != ctx {
			t.Error("Expected parentContext to be set")
		}

		if _, ok := data["startTimestamp"].(time.Time); !ok {
			t.Error("Expected startTimestamp to be set")
		}
	})

	t.Run("consumerReadMessageOnExit", func(t *testing.T) {
		call := &mockCallContext{}
		ctx := context.Background()
		
		// Set up data as would be done by OnEnter
		call.data = map[string]interface{}{
			"parentContext":  ctx,
			"startTimestamp": time.Now(),
		}

		message := kafka.Message{
			Topic: "test-topic",
			Value: []byte("test message"),
		}

		// Should not panic
		consumerReadMessageOnExit(call, message, nil)
		consumerReadMessageOnExit(call, message, errors.New("test error"))
	})

	t.Run("readerReadMessageOnEnter", func(t *testing.T) {
		call := &mockCallContext{}
		ctx := context.Background()

		// Should not panic and should set data
		readerReadMessageOnEnter(call, ctx)

		// Verify data was set
		if call.data == nil {
			t.Error("readerReadMessageOnEnter should set instrumentation data")
		}
	})

	t.Run("readerReadMessageOnExit", func(t *testing.T) {
		call := &mockCallContext{}
		ctx := context.Background()
		
		// Set up data as would be done by OnEnter
		call.data = map[string]interface{}{
			"parentContext":  ctx,
			"startTimestamp": time.Now(),
		}

		message := kafka.Message{
			Topic: "test-topic",
			Value: []byte("test message"),
		}

		// Should not panic
		readerReadMessageOnExit(call, message, nil)
	})

	t.Run("readerCommitMessagesOnEnter", func(t *testing.T) {
		call := &mockCallContext{}
		ctx := context.Background()
		msgs := []kafka.Message{
			{Topic: "topic1", Value: []byte("msg1")},
			{Topic: "topic2", Value: []byte("msg2")},
			{Topic: "topic1", Value: []byte("msg3")},
		}

		// Should not panic and should set data
		readerCommitMessagesOnEnter(call, ctx, msgs)

		// Verify data was set
		data, ok := call.data.(map[string]interface{})
		if !ok {
			t.Error("Expected data to be map[string]interface{}")
		}

		if data["messageCount"] != 3 {
			t.Error("Expected messageCount to be 3")
		}

		topics, ok := data["topics"].([]string)
		if !ok {
			t.Error("Expected topics to be []string")
		}

		// Should extract unique topics
		if len(topics) != 2 {
			t.Errorf("Expected 2 unique topics, got %d", len(topics))
		}
	})

	t.Run("readerCommitMessagesOnExit", func(t *testing.T) {
		call := &mockCallContext{}
		ctx := context.Background()
		
		// Set up data as would be done by OnEnter
		call.data = map[string]interface{}{
			"parentContext":  ctx,
			"startTimestamp": time.Now(),
			"messageCount":   3,
			"topics":        []string{"topic1", "topic2"},
		}

		// Should not panic
		readerCommitMessagesOnExit(call, nil)
		readerCommitMessagesOnExit(call, errors.New("commit error"))
	})
}

func TestProducerSetupFunctions(t *testing.T) {
	// Setup test metrics provider
	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	// Enable Kafka instrumentation for testing
	kafkaEnabler = kafkaInnerEnabler{enabled: true}

	t.Run("producerWriteOnEnter", func(t *testing.T) {
		call := &mockCallContext{}
		msgs := []kafka.Message{
			{Topic: "test-topic", Value: []byte("msg1")},
			{Topic: "test-topic", Value: []byte("msg2")},
		}

		// Should not panic and should set data
		producerWriteOnEnter(call, nil, msgs)

		// Verify data was set
		data, ok := call.data.(map[string]interface{})
		if !ok {
			t.Error("Expected data to be map[string]interface{}")
		}

		if data["messageCount"] != 2 {
			t.Error("Expected messageCount to be 2")
		}

		if data["topic"] != "test-topic" {
			t.Error("Expected topic to be test-topic")
		}
	})

	t.Run("producerWriteOnExit", func(t *testing.T) {
		call := &mockCallContext{}
		ctx := context.Background()
		
		// Set up data as would be done by OnEnter
		call.data = map[string]interface{}{
			"parentContext":  ctx,
			"startTimestamp": time.Now(),
			"topic":         "test-topic",
			"messageCount":  2,
		}

		// Should not panic
		producerWriteOnExit(call, 2, nil)
		producerWriteOnExit(call, 0, errors.New("write error"))
	})

	t.Run("writerWriteMessagesOnEnter", func(t *testing.T) {
		call := &mockCallContext{}
		ctx := context.Background()
		msgs := []kafka.Message{
			{Topic: "test-topic", Value: []byte("msg1")},
		}

		// Should not panic and should set data
		writerWriteMessagesOnEnter(call, ctx, msgs)

		// Verify data was set
		data, ok := call.data.(map[string]interface{})
		if !ok {
			t.Error("Expected data to be map[string]interface{}")
		}

		if data["parentContext"] != ctx {
			t.Error("Expected parentContext to be set")
		}

		if data["messageCount"] != 1 {
			t.Error("Expected messageCount to be 1")
		}

		if data["topic"] != "test-topic" {
			t.Error("Expected topic to be test-topic")
		}

		if data["async"] != true {
			t.Error("Expected async to be true for writer operations")
		}
	})

	t.Run("writerWriteMessagesOnExit", func(t *testing.T) {
		call := &mockCallContext{}
		ctx := context.Background()
		
		// Set up data as would be done by OnEnter
		call.data = map[string]interface{}{
			"parentContext":  ctx,
			"startTimestamp": time.Now(),
			"topic":         "test-topic",
			"messageCount":  1,
			"async":         true,
		}

		// Should not panic
		writerWriteMessagesOnExit(call, nil)
		writerWriteMessagesOnExit(call, errors.New("write error"))
	})
}

func TestExtractTopicsFromMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []kafka.Message
		expected []string
	}{
		{
			name:     "empty messages",
			messages: []kafka.Message{},
			expected: []string{},
		},
		{
			name: "single topic",
			messages: []kafka.Message{
				{Topic: "topic1"},
				{Topic: "topic1"},
				{Topic: "topic1"},
			},
			expected: []string{"topic1"},
		},
		{
			name: "multiple topics",
			messages: []kafka.Message{
				{Topic: "topic1"},
				{Topic: "topic2"},
				{Topic: "topic1"},
				{Topic: "topic3"},
				{Topic: "topic2"},
			},
			expected: []string{"topic1", "topic2", "topic3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTopicsFromMessages(tt.messages)
			
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d topics, got %d", len(tt.expected), len(result))
			}

			// Convert to map for easier comparison (order doesn't matter)
			resultMap := make(map[string]bool)
			for _, topic := range result {
				resultMap[topic] = true
			}

			expectedMap := make(map[string]bool)
			for _, topic := range tt.expected {
				expectedMap[topic] = true
			}

			for topic := range expectedMap {
				if !resultMap[topic] {
					t.Errorf("Expected topic %s not found in result", topic)
				}
			}

			for topic := range resultMap {
				if !expectedMap[topic] {
					t.Errorf("Unexpected topic %s found in result", topic)
				}
			}
		})
	}
}

func TestKafkaEnablerBehavior(t *testing.T) {
	t.Run("disabled instrumentation", func(t *testing.T) {
		// Temporarily disable instrumentation
		originalEnabler := kafkaEnabler
		kafkaEnabler = kafkaInnerEnabler{enabled: false}
		defer func() { kafkaEnabler = originalEnabler }()

		call := &mockCallContext{}
		ctx := context.Background()

		// All functions should return early when disabled
		consumerReadMessageOnEnter(call, nil, ctx)
		if call.data != nil {
			t.Error("Should not set data when instrumentation is disabled")
		}

		readerReadMessageOnEnter(call, ctx)
		if call.data != nil {
			t.Error("Should not set data when instrumentation is disabled")
		}

		producerWriteOnEnter(call, nil, []kafka.Message{})
		if call.data != nil {
			t.Error("Should not set data when instrumentation is disabled")
		}
	})

	t.Run("enabled instrumentation", func(t *testing.T) {
		// Ensure instrumentation is enabled
		kafkaEnabler = kafkaInnerEnabler{enabled: true}

		call := &mockCallContext{}
		ctx := context.Background()

		// Functions should work when enabled
		consumerReadMessageOnEnter(call, nil, ctx)
		if call.data == nil {
			t.Error("Should set data when instrumentation is enabled")
		}
	})
}

func TestErrorHandlingInSetupFunctions(t *testing.T) {
	// Enable Kafka instrumentation for testing
	kafkaEnabler = kafkaInnerEnabler{enabled: true}

	t.Run("exit functions with missing data", func(t *testing.T) {
		call := &mockCallContext{}
		// No data set, should handle gracefully

		// These should not panic even without proper setup data
		consumerReadMessageOnExit(call, kafka.Message{}, nil)
		readerReadMessageOnExit(call, kafka.Message{}, nil)
		producerWriteOnExit(call, 0, nil)
		writerWriteMessagesOnExit(call, nil)
		readerCommitMessagesOnExit(call, nil)
	})

	t.Run("exit functions with invalid data types", func(t *testing.T) {
		call := &mockCallContext{}
		call.data = "invalid data type" // Should be map[string]interface{}

		// These should handle invalid data gracefully and not panic
		consumerReadMessageOnExit(call, kafka.Message{}, nil)
		readerReadMessageOnExit(call, kafka.Message{}, nil)
		producerWriteOnExit(call, 0, nil)
		writerWriteMessagesOnExit(call, nil)
		readerCommitMessagesOnExit(call, nil)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("extractTopicFromConnection", func(t *testing.T) {
		// Test with nil connection
		topic := extractTopicFromConnection(nil)
		if topic != "" {
			t.Error("Expected empty string for nil connection")
		}

		// Note: We can't easily test with real connection without 
		// complex setup, but the function should handle nil gracefully
	})

	t.Run("createAddrFromConnection", func(t *testing.T) {
		// Test with nil connection
		addr := createAddrFromConnection(nil)
		if addr != nil {
			t.Error("Expected nil for nil connection")
		}

		// Note: We can't easily test with real connection without 
		// complex setup, but the function should handle nil gracefully
	})
}

func TestProcessingHelperIntegration(t *testing.T) {
	// Setup test metrics provider
	exporter, err := prometheus.New()
	if err != nil {
		t.Fatalf("Failed to create prometheus exporter: %v", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	helper := NewProcessingHelper()
	ctx := context.Background()

	t.Run("end processing without start", func(t *testing.T) {
		// Should handle gracefully when no start time in context
		helper.EndProcessing(ctx, "test-topic", "test-group", nil)
	})

	t.Run("multiple start/end cycles", func(t *testing.T) {
		// Test multiple processing cycles
		for i := 0; i < 5; i++ {
			processingCtx := helper.StartProcessing(ctx)
			time.Sleep(1 * time.Millisecond) // Simulate processing
			helper.EndProcessing(processingCtx, "test-topic", "test-group", nil)
		}
	})

	t.Run("nested processing contexts", func(t *testing.T) {
		// Test nested processing contexts
		ctx1 := helper.StartProcessing(ctx)
		ctx2 := helper.StartProcessing(ctx1)
		
		time.Sleep(1 * time.Millisecond)
		
		// End in reverse order
		helper.EndProcessing(ctx2, "inner-topic", "inner-group", nil)
		helper.EndProcessing(ctx1, "outer-topic", "outer-group", nil)
	})
}