//go:build ignore

package goredis

import (
	"context"
	redis "github.com/redis/go-redis/v9"
	"net"
	"strings"
)

var goRedisInstrumenter = BuildGoRedisOtelInstrumenter()

func afterNewRedisClient(call redis.CallContext, client *redis.Client) {
	client.AddHook(newOtRedisHook(client.Options().Addr))
}

func afterNewFailOverRedisClient(call redis.CallContext, client *redis.Client) {
	client.AddHook(newOtRedisHook(client.Options().Addr))
}

func afterNewClusterClient(call redis.CallContext, client *redis.ClusterClient) {
	client.OnNewNode(func(rdb *redis.Client) {
		rdb.AddHook(newOtRedisHook(rdb.Options().Addr))
	})
}

func afterNewRingClient(call redis.CallContext, client *redis.Ring) {
	client.OnNewNode(func(rdb *redis.Client) {
		rdb.AddHook(newOtRedisHook(rdb.Options().Addr))
	})
}

type otRedisHook struct {
	Addr string
}

func newOtRedisHook(addr string) *otRedisHook {
	return &otRedisHook{
		Addr: addr,
	}
}

func (o *otRedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := next(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		return conn, err
	}
}

func (o *otRedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if strings.Contains(cmd.FullName(), "ping") || strings.Contains(cmd.FullName(), "PING") {
			return next(ctx, cmd)
		}
		request := goRedisRequest{
			cmd:      cmd,
			endpoint: o.Addr,
		}
		ctx = goRedisInstrumenter.Start(ctx, request)
		if err := next(ctx, cmd); err != nil {
			goRedisInstrumenter.End(ctx, request, nil, err)
			return err
		}
		goRedisInstrumenter.End(ctx, request, nil, nil)
		return nil
	}
}

func (o *otRedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if err := next(ctx, cmds); err != nil {
			return err
		}
		return nil
	}
}
