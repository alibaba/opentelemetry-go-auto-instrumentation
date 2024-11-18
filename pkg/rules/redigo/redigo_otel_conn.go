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

package redigo

import (
	"container/list"
	"context"
	"github.com/gomodule/redigo/redis"
	"os"
	"strconv"
	"time"
)

const max_queue_length = 2048

var configuredQueueLength int

var commandQueue = list.New()

var redigoInstrumenter = BuildRedigoInstrumenter()

type armsConn struct {
	redis.Conn
	endpoint string
	ctx      context.Context
}

func (a *armsConn) Close() error {
	return a.Conn.Close()
}

func (a *armsConn) Err() error {
	return a.Conn.Err()
}

func (a *armsConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	req := &redigoRequest{
		args:     args,
		endpoint: a.endpoint,
		cmd:      commandName,
	}
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	startTime := time.Now()
	reply, err = a.Conn.Do(commandName, args...)
	endTime := time.Now()
	redigoInstrumenter.StartAndEnd(ctx, req, nil, err, startTime, endTime)
	return
}

func (a *armsConn) Send(commandName string, args ...interface{}) error {
	now := time.Now()
	req := &redigoRequest{
		args:      args,
		endpoint:  a.endpoint,
		cmd:       commandName,
		startTime: now,
	}
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	req.ctx = ctx
	push(req)
	return a.Conn.Send(commandName, args...)
}

func (a *armsConn) Flush() error {
	return a.Conn.Flush()
}

func (a *armsConn) Receive() (reply interface{}, err error) {
	reply, err = a.Conn.Receive()
	req := pop()
	if req != nil {
		now := time.Now()
		redigoInstrumenter.StartAndEnd(req.ctx, req, nil, err, req.startTime, now)
	}
	return
}

func push(request *redigoRequest) {
	if commandQueue != nil && commandQueue.Len() > getMaxQueueLength() {
		return
	}
	commandQueue.PushBack(request)
}

func pop() *redigoRequest {
	front := commandQueue.Front()
	commandQueue.Remove(front)
	p, ok := front.Value.(*redigoRequest)
	if ok {
		return p
	}
	return nil
}

func getMaxQueueLength() int {
	if configuredQueueLength == 0 {
		var e = os.Getenv("MAX_REDIGO_QUEUE_LENGTH")
		if e != "" {
			configuredQueueLength, _ = strconv.Atoi(os.Getenv(e))
		} else {
			configuredQueueLength = max_queue_length
		}
	}
	return configuredQueueLength
}
