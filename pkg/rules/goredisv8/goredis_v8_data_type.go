// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goredisv8

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
	"unicode/utf8"
	"unsafe"
)

const redisV8Context = "redis-v8-context"

var pipelineCmd = redis.NewCmd(context.Background(), "pipeline")

type redisv8Data struct {
	cmd  redis.Cmder
	Host string
}

func redisV8String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func redisV8AppendUTF8String(dst []byte, src []byte) []byte {
	dst = append(dst, src...)
	return dst
}

func redisV8Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
func redisV8AppendArg(b []byte, v interface{}) []byte {
	switch v := v.(type) {
	case nil:
		return append(b, "<nil>"...)
	case string:
		bts := redisV8Bytes(v)
		if utf8.Valid(bts) {
			return redisV8AppendUTF8String(b, bts)
		} else {
			return redisV8AppendUTF8String(b, redisV8Bytes("<string>"))
		}
	case []byte:
		if utf8.Valid(v) {
			return redisV8AppendUTF8String(b, v)
		} else {
			return redisV8AppendUTF8String(b, redisV8Bytes("<byte>"))
		}
	case int:
		return strconv.AppendInt(b, int64(v), 10)
	case int8:
		return strconv.AppendInt(b, int64(v), 10)
	case int16:
		return strconv.AppendInt(b, int64(v), 10)
	case int32:
		return strconv.AppendInt(b, int64(v), 10)
	case int64:
		return strconv.AppendInt(b, v, 10)
	case uint:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint8:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint16:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(b, v, 'f', -1, 64)
	case bool:
		if v {
			return append(b, "true"...)
		}
		return append(b, "false"...)
	case time.Time:
		return v.AppendFormat(b, time.RFC3339Nano)
	default:
		return append(b, fmt.Sprint(v)...)
	}
}
