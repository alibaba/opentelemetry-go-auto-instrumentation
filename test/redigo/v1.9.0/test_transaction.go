package main

import (
	"fmt"
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
	c.Send("MULTI")
	c.Send("INCR", "foo")
	c.Send("INCR", "bar")
	r, err := c.Do("EXEC")
	fmt.Println(r) // prints [1, 1]

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {

	}, 1)
}
