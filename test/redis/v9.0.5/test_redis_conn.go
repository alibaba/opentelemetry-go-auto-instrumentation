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
	cn := rdb.Conn()
	defer cn.Close()

	if err := cn.ClientSetName(ctx, "myclient").Err(); err != nil {
		panic(err)
	}

	name, err := cn.ClientGetName(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client name", name)

	_, err = cn.Set(ctx, "a", "b", 5*time.Second).Result()
	if err != nil {
		panic(err)
	}
	_, err = cn.Get(ctx, "a").Result()
	if err != nil {
		panic(err)
	}
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "client", "", "redis", "", "localhost", "client setname myclient: false", "client")
		verifier.VerifyDbAttributes(stubs[1][0], "client", "", "redis", "", "localhost", "client getname: ", "client")
		verifier.VerifyDbAttributes(stubs[2][0], "set", "", "redis", "", "localhost", "set a b ex 5: ", "set")
		verifier.VerifyDbAttributes(stubs[3][0], "get", "", "redis", "", "localhost", "get a: ", "get")
	})
}
