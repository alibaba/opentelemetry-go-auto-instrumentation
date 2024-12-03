// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package gin

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/gin-gonic/gin"
)

func htmlOnEnter(call api.CallContext, c *gin.Context, code int, name string, obj any) {
	if !ginEnabler.Enable() {
		return
	}
	if c == nil {
		return
	}
	lcs := trace.LocalRootSpanFromGLS()
	if lcs != nil && c.FullPath() != "" && c.Request != nil && c.Request.URL != nil && (c.FullPath() != c.Request.URL.Path) {
		lcs.SetName(c.FullPath())
	}
}
