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
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/segmentio/kafka-go"
	"net"
	"time"
	_ "unsafe"
)

//go:linkname producerWriteMessagesOnEnter github.com/segmentio/kafka-go.producerWriteMessagesOnEnter
func producerWriteMessagesOnEnter(call api.CallContext, writer *kafka.Writer, ctx context.Context, messages ...kafka.Message) {
	if !kafkaEnabler.Enable() {
		return
	}

	// Create pointers to messages for instrumentation
	messagePointers := make([]*kafka.Message, len(messages))
	for i := range messages {
		// Create a copy to avoid issues with loop variable
		messageCopy := messages[i]
		messagePointers[i] = &messageCopy
	}

	// Prepare request data for instrumentation
	producerRequest := kafkaProducerReq{
		topic: writer.Topic,
		addr:  writer.Addr,
		async: writer.Async,
		msgs:  messagePointers,
	}

	// Start instrumentation and get instrumented context
	instrumentedContext := producerInstrumenter.Start(ctx, producerRequest)

	// Store data for later use in exit hook
	instrumentationData := map[string]interface{}{
		"instrumentedContext": instrumentedContext,
		"producerRequest":     producerRequest,
		"startTimestamp":      time.Now(),
	}
	call.SetData(instrumentationData)

	// Prepare message copies to maintain idempotency
	messageCopies := make([]kafka.Message, len(messagePointers))
	for i, msg := range messagePointers {
		messageCopies[i] = *msg
	}

	// Update the call parameters with copied messages
	call.SetParam(2, messageCopies)
}

//go:linkname producerWriteMessagesOnExit github.com/segmentio/kafka-go.producerWriteMessagesOnExit
func producerWriteMessagesOnExit(call api.CallContext, err error) {
	if !kafkaEnabler.Enable() {
		return
	}

	// Retrieve stored instrumentation data
	instrumentationData := call.GetData().(map[string]interface{})
	instrumentedContext := instrumentationData["instrumentedContext"].(context.Context)
	producerRequest := instrumentationData["producerRequest"].(kafkaProducerReq)

	// End instrumentation with results (includes metrics recording)
	producerInstrumenter.End(instrumentedContext, producerRequest, nil, err)
}

//go:linkname producerWriteOnEnter github.com/segmentio/kafka-go.producerWriteOnEnter
func producerWriteOnEnter(call api.CallContext, conn *kafka.Conn, msgs []kafka.Message) {
	if !kafkaEnabler.Enable() {
		return
	}

	// Extract topic from connection or messages
	topic := extractTopicFromConnection(conn)
	if topic == "" && len(msgs) > 0 {
		topic = msgs[0].Topic
	}

	instrumentationData := map[string]interface{}{
		"parentContext":   context.Background(),
		"startTimestamp":  time.Now(),
		"topic":          topic,
		"messageCount":   len(msgs),
		"async":          false, // Write operation is synchronous
	}
	call.SetData(instrumentationData)
}

//go:linkname producerWriteOnExit github.com/segmentio/kafka-go.producerWriteOnExit
func producerWriteOnExit(call api.CallContext, n int, err error) {
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
	topic := instrumentationData["topic"].(string)
	messageCount := instrumentationData["messageCount"].(int)
	
	// Create producer request
	msgs := make([]*kafka.Message, messageCount)
	producerRequest := kafkaProducerReq{
		msgs:  msgs,
		topic: topic,
		async: false,
	}

	// Record both traces and metrics
	producerInstrumenter.StartAndEnd(
		parentContext,
		producerRequest,
		nil,
		err,
		startTimestamp,
		endTimestamp,
	)
}

//go:linkname producerWriteToOnEnter github.com/segmentio/kafka-go.producerWriteToOnEnter
func producerWriteToOnEnter(call api.CallContext, conn *kafka.Conn, msgs []kafka.Message) {
	if !kafkaEnabler.Enable() {
		return
	}

	// Extract topic from connection or messages
	topic := extractTopicFromConnection(conn)
	if topic == "" && len(msgs) > 0 {
		topic = msgs[0].Topic
	}

	instrumentationData := map[string]interface{}{
		"parentContext":   context.Background(),
		"startTimestamp":  time.Now(),
		"topic":          topic,
		"messageCount":   len(msgs),
		"async":          false,
	}
	call.SetData(instrumentationData)
}

//go:linkname producerWriteToOnExit github.com/segmentio/kafka-go.producerWriteToOnExit
func producerWriteToOnExit(call api.CallContext, n int, err error) {
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
	topic := instrumentationData["topic"].(string)
	messageCount := instrumentationData["messageCount"].(int)
	
	// Create producer request
	msgs := make([]*kafka.Message, messageCount)
	producerRequest := kafkaProducerReq{
		msgs:  msgs,
		topic: topic,
		async: false,
	}

	// Record both traces and metrics
	producerInstrumenter.StartAndEnd(
		parentContext,
		producerRequest,
		nil,
		err,
		startTimestamp,
		endTimestamp,
	)
}

//go:linkname writerWriteMessagesOnEnter github.com/segmentio/kafka-go.writerWriteMessagesOnEnter
func writerWriteMessagesOnEnter(call api.CallContext, ctx context.Context, msgs []kafka.Message) {
	if !kafkaEnabler.Enable() {
		return
	}

	// Extract topic from messages
	topic := ""
	if len(msgs) > 0 {
		topic = msgs[0].Topic
	}

	instrumentationData := map[string]interface{}{
		"parentContext":   ctx,
		"startTimestamp":  time.Now(),
		"topic":          topic,
		"messageCount":   len(msgs),
		"async":          true, // Writer operations can be async
	}
	call.SetData(instrumentationData)
}

//go:linkname writerWriteMessagesOnExit github.com/segmentio/kafka-go.writerWriteMessagesOnExit
func writerWriteMessagesOnExit(call api.CallContext, err error) {
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
	topic := instrumentationData["topic"].(string)
	messageCount := instrumentationData["messageCount"].(int)
	async := instrumentationData["async"].(bool)
	
	// Create producer request
	msgs := make([]*kafka.Message, messageCount)
	producerRequest := kafkaProducerReq{
		msgs:  msgs,
		topic: topic,
		async: async,
	}

	// Record both traces and metrics
	producerInstrumenter.StartAndEnd(
		parentContext,
		producerRequest,
		nil,
		err,
		startTimestamp,
		endTimestamp,
	)
}

// Helper function to extract topic from connection
func extractTopicFromConnection(conn *kafka.Conn) string {
	if conn == nil {
		return ""
	}
	
	// Try to extract topic from connection address or metadata
	// This is a simplified implementation - in real scenarios, 
	// you might need to access internal connection state
	return ""
}

// Helper function to create address from connection
func createAddrFromConnection(conn *kafka.Conn) net.Addr {
	if conn == nil {
		return nil
	}
	
	// Return the connection's remote address
	// This is a simplified implementation
	return nil
}