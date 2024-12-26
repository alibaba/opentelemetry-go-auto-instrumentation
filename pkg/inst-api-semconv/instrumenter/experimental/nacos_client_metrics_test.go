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

package experimental

import (
	"go.opentelemetry.io/otel/sdk/metric"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitNacosExperimentalMetrics_GlobalMeterNil_NoMetricsInitialized(t *testing.T) {
	InitNacosExperimentalMetrics(nil)
	assert.Nil(t, ClientServiceInfoMapSize)
	assert.Nil(t, ClientConfigCacheMapSize)
	assert.Nil(t, ClientDomBeatMapSize)
	assert.Nil(t, ClientConfigRequestDuration)
	assert.Nil(t, ClientNamingRequestDuration)
}

func TestInitNacosExperimentalMetrics_GlobalMeterNotNull_AllMetricsInitialized(t *testing.T) {
	mp := metric.NewMeterProvider()
	InitNacosExperimentalMetrics(mp.Meter("a"))
	assert.NotNil(t, ClientServiceInfoMapSize)
	assert.NotNil(t, ClientConfigCacheMapSize)
	assert.NotNil(t, ClientDomBeatMapSize)
	assert.NotNil(t, ClientConfigRequestDuration)
	assert.NotNil(t, ClientNamingRequestDuration)
}
