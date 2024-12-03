// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package gin

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
)

var ginEnabler = instrumenter.NewDefaultInstrumentEnabler()
