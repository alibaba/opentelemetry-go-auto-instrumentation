package main

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"os"
)

type MyHash struct {
	Key1 string `redis:"key1"`
	Key2 int    `redis:"key2"`
}

func main() {
	ctx := context.Background()
	rdb := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"shard1": "localhost:" + os.Getenv("REDIS_PORT"),
		},
		Password: "Hello1234",
	})
	_, err := rdb.HSet(ctx, "a", MyHash{
		Key1: "1",
		Key2: 2,
	}).Result()
	if err != nil {
		panic(err)
	}
	val := rdb.HVals(ctx, "a").Val()
	fmt.Printf("%v\n", val)
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "command", "", "redis", "", "localhost", "command: map[]", "command")
		verifier.VerifyDbAttributes(stubs[1][0], "hset", "", "redis", "", "localhost", "hset a key1 1 key2 2: 0", "hset")
		verifier.VerifyDbAttributes(stubs[2][0], "hvals", "", "redis", "", "localhost", "hvals a: []", "hvals")
	})
}
