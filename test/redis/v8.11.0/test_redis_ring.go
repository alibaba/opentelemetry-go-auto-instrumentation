// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"os"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"shard1": "localhost:" + os.Getenv("REDIS_PORT"),
		},
	})
	_, err := rdb.HSet(ctx, "a", map[string]string{
		"a": "b",
	}).Result()
	if err != nil {
		panic(err)
	}
	val := rdb.HVals(ctx, "a").Val()
	fmt.Printf("%v\n", val)
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "hset", "redis", "shard1", "hset a a b", "hset")
		verifier.VerifyDbAttributes(stubs[1][0], "hvals", "redis", "shard1", "hvals a", "hvals")
	}, 3)
}
