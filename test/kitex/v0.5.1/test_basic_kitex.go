// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"time"

	"github.com/cloudwego/kitex/client"
	"kitex/v0.5.1/kitex_gen/api"
	"kitex/v0.5.1/kitex_gen/api/hello"
)

func main() {
	go func() {
		svr := hello.NewServer(new(HelloImpl))

		err := svr.Run()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	time.Sleep(3 * time.Second)
	client, err := hello.NewClient("hello", client.WithHostPorts("0.0.0.0:8888"))
	if err != nil {
		log.Fatal(err)
	}
	req := &api.Request{Message: "my request"}
	resp, err := client.Echo(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
	time.Sleep(time.Second)
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyRpcClientAttributes(stubs[0][0], "api.Hello/api.Hello/echo", "kitex", "api.Hello", "api.Hello/echo")
		verifier.VerifyRpcServerAttributes(stubs[0][1], "api.Hello/api.Hello/echo", "kitex", "api.Hello", "api.Hello/echo")
	}, 1)
}
