package instrumenter

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"log"
	"testing"
	"time"
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
	w := OperationListenerWrapper{listener: &testListener{}}
	newCtx := w.OnBeforeStart(context.Background(), time.UnixMilli(123412341234))
	wListener := w.listener.(*testListener)
	if wListener.startTime.UnixMilli() != 123412341234 {
		log.Fatal("start time is not equal to new start time")
	}
	if newCtx.Value("test1") != "a" {
		log.Fatal("key test1 is not equal to new key value")
	}
}

func TestOnBeforeEnd(t *testing.T) {
	w := OperationListenerWrapper{listener: &testListener{}}
	w.OnBeforeEnd(context.Background(), []attribute.KeyValue{{
		Key:   "123",
		Value: attribute.StringValue("abcde"),
	}}, time.UnixMilli(123412341234))
	wListener := w.listener.(*testListener)
	if wListener.startAttributes[0].Key != "123" {
		log.Fatal("start attribute key is not equal to new start attribute key")
	}
	if wListener.startAttributes[0].Value.AsString() != "abcde" {
		log.Fatal("start attribute value is not equal to new start attribute value")
	}
}

func TestOnAfterStart(t *testing.T) {
	w := OperationListenerWrapper{listener: &testListener{}}
	w.OnAfterStart(context.Background(), time.UnixMilli(123412341234))
	wListener := w.listener.(*testListener)
	if wListener.startTime.UnixMilli() != 123412341234 {
		log.Fatal("start time is not equal to new start time")
	}
}

func TestOnAfterEnd(t *testing.T) {
	w := OperationListenerWrapper{listener: &testListener{}}
	w.OnAfterEnd(context.Background(), []attribute.KeyValue{{
		Key:   "123",
		Value: attribute.StringValue("abcde"),
	}}, time.UnixMilli(123412341234))
	wListener := w.listener.(*testListener)
	if wListener.endAttributes[0].Key != "123" {
		log.Fatal("start attribute key is not equal to new start attribute key")
	}
	if wListener.endAttributes[0].Value.AsString() != "abcde" {
		log.Fatal("start attribute value is not equal to new start attribute value")
	}
}
