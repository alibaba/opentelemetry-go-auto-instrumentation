//go:build ignore

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package baggage // import "go.opentelemetry.io/otel/baggage"

import (
	"context"

	"go.opentelemetry.io/otel/internal/baggage"
)

// ContextWithBaggage returns a copy of parent with baggage.
func ContextWithBaggage(parent context.Context, b Baggage) context.Context {
	SetBaggageToGLS(&b)
	// Delegate so any hooks for the OpenTracing bridge are handled.
	return baggage.ContextWithList(parent, b.list)
}

// ContextWithoutBaggage returns a copy of parent with no baggage.
func ContextWithoutBaggage(parent context.Context) context.Context {
	// Delegate so any hooks for the OpenTracing bridge are handled.
	return baggage.ContextWithList(parent, nil)
}

// FromContext returns the baggage contained in ctx.
func FromContext(ctx context.Context) Baggage {
	if ctx == nil || baggage.ListFromContext(ctx) == nil {
		if b := GetBaggageFromGLS(); b != nil {
			return *b
		}
	}

	if ctx == nil {
		return Baggage{}
	} else {
		// Delegate so any hooks for the OpenTracing bridge are handled.
		return Baggage{list: baggage.ListFromContext(ctx)}
	}
}
