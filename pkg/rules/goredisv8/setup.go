package goredisv8

import (
	"context"
	"errors"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	redis "github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/trace"
)

var redisv8Instrumenter = BuildRedisv8Instrumenter()

var rv8Enabler = instrumenter.NewDefaultInstrumentEnabler()

var redisV8StartOptions = []trace.SpanStartOption{}

func afterNewRedisV8Client(call api.CallContext, client *redis.Client) {
	if !rv8Enabler.Enable() {
		return
	}
	client.AddHook(newOtRedisV8Hook(client.Options().Addr))
}

func afterNewFailOverRedisV8Client(call api.CallContext, client *redis.Client) {
	if !rv8Enabler.Enable() {
		return
	}
	client.AddHook(newOtRedisV8Hook(client.Options().Addr))
}

func afterNewConnRedisV8Client(call api.CallContext, conn *redis.Conn) {
	if !rv8Enabler.Enable() {
		return
	}
	conn.AddHook(newOtRedisV8Hook(conn.String()))
}

func afterNewClusterV8Client(call api.CallContext, client *redis.ClusterClient) {
	if !rv8Enabler.Enable() {
		return
	}
	client.AddHook(newOtRedisV8Hook(strings.Join(client.Options().Addrs, ",")))
}

func afterNewRingV8Client(call api.CallContext, client *redis.Ring) {
	if !rv8Enabler.Enable() {
		return
	}
	addrBuilder := strings.Builder{}
	for addr, _ := range client.Options().Addrs {
		addrBuilder.WriteString(addr)
	}
	client.AddHook(newOtRedisV8Hook(addrBuilder.String()))
}

type otRedisV8Hook struct {
	Addr string
}

func newOtRedisV8Hook(addr string) *otRedisV8Hook {
	return &otRedisV8Hook{
		Addr: addr,
	}
}

func (o *otRedisV8Hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	request := redisv8Data{
		cmd:  cmd,
		Host: o.Addr,
	}
	newCtx := redisv8Instrumenter.Start(ctx, request, redisV8StartOptions...)
	ctx = context.WithValue(ctx, redisV8Context, newCtx)
	return ctx, nil
}

func (o *otRedisV8Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	request := redisv8Data{
		cmd:  cmd,
		Host: o.Addr,
	}
	redisV8Ctx, ok := ctx.Value(redisV8Context).(context.Context)
	if !ok {
		redisV8Ctx = ctx
	}
	redisv8Instrumenter.End(redisV8Ctx, request, nil, cmd.Err())
	return nil
}

func (o *otRedisV8Hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	request := redisv8Data{
		cmd:  pipelineCmd,
		Host: o.Addr,
	}
	newCtx := redisv8Instrumenter.Start(ctx, request, redisV8StartOptions...)
	ctx = context.WithValue(ctx, redisV8Context, newCtx)
	return ctx, nil
}

func (o *otRedisV8Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	request := redisv8Data{
		cmd:  pipelineCmd,
		Host: o.Addr,
	}
	var tError error
	hasError := false
	errSb := strings.Builder{}
	for _, cmd := range cmds {
		if cmd.Err() != nil {
			errSb.WriteString(cmd.Err().Error())
			hasError = true
		}
	}
	tError = errors.New(errSb.String())
	redisV8Ctx, ok := ctx.Value(redisV8Context).(context.Context)
	if !ok {
		redisV8Ctx = ctx
	}
	if hasError {
		redisv8Instrumenter.End(redisV8Ctx, request, nil, tError)
	} else {
		redisv8Instrumenter.End(redisV8Ctx, request, nil, nil)
	}
	return nil
}
