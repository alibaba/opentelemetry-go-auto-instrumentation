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

package main

import (
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/databasesql"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/echo"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/gin"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/goredis"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/gorm"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/grpc"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/hertz"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/http"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/logrus"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/mongo"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/mux"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/otsdk"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/runtime"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/test"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/zap"
)
