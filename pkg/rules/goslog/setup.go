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

//go:build ignore

package golog

import (
	"context"
	"go.opentelemetry.io/otel/sdk/trace"
	"log/slog"
)

func goSlogWriteOnEnter(call slog.CallContext, ce *slog.Logger, ctx context.Context, level slog.Level, msg string, args ...any) {
	traceId, spanId := trace.GetTraceAndSpanId()
	if traceId != "" {
		msg = msg + " trace_id=" + traceId
	}
	if spanId != "" {
		msg = msg + " span_id=" + spanId
	}
	call.SetParam(3, msg)
	return
}
