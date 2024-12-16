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
	"os"
	"testing"
)

func TestGetSpanSuppressionStrategyFromEnv(t *testing.T) {
	tests := map[string]SpanSuppressorStrategy{
		"none":      &NoneStrategy{},
		"span-kind": &SpanKindStrategy{},
		"":          &SemConvStrategy{},
		"unknown":   &SemConvStrategy{},
	}

	for value, expectedStrategy := range tests {
		os.Setenv("OTEL_INSTRUMENTATION_EXPERIMENTAL_SPAN_SUPPRESSION_STRATEGY", value)
		defer os.Unsetenv("OTEL_INSTRUMENTATION_EXPERIMENTAL_SPAN_SUPPRESSION_STRATEGY")

		actualStrategy := getSpanSuppressionStrategyFromEnv()

		if expectedStrategy != actualStrategy {
			panic("Expected strategy does not match actual strategy")
		}
	}
}
