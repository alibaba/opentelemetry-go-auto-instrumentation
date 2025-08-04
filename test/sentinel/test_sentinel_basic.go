// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

package main

import (
	"github.com/alibaba/loongsuite-go-agent/test/verifier"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func main() {
	api.InitDefault()
	flow.LoadRules([]*flow.Rule{
		{
			Resource:               "test",
			Threshold:              20,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
		{
			Resource:               "test1",
			Threshold:              0,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})
	e, _ := api.Entry(
		"test",
		api.WithResourceType(base.ResTypeWeb),
		api.WithTrafficType(base.Inbound),
	)
	if e != nil {
		e.Exit()
	}
	e, _ = api.Entry(
		"test_block",
		api.WithResourceType(base.ResTypeWeb),
		api.WithTrafficType(base.Inbound),
	)
	if e != nil {
		e.Exit()
	}

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifySentinelAttributes(stubs[0][0], "test", "Inbound", false)
		verifier.VerifySentinelAttributes(stubs[1][0], "test_block", "Inbound", true)
	}, 2)
}
