// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"net/url"
	"testing"
)

func TestDefaultUrlFilter(t *testing.T) {
	filter := DefaultUrlFilter{}
	testCases := []struct {
		input    *url.URL
		expected bool
	}{
		{
			input:    &url.URL{Scheme: "http", Host: "example.com"},
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.input.String(), func(t *testing.T) {
			result := filter.FilterUrl(tc.input)
			if result != tc.expected {
				t.Errorf("FilterUrl(%v) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}
