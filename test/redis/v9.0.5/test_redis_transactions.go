package main

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + os.Getenv("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pipe := rdb.Pipeline()

	incr := pipe.Incr(ctx, "pipeline_counter")
	pipe.Expire(ctx, "pipeline_counter", time.Hour)

	_, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}

	// The value is available only after Exec is called.
	fmt.Println(incr.Val())
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "pipeline", "", "redis", "", "localhost", "pipeline incr/expire/", "pipeline")
	})
}
