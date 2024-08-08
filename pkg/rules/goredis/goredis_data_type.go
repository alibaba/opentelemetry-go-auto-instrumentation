//go:build ignore

package goredis

import (
	redis "github.com/redis/go-redis/v9"
)

type goRedisRequest struct {
	cmd      redis.Cmder
	endpoint string
}
