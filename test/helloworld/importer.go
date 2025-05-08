// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

package main

import (
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg" // use otel setup
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/golog"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/goslog"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/grpc"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/http"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/otel-context"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/test/fmt1"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/test/fmt4"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/test/fmt5"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/test/fmt6"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/test/fmt7"
)
