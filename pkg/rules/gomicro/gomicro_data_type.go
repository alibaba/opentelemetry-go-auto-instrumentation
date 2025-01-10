// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

package gomicro

import (
	"context"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/server"
	"go.opentelemetry.io/otel/propagation"
)

var goMicroClientInstrument = BuildGoMicroClientInstrumenter()

type requestType int

const (
	MessageRequest requestType = iota
	CallRequest
	StreamRequest
)

type goMicroRequest struct {
	reqType     requestType
	request     client.Request
	msg         client.Message
	ctx         context.Context
	propagators propagation.TextMapCarrier
}

type goMicroServerRequest struct {
	reqType     requestType
	request     server.Request
	msg         server.Message
	ctx         context.Context
	propagators propagation.TextMapCarrier
}

type goMicroResponse struct {
	response interface{}
	ctx      context.Context
	err      error
}
