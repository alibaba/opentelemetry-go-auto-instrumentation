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
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"sync"
	"time"
)

// KafkaMetrics provides OpenTelemetry metrics collection for Kafka operations
type KafkaMetrics struct {
	meter                        metric.Meter
	operationDuration           metric.Float64Histogram
	sentMessages                metric.Int64Counter
	consumedMessages            metric.Int64Counter
	processDuration             metric.Float64Histogram
	initOnce                    sync.Once
}

var (
	// Global instance to avoid multiple meter creation
	globalKafkaMetrics *KafkaMetrics
	metricsInitOnce    sync.Once
)

// GetKafkaMetrics returns the global KafkaMetrics instance
func GetKafkaMetrics() *KafkaMetrics {
	metricsInitOnce.Do(func() {
		globalKafkaMetrics = &KafkaMetrics{}
		globalKafkaMetrics.init()
	})
	return globalKafkaMetrics
}

// init initializes the OpenTelemetry metrics instruments
func (km *KafkaMetrics) init() {
	km.initOnce.Do(func() {
		// Create meter with instrumentation scope
		km.meter = otel.Meter(
			"github.com/segmentio/kafka-go",
			metric.WithInstrumentationVersion("1.0.0"),
		)

		var err error

		// messaging.client.operation.duration - Duration of messaging operation initiated by a producer or consumer client
		km.operationDuration, err = km.meter.Float64Histogram(
			"messaging.client.operation.duration",
			metric.WithDescription("Duration of messaging operation initiated by a producer or consumer client"),
			metric.WithUnit("s"),
			metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10),
		)
		if err != nil {
			fmt.Printf("Failed to create messaging.client.operation.duration metric: %v\n", err)
		}

		// messaging.client.sent.messages - Number of messages producer attempted to send to the broker
		km.sentMessages, err = km.meter.Int64Counter(
			"messaging.client.sent.messages",
			metric.WithDescription("Number of messages producer attempted to send to the broker"),
			metric.WithUnit("{message}"),
		)
		if err != nil {
			fmt.Printf("Failed to create messaging.client.sent.messages metric: %v\n", err)
		}

		// messaging.client.consumed.messages - Number of messages that were delivered to the application
		km.consumedMessages, err = km.meter.Int64Counter(
			"messaging.client.consumed.messages",
			metric.WithDescription("Number of messages that were delivered to the application"),
			metric.WithUnit("{message}"),
		)
		if err != nil {
			fmt.Printf("Failed to create messaging.client.consumed.messages metric: %v\n", err)
		}

		// messaging.process.duration - Duration of processing operation
		km.processDuration, err = km.meter.Float64Histogram(
			"messaging.process.duration",
			metric.WithDescription("Duration of processing operation"),
			metric.WithUnit("s"),
			metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10),
		)
		if err != nil {
			fmt.Printf("Failed to create messaging.process.duration metric: %v\n", err)
		}
	})
}

// RecordOperationDuration records the duration of a messaging operation
func (km *KafkaMetrics) RecordOperationDuration(ctx context.Context, duration time.Duration, operation string, topic string, err error) {
	if km.operationDuration == nil {
		return
	}

	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKafka,
		semconv.MessagingOperationName(operation),
		semconv.MessagingDestinationName(topic),
	}

	// Add error type if operation failed
	if err != nil {
		attrs = append(attrs, semconv.ErrorTypeKey.String(err.Error()))
	}

	km.operationDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// RecordSentMessages records the number of messages sent by producer
func (km *KafkaMetrics) RecordSentMessages(ctx context.Context, count int64, operation string, topic string, err error) {
	if km.sentMessages == nil {
		return
	}

	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKafka,
		semconv.MessagingOperationName(operation),
		semconv.MessagingDestinationName(topic),
	}

	// Add error type if operation failed
	if err != nil {
		attrs = append(attrs, semconv.ErrorTypeKey.String(err.Error()))
	}

	km.sentMessages.Add(ctx, count, metric.WithAttributes(attrs...))
}

// RecordConsumedMessages records the number of messages consumed
func (km *KafkaMetrics) RecordConsumedMessages(ctx context.Context, count int64, operation string, topic string, consumerGroup string, err error) {
	if km.consumedMessages == nil {
		return
	}

	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKafka,
		semconv.MessagingOperationName(operation),
		semconv.MessagingDestinationName(topic),
	}

	// Add consumer group if available
	if consumerGroup != "" {
		attrs = append(attrs, semconv.MessagingConsumerGroupName(consumerGroup))
	}

	// Add error type if operation failed
	if err != nil {
		attrs = append(attrs, semconv.ErrorTypeKey.String(err.Error()))
	}

	km.consumedMessages.Add(ctx, count, metric.WithAttributes(attrs...))
}

// RecordProcessDuration records the duration of message processing
func (km *KafkaMetrics) RecordProcessDuration(ctx context.Context, duration time.Duration, operation string, topic string, consumerGroup string, err error) {
	if km.processDuration == nil {
		return
	}

	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKafka,
		semconv.MessagingOperationName(operation),
		semconv.MessagingDestinationName(topic),
	}

	// Add consumer group if available
	if consumerGroup != "" {
		attrs = append(attrs, semconv.MessagingConsumerGroupName(consumerGroup))
	}

	// Add error type if operation failed
	if err != nil {
		attrs = append(attrs, semconv.ErrorTypeKey.String(err.Error()))
	}

	km.processDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// Helper functions for common metric recording patterns

// RecordProducerSend records metrics for producer send operations
func (km *KafkaMetrics) RecordProducerSend(ctx context.Context, topic string, messageCount int64, duration time.Duration, err error) {
	// Record operation duration
	km.RecordOperationDuration(ctx, duration, "send", topic, err)
	
	// Record sent messages
	km.RecordSentMessages(ctx, messageCount, "send", topic, err)
}

// RecordConsumerReceive records metrics for consumer receive operations
func (km *KafkaMetrics) RecordConsumerReceive(ctx context.Context, topic string, consumerGroup string, messageCount int64, duration time.Duration, err error) {
	// Record operation duration
	km.RecordOperationDuration(ctx, duration, "receive", topic, err)
	
	// Record consumed messages
	km.RecordConsumedMessages(ctx, messageCount, "receive", topic, consumerGroup, err)
}

// RecordMessageProcess records metrics for message processing operations
func (km *KafkaMetrics) RecordMessageProcess(ctx context.Context, topic string, consumerGroup string, duration time.Duration, err error) {
	// Record process duration
	km.RecordProcessDuration(ctx, duration, "process", topic, consumerGroup, err)
}