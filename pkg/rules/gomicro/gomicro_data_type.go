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
