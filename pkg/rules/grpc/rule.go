// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rule

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

func init() {
	//client
	api.NewRule("google.golang.org/grpc", "DialContext", "", "grpcClientOnEnter", "grpcClientOnExit").
		WithFileDeps("config.go", "grpc_data_type.go", "interceptorinfo.go", "metadata_supplier.go", "grpc_otel_instrumenter.go").
		WithFileDeps("grpc_attrs_extractor.go", "grpc_attrs_getter.go", "grpc_span_name_extractor.go", "grpc_status_code_converter.go").
		WithVersion("[1.44.0,1.66.2)").
		Register()

	api.NewRule("google.golang.org/grpc", "NewClient", "", "grpcNewClientOnEnter", "grpcNewClientOnExit").
		WithFileDeps("config.go", "grpc_data_type.go", "interceptorinfo.go", "metadata_supplier.go", "grpc_otel_instrumenter.go").
		WithFileDeps("grpc_attrs_extractor.go", "grpc_attrs_getter.go", "grpc_span_name_extractor.go", "grpc_status_code_converter.go").
		WithVersion("[1.63.0,1.66.2)").
		Register()

	//server
	api.NewRule("google.golang.org/grpc", "NewServer", "", "grpcServerOnEnter", "grpcServerOnExit").
		WithVersion("[1.44.0,1.66.2)").
		Register()

}
