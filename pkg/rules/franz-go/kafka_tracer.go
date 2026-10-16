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

package franz_go

import (
	"context"
	"strconv"

	"github.com/twmb/franz-go/pkg/kgo"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	_ kgo.HookProduceRecordBuffered   = new(KafkaTracer)
	_ kgo.HookProduceRecordUnbuffered = new(KafkaTracer)
	_ kgo.HookFetchRecordBuffered     = new(KafkaTracer)
	_ kgo.HookFetchRecordUnbuffered   = new(KafkaTracer)
)

type KafkaTracer struct{}

// detachSpanFromGLS removes the span created by Start from the caller's
// goroutine-local TraceContext while keeping it alive in ctx.
//
// Produce buffered hooks run on the Produce caller; completion hooks can run
// on another goroutine. Fetch buffered hooks run while records are fetched,
// while unbuffered hooks run when records are polled or discarded. Start/End
// therefore need not access the same GLS stack. Detach where Start runs, and
// keep the span in Record.Context until the corresponding unbuffered hook.
func detachSpanFromGLS(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}
	sdktrace.TraceContextDelSpan(span)
}

func (t *KafkaTracer) OnProduceRecordBuffered(r *kgo.Record) {
	if r.Context == nil {
		r.Context = context.Background()
	}
	request := franzKafkaProducerRequest{
		topic:   r.Topic,
		Record:  r,
		Context: r.Context,
	}
	ctx := franzKafkaProducerInstrumenter.Start(r.Context, request)
	r.Context = ctx
	// Inject already happened in Start. Detach from caller GLS so the async
	// Unbuffered End does not leak the span on the Produce caller.
	detachSpanFromGLS(ctx)
}

func (t *KafkaTracer) OnProduceRecordUnbuffered(r *kgo.Record, err error) {
	request := franzKafkaProducerRequest{
		topic:   r.Topic,
		Record:  r,
		Context: r.Context,
	}
	offsetStr := strconv.FormatInt(r.Offset, 10)
	response := franzKafkaProducerResponse{
		offset: &offsetStr,
	}
	franzKafkaProducerInstrumenter.End(r.Context, request, response, err)
}

func (t *KafkaTracer) OnFetchRecordBuffered(r *kgo.Record) {
	if r.Context == nil {
		r.Context = context.Background()
	}
	request := franzKafkaConsumerRequest{
		topic:   r.Topic,
		Record:  r,
		Context: r.Context,
	}
	newCtx := franzKafkaConsumerInstrumenter.Start(r.Context, request)
	r.Context = newCtx
	// Detach from the fetch goroutine before polling can end the span elsewhere.
	detachSpanFromGLS(newCtx)
}

func (t *KafkaTracer) OnFetchRecordUnbuffered(r *kgo.Record, _ bool) {
	request := franzKafkaConsumerRequest{
		topic:   r.Topic,
		Record:  r,
		Context: r.Context,
	}
	franzKafkaConsumerInstrumenter.End(r.Context, request, franzKafkaConsumerResponse{}, nil)
}
