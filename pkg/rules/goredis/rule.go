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

	api.NewRule("github.com/redis/go-redis/v9", "NewClusterClient", "", "", "afterNewClusterClient").
		WithVersion("[9.0.5,9.5.2)").
		WithFileDeps("goredis_data_type.go", "goredis_otel_instrumenter.go").
		Register()

	api.NewRule("github.com/redis/go-redis/v9", "NewRing", "", "", "afterNewRingClient").
		WithVersion("[9.0.5,9.5.2)").
		WithFileDeps("goredis_data_type.go", "goredis_otel_instrumenter.go").
		Register()
}
