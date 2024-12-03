// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"encoding/json"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"net/http"
)

func httpClientEnterHook(call api.CallContext, t *http.Transport, req *http.Request) {
	header, _ := json.Marshal(req.Header)
	fmt.Println("request header is ", string(header))
}

func httpClientExitHook(call api.CallContext, res *http.Response, err error) {
	header, _ := json.Marshal(res.Header)
	fmt.Println("response header is ", string(header))
}
