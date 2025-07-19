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
	"github.com/alibaba/loongsuite-go-agent/pkg/api"
	"github.com/segmentio/kafka-go"
	"time"
	_ "unsafe"
)

//go:linkname consumerReadMessageOnEnter github.com/segmentio/kafka-go.consumerReadMessageOnEnter
func consumerReadMessageOnEnter(call api.CallContext, _ interface{}, ctx context.Context) {
	if !kafkaEnabler.Enable() {
		return
	}

	instrumentationData := map[string]interface{}{
		"parentContext":   ctx,
		"startTimestamp":  time.Now(),
	}
	call.SetData(instrumentationData)
}

//go:linkname consumerReadMessageOnExit github.com/segmentio/kafka-go.consumerReadMessageOnExit
func consumerReadMessageOnExit(call api.CallContext, message kafka.Message, err error) {
	if !kafkaEnabler.Enable() {
		return
	}

	instrumentationData, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}

	parentContext := instrumentationData["parentContext"].(context.Context)
	startTimestamp := instrumentationData["startTimestamp"].(time.Time)
	endTimestamp := time.Now()

	consumerRequest := kafkaConsumerReq{msg: message}
	
	// Record both traces and metrics
	consumerInstrumenter.StartAndEnd(
		parentContext,
		consumerRequest,
		nil,
		err,
		startTimestamp,
		endTimestamp,
	)
}

//go:linkname readerReadMessageOnEnter github.com/segmentio/kafka-go.readerReadMessageOnEnter
func readerReadMessageOnEnter(call api.CallContext, ctx context.Context) {
	if !kafkaEnabler.Enable() {
		return
	}

	instrumentationData := map[string]interface{}{
		"parentContext":   ctx,
		"startTimestamp":  time.Now(),
	}
	call.SetData(instrumentationData)
}

//go:linkname readerReadMessageOnExit github.com/segmentio/kafka-go.readerReadMessageOnExit
func readerReadMessageOnExit(call api.CallContext, message kafka.Message, err error) {
	if !kafkaEnabler.Enable() {
		return
	}

	instrumentationData, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}

	parentContext := instrumentationData["parentContext"].(context.Context)
	startTimestamp := instrumentationData["startTimestamp"].(time.Time)
	endTimestamp := time.Now()

	consumerRequest := kafkaConsumerReq{msg: message}
	
	// Record both traces and metrics
	consumerInstrumenter.StartAndEnd(
		parentContext,
		consumerRequest,
		nil,
		err,
		startTimestamp,
		endTimestamp,
	)
}

//go:linkname readerFetchMessageOnEnter github.com/segmentio/kafka-go.readerFetchMessageOnEnter
func readerFetchMessageOnEnter(call api.CallContext, ctx context.Context) {
	if !kafkaEnabler.Enable() {
		return
	}

	instrumentationData := map[string]interface{}{
		"parentContext":   ctx,
		"startTimestamp":  time.Now(),
	}
	call.SetData(instrumentationData)
}

//go:linkname readerFetchMessageOnExit github.com/segmentio/kafka-go.readerFetchMessageOnExit
func readerFetchMessageOnExit(call api.CallContext, message kafka.Message, err error) {
	if !kafkaEnabler.Enable() {
		return
	}

	instrumentationData, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}

	parentContext := instrumentationData["parentContext"].(context.Context)
	startTimestamp := instrumentationData["startTimestamp"].(time.Time)
	endTimestamp := time.Now()

	consumerRequest := kafkaConsumerReq{msg: message}
	
	// Record both traces and metrics
	consumerInstrumenter.StartAndEnd(
		parentContext,
		consumerRequest,
		nil,
		err,
		startTimestamp,
		endTimestamp,
	)
}

//go:linkname readerCommitMessagesOnEnter github.com/segmentio/kafka-go.readerCommitMessagesOnEnter
func readerCommitMessagesOnEnter(call api.CallContext, ctx context.Context, msgs []kafka.Message) {
	if !kafkaEnabler.Enable() {
		return
	}

	instrumentationData := map[string]interface{}{
		"parentContext":   ctx,
		"startTimestamp":  time.Now(),
		"messageCount":    len(msgs),
		"topics":         extractTopicsFromMessages(msgs),
	}
	call.SetData(instrumentationData)
}

//go:linkname readerCommitMessagesOnExit github.com/segmentio/kafka-go.readerCommitMessagesOnExit
func readerCommitMessagesOnExit(call api.CallContext, err error) {
	if !kafkaEnabler.Enable() {
		return
	}

	instrumentationData, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}

	parentContext := instrumentationData["parentContext"].(context.Context)
	startTimestamp := instrumentationData["startTimestamp"].(time.Time)
	endTimestamp := time.Now()
	messageCount := instrumentationData["messageCount"].(int)
	topics := instrumentationData["topics"].([]string)

	// Record metrics for commit operation
	duration := endTimestamp.Sub(startTimestamp)
	for _, topic := range topics {
		kafkaMetrics.RecordOperationDuration(parentContext, duration, "commit", topic, err)
	}
}

// Helper function to extract unique topics from messages
func extractTopicsFromMessages(msgs []kafka.Message) []string {
	topicSet := make(map[string]bool)
	for _, msg := range msgs {
		topicSet[msg.Topic] = true
	}
	
	topics := make([]string, 0, len(topicSet))
	for topic := range topicSet {
		topics = append(topics, topic)
	}
	return topics
}

// ProcessingHelper provides utilities for message processing metrics
type ProcessingHelper struct {
	metrics *KafkaMetrics
}

// NewProcessingHelper creates a new processing helper
func NewProcessingHelper() *ProcessingHelper {
	return &ProcessingHelper{
		metrics: GetKafkaMetrics(),
	}
}

// StartProcessing marks the beginning of message processing and returns context with timing info
func (h *ProcessingHelper) StartProcessing(ctx context.Context) context.Context {
	return context.WithValue(ctx, "kafka_process_start_time", time.Now())
}

// EndProcessing marks the end of message processing and records metrics
func (h *ProcessingHelper) EndProcessing(ctx context.Context, topic string, consumerGroup string, err error) {
	if startTime, ok := ctx.Value("kafka_process_start_time").(time.Time); ok {
		duration := time.Since(startTime)
		h.metrics.RecordMessageProcess(ctx, topic, consumerGroup, duration, err)
	}
}

// Global helper instance for convenience
var ProcessHelper = NewProcessingHelper()