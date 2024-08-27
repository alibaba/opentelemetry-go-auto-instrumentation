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

package instrumenter

import (
	"errors"
	"testing"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type testSpan struct {
	trace.Span
	status *codes.Code
}

func (ts testSpan) SetStatus(status codes.Code, desc string) {
	*ts.status = status
}

func TestDefaultSpanStatusExtractor(t *testing.T) {
	unset := codes.Unset
	ts := testSpan{Span: noop.Span{}, status: &unset}
	d := defaultSpanStatusExtractor[interface{}, interface{}]{}
	d.Extract(ts, nil, nil, errors.New(""))
	if *ts.status != codes.Error {
		t.Fatal("expected error code")
	}
}

func TestAlwaysInternalExtractor(t *testing.T) {
	a := &AlwaysInternalExtractor[any]{}
	kind := a.Extract(nil)
	if kind != trace.SpanKindInternal {
		t.Fatal("expected internal kind")
	}
}

func TestAlwaysServerExtractor(t *testing.T) {
	a := &AlwaysServerExtractor[any]{}
	kind := a.Extract(nil)
	if kind != trace.SpanKindServer {
		t.Fatal("expected server kind")
	}
}

func TestAlwaysClientExtractor(t *testing.T) {
	a := &AlwaysClientExtractor[any]{}
	kind := a.Extract(nil)
	if kind != trace.SpanKindClient {
		t.Fatal("expected client kind")
	}
}

func TestAlwaysConsumerExtractor(t *testing.T) {
	a := &AlwaysConsumerExtractor[any]{}
	kind := a.Extract(nil)
	if kind != trace.SpanKindConsumer {
		t.Fatal("expected consumer kind")
	}
}

func TestAlwaysProducerExtractor(t *testing.T) {
	a := &AlwaysProducerExtractor[any]{}
	kind := a.Extract(nil)
	if kind != trace.SpanKindProducer {
		t.Fatal("expected producer kind")
	}
}
