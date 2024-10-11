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

package rule

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {

	api.NewRule("github.com/go-kratos/kratos/v2/transport/http", "NewServer", "", "KratosNewHTTPServiceOnEnter", "").
		WithVersion("[2.6.3,2.8.1)").
		WithFileDeps("kratos_data_type.go").
		WithFileDeps("kratos_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/go-kratos/kratos/v2/transport/grpc", "NewServer", "", "KratosNewGRPCServiceOnEnter", "").
		WithVersion("[2.6.3,2.8.1)").
		WithFileDeps("kratos_data_type.go").
		WithFileDeps("kratos_otel_instrumenter.go").
		Register()

}
