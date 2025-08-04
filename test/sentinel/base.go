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

package main

import (
	"time"

	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
)

func SetUpResource() {
	api.InitDefault()
	flow.LoadRules([]*flow.Rule{
		{
			Resource:               "test",
			Threshold:              20,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})
	r := gin.New()
	r.Use(func(c *gin.Context) {
		e, err := api.Entry(
			"test",
			api.WithResourceType(base.ResTypeWeb),
			api.WithTrafficType(base.Inbound),
		)
		if err != nil {
			c.AbortWithStatus(429) // Too Many Requests
			return
		}
		defer e.Exit()
		c.Next()
	})
	r.GET("/test", func(c *gin.Context) {
		time.Sleep(time.Duration(time.Now().UnixNano()%300) * time.Millisecond)
		c.String(200, "Hello, world!")

	})
	r.Run(":8080")
}
