// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//go:build ignore

package rule

import (
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func zapLogWriteOnEnter(call zapcore.CallContext, ce *zapcore.CheckedEntry, fields ...zap.Field) {
	var fieldsTemp []zap.Field
	traceId, spanId := trace.GetTraceAndSpanId()
	if traceId != "" {
		fieldsTemp = append(fieldsTemp, zap.String("trace_id", traceId))
	}
	if spanId != "" {
		fieldsTemp = append(fieldsTemp, zap.String("span_id", spanId))
	}
	if fields == nil {
		fields = fieldsTemp
	}
	call.SetParam(1, fields)
	return
}
