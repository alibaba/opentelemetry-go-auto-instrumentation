package main

import (
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"os"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/gomodule/redigo/redis"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:"+os.Getenv("REDIS_PORT"))
	if err != nil {
		panic(err)
	}
	defer c.Close()
	c.Do("SET", "foo", "bar")
	_, err = c.Do("UNKNOWN", "nononononono")
	println(err.Error())
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		
	}, 2)
}
