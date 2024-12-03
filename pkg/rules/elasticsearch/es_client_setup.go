// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package elasticsearch

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"net/http"
	"strings"
)

var esInstrumenter = BuildElasticSearchInstrumenter()

var esEnabler = instrumenter.NewDefaultInstrumentEnabler()

func beforeElasticSearchPerform(call api.CallContext, client *elasticsearch.BaseClient, request *http.Request) {
	if !esEnabler.Enable() {
		return
	}
	var addresses []string
	for _, u := range client.Transport.(*elastictransport.Client).URLs() {
		addresses = append(addresses, u.String())
	}
	op, params := getEsOpAndParams(request)
	er := &esRequest{
		request: request,
		address: strings.Join(addresses, ","),
		op:      op,
		params:  params,
	}
	newCtx := esInstrumenter.Start(request.Context(), er)
	call.SetKeyData("ctx", newCtx)
	call.SetKeyData("request", er)
}

func afterElasticSearchPerform(call api.CallContext, response *http.Response, err error) {
	if !esEnabler.Enable() {
		return
	}
	newCtx := call.GetKeyData("ctx").(context.Context)
	er := call.GetKeyData("request").(*esRequest)
	esInstrumenter.End(newCtx, er, response, err)
}

func getEsOpAndParams(req *http.Request) (string, []any) {
	if req == nil || req.URL == nil {
		return "UNKNOWN", nil
	}
	path := req.URL.Path
	paths := strings.Split(path, "/")
	if len(paths) <= 1 {
		return "UNKNOWN", nil
	}
	if len(paths) == 2 {
		return strings.ToLower(req.Method), nil
	}
	params := make([]any, len(paths)-2)
	// path[0] should be the index name
	for i := 2; i < len(paths); i++ {
		params[i-2] = paths[i]
	}
	return paths[2], params
}
