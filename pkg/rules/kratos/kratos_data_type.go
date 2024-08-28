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
//go:build ignore

package rule

import (
	"context"
	kt "github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"net/url"
	"strings"
)

var kratosServerInstrument = BuildKratosServerInstrumenter()

var kratosClientInstrument = BuildKratosClientInstrumenter()

var kratosPropagators = propagation.NewCompositeTextMapPropagator(OtelMetaData{}, propagation.Baggage{}, propagation.TraceContext{})

type kratosRequest struct {
	method        string
	addr          string
	header        transport.Header
	componentName string
	httpMethod    string
}

type kratosResponse struct {
	statusCode int
}

type TracerType int

const (
	TypeServer TracerType = iota
	TypeClient
)

const serviceHeader = "x-md-service-name"

// OtelMetaData is tracing metadata propagator
type OtelMetaData struct{}

var _ propagation.TextMapPropagator = OtelMetaData{}

// Inject sets metadata key-values from ctx into the carrier.
func (b OtelMetaData) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	app, ok := kt.FromContext(ctx)
	if ok {
		carrier.Set(serviceHeader, app.Name())
	}
}

// Extract returns a copy of parent with the metadata from the carrier added.
func (b OtelMetaData) Extract(parent context.Context, carrier propagation.TextMapCarrier) context.Context {
	name := carrier.Get(serviceHeader)
	if name == "" {
		return parent
	}
	if md, ok := metadata.FromServerContext(parent); ok {
		md.Set(serviceHeader, name)
		return parent
	}
	md := metadata.New()
	md.Set(serviceHeader, name)
	parent = metadata.NewServerContext(parent, md)
	return parent
}

// Fields returns the keys who's values are set with Inject.
func (b OtelMetaData) Fields() []string {
	return []string{serviceHeader}
}

// parseKratosFullMethod returns a span name following the OpenTelemetry semantic
// conventions as well as all applicable span attribute.KeyValue attributes based
// on a gRPC's FullMethod.
func parseKratosFullMethod(fullMethod string) (string, []attribute.KeyValue) {
	name := strings.TrimLeft(fullMethod, "/")
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 { //nolint:gomnd
		// Invalid format, does not follow `/package.service/method`.
		return name, []attribute.KeyValue{attribute.Key("rpc.operation").String(fullMethod)}
	}

	var attrs []attribute.KeyValue
	if service := parts[0]; service != "" {
		attrs = append(attrs, semconv.RPCService(service))
	}
	if method := parts[1]; method != "" {
		attrs = append(attrs, semconv.RPCMethod(method))
	}
	return name, attrs
}

func otelParseTarget(endpoint string) (address string, err error) {
	var u *url.URL
	u, err = url.Parse(endpoint)
	if err != nil {
		if u, err = url.Parse("http://" + endpoint); err != nil {
			return "", err
		}
		return u.Host, nil
	}
	if len(u.Path) > 1 {
		return u.Path[1:], nil
	}
	return endpoint, nil
}
