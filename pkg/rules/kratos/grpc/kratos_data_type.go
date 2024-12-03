// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package grpc

type kratosRequest struct {
	protocolType    string
	serviceName     string
	serviceId       string
	serviceVersion  string
	serviceEndpoint []string
	serviceMeta     map[string]string
}
