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
