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

package amqp091go

import (
	"strings"
)

type RabbitRequest struct {
	exchange        string
	routingKey      string
	operationName   string
	destinationName string
	messageId       string
	bodySize        int64
	conversationID  string
}

func (r *RabbitRequest) Get(key string) string {
	if r == nil || r.conversationID == "" {
		return ""
	}
	vs := strings.Split(r.conversationID, ":")
	if len(vs) < 2 {
		return ""
	}
	if vs[0] == key {
		return vs[1]
	}
	return ""
}
func (r *RabbitRequest) Set(key string, value string) {
	if r == nil || r.conversationID != "" {
		return
	}
	vs := key + ":" + value
	r.conversationID = vs
}
func (r *RabbitRequest) Keys() []string {
	if r == nil || r.conversationID == "" {
		return []string{}
	}
	vs := strings.Split(r.conversationID, ":")
	if len(vs) < 2 {
		return []string{}
	}
	return []string{vs[0]}
}
