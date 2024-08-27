// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package verifier

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/attribute"
)

const IS_IN_TEST = "IN_OTEL_TEST"

func GetAttribute(attrs []attribute.KeyValue, name string) attribute.Value {
	for _, attr := range attrs {
		if string(attr.Key) == name {
			return attr.Value
		}
	}
	return attribute.Value{}
}

func Assert(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Sprintf(format, args...))
	}
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		panic("Failed to create a free port: " + err.Error())
	}
	cli, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic("Failed to create a free port: " + err.Error())
	}
	defer cli.Close()
	return cli.Addr().(*net.TCPAddr).Port, nil
}

func GetServer(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp.Status, nil
}

func IsInTest() bool {
	return os.Getenv(IS_IN_TEST) == "true"
}
