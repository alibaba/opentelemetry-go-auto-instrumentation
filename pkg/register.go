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

package pkg

import (
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/databasesql"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/echo"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/elasticsearch"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/fasthttp"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/fiberv2"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/gin"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/golog"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/gomicro"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/goredis"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/goredisv8"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/gorestful"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/gorm"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/goslog"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/grpc"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/hertz/client"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/hertz/server"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/http"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/iris"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/kitex"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/kratos/grpc"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/kratos/http"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/langchain"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/logrus"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/mongo"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/mux"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/nacos/config"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/nacos/dom"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/nacos/service"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/nacos/service_holder"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/nacos2_1_0/service_holder"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/nacos2_1_1/dom"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/otel-context"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/redigo"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/runtime"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/trpc"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/zap"
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/rules/zerolog"
)
