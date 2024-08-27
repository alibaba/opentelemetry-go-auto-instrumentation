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

package goredis

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	api.NewRule("github.com/redis/go-redis/v9", "NewClient", "", "", "afterNewRedisClient").
		WithFileDeps("goredis_data_type.go", "goredis_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/redis/go-redis/v9", "NewFailoverClient", "", "", "afterNewFailOverRedisClient").
		WithVersion("[9.0.5,9.5.2)").
		WithFileDeps("goredis_data_type.go", "goredis_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/redis/go-redis/v9", "NewSentinelClient", "", "", "afterNewSentinelClient").
		WithVersion("[9.0.5,9.5.2)").
		WithFileDeps("goredis_data_type.go", "goredis_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/redis/go-redis/v9", "Conn", "*Client", "", "afterClientConn").
		WithVersion("[9.0.5,9.5.2)").
		WithFileDeps("goredis_data_type.go", "goredis_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/redis/go-redis/v9", "NewClusterClient", "", "", "afterNewClusterClient").
		WithVersion("[9.0.5,9.5.2)").
		WithFileDeps("goredis_data_type.go", "goredis_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/redis/go-redis/v9", "NewRing", "", "", "afterNewRingClient").
		WithVersion("[9.0.5,9.5.2)").
		WithFileDeps("goredis_data_type.go", "goredis_otel_instrumenter.go").
		Register()
}
