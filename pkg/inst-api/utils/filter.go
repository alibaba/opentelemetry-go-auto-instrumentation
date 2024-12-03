// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package utils

import "net/url"

type UrlFilter interface {
	FilterUrl(url *url.URL) bool
}

type SpanNameFilter interface {
	FilterSpanName(spanName string) bool
}

type DefaultUrlFilter struct {
}

func (d DefaultUrlFilter) FilterUrl(url *url.URL) bool {
	return false
}
