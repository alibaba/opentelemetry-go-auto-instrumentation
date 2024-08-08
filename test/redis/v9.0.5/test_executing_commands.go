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
	_, err := rdb.Set(ctx, "a", "b", 5*time.Second).Result()
	if err != nil {
		panic(err)
	}
	val, err := rdb.Get(ctx, "a").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "set", "", "redis", "", "localhost", "set a b ex 5: ", "set")
		verifier.VerifyDbAttributes(stubs[1][0], "get", "", "redis", "", "localhost", "get a: ", "get")
	})
}
