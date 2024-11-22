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

func beforeElasticSearchPerform(call api.CallContext, client *elasticsearch.Client, request *http.Request) {
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
	call.SetParam(0, newCtx)
	call.SetParam(1, er)
}

func afterElasticSearchPerform(call api.CallContext, response *http.Response, err error) {
	if !esEnabler.Enable() {
		return
	}
	newCtx := call.GetParam(0).(context.Context)
	er := call.GetParam(1).(*esRequest)
	esInstrumenter.End(newCtx, er, response, err)
}

func getEsOpAndParams(req *http.Request) (string, []any) {
	if req == nil || req.URL == nil {
		return "UNKNOWN", nil
	}
	path := req.URL.Path
	paths := strings.Split(path, "/")
	if len(paths) == 0 {
		return "UNKNOWN", nil
	}
	if len(paths) == 1 {
		return "create", nil
	}
	params := make([]any, len(paths)-2)
	// path[0] should be the index name
	for i := 2; i < len(paths); i++ {
		params[i-2] = paths[i]
	}
	return paths[1], params
}
