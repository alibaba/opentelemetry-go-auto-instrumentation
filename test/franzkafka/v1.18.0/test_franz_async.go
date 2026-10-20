// Copyright (c) 2026 Alibaba Group Holding Ltd.
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

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alibaba/loongsuite-go/pkg/testaccess"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

const recordCount = 350 // Exceeds the GLS limit of 300 spans.

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	topic := fmt.Sprintf("franz-async-%d", time.Now().UnixNano())
	client, err := kgo.NewClient(kgo.SeedBrokers("localhost:"+os.Getenv("FRANZ_KAFKA_PORT")), kgo.MaxBufferedBytes(1024))
	check(err == nil, "create client: %v", err)
	defer client.Close()
	responses, err := kadm.NewClient(client).CreateTopics(ctx, 1, 1, nil, topic)
	check(err == nil, "create topic: %v", err)
	for _, response := range responses {
		check(response.Err == nil, "create topic: %v", response.Err)
	}

	testaccess.ResetTestSpans()
	producers := make(map[trace.SpanID]trace.SpanContext, recordCount)
	for i := 0; i < recordCount; i++ {
		r := &kgo.Record{Topic: topic, Value: []byte("async record")}
		send(ctx, client, context.Background(), r, trace.SpanContext{}, false)
		sc := trace.SpanContextFromContext(r.Context)
		producers[sc.SpanID()] = sc
	}
	spans := exportedSpans(recordCount)
	traces := make(map[trace.TraceID]bool, recordCount)
	for _, span := range spans {
		check(!span.Parent.IsValid(), "background produce inherited parent %s", span.Parent.SpanID())
		traces[span.SpanContext.TraceID()] = true
	}
	check(len(traces) == recordCount, "expected %d independent traces, got %d", recordCount, len(traces))

	// Buffered fetches start on fetch goroutines; polling ends spans elsewhere.
	// Every consumer must retain the producer context propagated in its headers.
	testaccess.ResetTestSpans()
	consumer, err := kgo.NewClient(kgo.SeedBrokers("localhost:"+os.Getenv("FRANZ_KAFKA_PORT")), kgo.ConsumePartitions(map[string]map[int32]kgo.Offset{topic: {0: kgo.NewOffset().AtStart()}}))
	check(err == nil, "create consumer: %v", err)
	received := 0
	for received < recordCount {
		fetches := consumer.PollFetches(ctx)
		check(ctx.Err() == nil, "poll timeout: %v", ctx.Err())
		check(len(fetches.Errors()) == 0, "fetch errors: %v", fetches.Errors())
		received += fetches.NumRecords()
	}
	consumer.Close()
	spans = exportedSpans(recordCount)
	seen := make(map[trace.SpanID]bool, recordCount)
	for _, span := range spans {
		parent, ok := producers[span.Parent.SpanID()]
		check(ok && parent.TraceID() == span.SpanContext.TraceID(), "consumer lost producer parent: %v", span.Parent)
		seen[span.Parent.SpanID()] = true
	}
	check(len(seen) == recordCount, "expected %d consumed producer contexts, got %d", recordCount, len(seen))
	assertCaller(trace.SpanContext{})

	for _, explicit := range []bool{true, false} {
		testaccess.ResetTestSpans()
		parentCtx, parent := otel.Tracer("franz-async-test").Start(context.Background(), "business-parent")
		for i := 0; i < 4; i++ {
			r := &kgo.Record{Topic: topic, Value: []byte("child")}
			if explicit {
				r.Context = parentCtx
			}
			// With no Record.Context, the valid caller GLS parent must still work.
			send(ctx, client, context.Background(), r, parent.SpanContext(), false)
		}
		parent.End()
		spans = exportedSpans(5)
		children := 0
		for _, span := range spans {
			if span.Name == "business-parent" {
				continue
			}
			check(span.Parent.SpanID() == parent.SpanContext().SpanID() && span.SpanContext.TraceID() == parent.SpanContext().TraceID(), "lost business parent (explicit=%v): %v", explicit, span.Parent)
			children++
		}
		check(children == 4, "expected four children, got %d", children)
		assertCaller(trace.SpanContext{})
	}

	testaccess.ResetTestSpans()
	remote := trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}, TraceFlags: trace.FlagsSampled, Remote: true})
	r := &kgo.Record{Topic: topic, Value: []byte("remote child"), Context: trace.ContextWithRemoteSpanContext(context.Background(), remote)}
	send(ctx, client, context.Background(), r, trace.SpanContext{}, false)
	spans = exportedSpans(1)
	check(spans[0].Parent.Equal(remote), "lost remote parent: %v", spans[0].Parent)

	testaccess.ResetTestSpans()
	// Oversized records fail through the asynchronous completion hooks too.
	send(ctx, client, context.Background(), &kgo.Record{Topic: topic, Value: make([]byte, 2048)}, trace.SpanContext{}, true)
	spans = exportedSpans(1)
	check(spans[0].Status.Code == codes.Error, "failed produce did not end with error status")
	fmt.Println("franz async regression passed: 350 independent producers, 350 propagated consumers, explicit/GLS/remote parents, and failed produce")
}

func send(waitCtx context.Context, client *kgo.Client, parent context.Context, record *kgo.Record, caller trace.SpanContext, wantError bool) {
	done := make(chan error, 1)
	client.Produce(parent, record, func(_ *kgo.Record, err error) { done <- err })
	assertCaller(caller)
	select {
	case err := <-done:
		check((err != nil) == wantError, "produce error = %v, want error = %v", err, wantError)
	case <-waitCtx.Done():
		panic(waitCtx.Err())
	}
	check(client.Flush(waitCtx) == nil, "flush failed")
	assertCaller(caller)
	sc := trace.SpanContextFromContext(record.Context)
	check(sc.IsValid(), "record lost its span context")
	carrier := propagation.MapCarrier{}
	for _, header := range record.Headers {
		carrier[header.Key] = string(header.Value)
	}
	propagated := trace.SpanContextFromContext(propagation.TraceContext{}.Extract(context.Background(), carrier))
	check(propagated.TraceID() == sc.TraceID() && propagated.SpanID() == sc.SpanID(), "record headers do not carry producer span")
}

func assertCaller(want trace.SpanContext) {
	span := trace.SpanFromContext(context.Background())
	check(span != nil, "SpanFromContext must preserve a non-nil noop span when GLS is empty")
	var got trace.SpanContext
	if span != nil {
		got = span.SpanContext()
	}
	check(got.TraceID() == want.TraceID() && got.SpanID() == want.SpanID(), "caller GLS changed: got %v, want %v", got, want)
}

func exportedSpans(want int) tracetest.SpanStubs {
	spans := *testaccess.GetTestSpans().(*tracetest.SpanStubs)
	check(len(spans) == want, "expected %d completed spans, got %d", want, len(spans))
	return spans
}

func check(ok bool, format string, args ...interface{}) {
	if !ok {
		panic(fmt.Sprintf(format, args...))
	}
}
