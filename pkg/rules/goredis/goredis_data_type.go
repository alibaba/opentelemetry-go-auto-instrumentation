// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package goredis

import (
	redis "github.com/redis/go-redis/v9"
)

type goRedisRequest struct {
	cmd      redis.Cmder
	endpoint string
}
