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

package instrumenter

import (
	"context"
	"log"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

type testListener struct {
	startTime       time.Time
	endTime         time.Time
	startAttributes []attribute.KeyValue
	endAttributes   []attribute.KeyValue
}

func (t *testListener) OnBeforeStart(parentContext context.Context, startTimestamp time.Time) context.Context {
	t.startTime = startTimestamp
	return context.WithValue(parentContext, "test1", "a")
}

func (t *testListener) OnBeforeEnd(ctx context.Context, startAttributes []attribute.KeyValue, startTimestamp time.Time) context.Context {
	t.startAttributes = startAttributes
	return context.WithValue(ctx, "test2", "a")
}

func (t *testListener) OnAfterStart(context context.Context, endTimestamp time.Time) {
	t.endTime = endTimestamp
}

func (t *testListener) OnAfterEnd(context context.Context, endAttributes []attribute.KeyValue, endTimestamp time.Time) {
	t.endAttributes = endAttributes
}

func TestShadower(t *testing.T) {
	originAttrs := []attribute.KeyValue{
		attribute.String("a", "b"),
		attribute.String("a1", "a1"),
		attribute.String("a2", "a2"),
		attribute.String("a3", "a3"),
	}

	n := NoopAttrsShadower{}
	num, newAttrs := n.Shadow(originAttrs)
	if num != len(originAttrs) {
		log.Fatal("origin attrs length is not equal to new attrs length")
	}
	for i := 0; i < num; i++ {
		if newAttrs[i].Value != originAttrs[i].Value {
			log.Fatal("origin attrs value is not equal to new attrs value")
		}
	}
}

func TestOnBeforeStart(t *testing.T) {
	w := &testListener{}
	newCtx := w.OnBeforeStart(context.Background(), time.UnixMilli(123412341234))
	if w.startTime.UnixMilli() != 123412341234 {
		log.Fatal("start time is not equal to new start time")
	}
	if newCtx.Value("test1") != "a" {
		log.Fatal("key test1 is not equal to new key value")
	}
}

func TestOnBeforeEnd(t *testing.T) {
	w := &testListener{}
	w.OnBeforeEnd(context.Background(), []attribute.KeyValue{{
		Key:   "123",
		Value: attribute.StringValue("abcde"),
	}}, time.UnixMilli(123412341234))
	if w.startAttributes[0].Key != "123" {
		log.Fatal("start attribute key is not equal to new start attribute key")
	}
	if w.startAttributes[0].Value.AsString() != "abcde" {
		log.Fatal("start attribute value is not equal to new start attribute value")
	}
}

func TestOnAfterStart(t *testing.T) {
	w := &testListener{}
	w.OnAfterStart(context.Background(), time.UnixMilli(123412341234))
	if w.endTime.UnixMilli() != 123412341234 {
		log.Fatal("start time is not equal to new start time")
	}
}

func TestOnAfterEnd(t *testing.T) {
	w := &testListener{}
	w.OnAfterEnd(context.Background(), []attribute.KeyValue{{
		Key:   "123",
		Value: attribute.StringValue("abcde"),
	}}, time.UnixMilli(123412341234))
	if w.endAttributes[0].Key != "123" {
		log.Fatal("start attribute key is not equal to new start attribute key")
	}
	if w.endAttributes[0].Value.AsString() != "abcde" {
		log.Fatal("start attribute value is not equal to new start attribute value")
	}
}
