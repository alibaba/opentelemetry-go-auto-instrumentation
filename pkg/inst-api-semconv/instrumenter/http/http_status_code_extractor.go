// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const invalidHttpStatusCode = "INVALID_HTTP_STATUS_CODE"

type HttpClientSpanStatusExtractor[REQUEST any, RESPONSE any] struct {
	Getter HttpCommonAttrsGetter[REQUEST, RESPONSE]
}

func (h HttpClientSpanStatusExtractor[REQUEST, RESPONSE]) Extract(span trace.Span, request REQUEST, response RESPONSE, err error) {
	statusCode := h.Getter.GetHttpResponseStatusCode(request, response, err)
	if statusCode >= 400 || statusCode < 100 {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Error, invalidHttpStatusCode)
		}
	}
}

type HttpServerSpanStatusExtractor[REQUEST any, RESPONSE any] struct {
	Getter HttpCommonAttrsGetter[REQUEST, RESPONSE]
}

func (h HttpServerSpanStatusExtractor[REQUEST, RESPONSE]) Extract(span trace.Span, request REQUEST, response RESPONSE, err error) {
	statusCode := h.Getter.GetHttpResponseStatusCode(request, response, err)
	if statusCode >= 500 || statusCode < 100 {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Error, invalidHttpStatusCode)
		}
	}
}
