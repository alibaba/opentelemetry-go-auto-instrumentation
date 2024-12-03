// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package meter

import (
	"go.opentelemetry.io/otel/metric"
	"sync"
)

var globalMeter metric.Meter
var mu sync.Mutex

func SetMeter(meter metric.Meter) {
	mu.Lock()
	defer mu.Unlock()
	globalMeter = meter
}

func GetMeter() metric.Meter {
	mu.Lock()
	defer mu.Unlock()
	return globalMeter
}
