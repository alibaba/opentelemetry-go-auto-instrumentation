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
	"go.opentelemetry.io/otel/metric"
	"log"
	"os"
)

var (
	ClientServiceInfoMapSize    metric.Int64ObservableGauge
	ClientConfigCacheMapSize    metric.Int64ObservableGauge
	ClientDomBeatMapSize        metric.Int64ObservableGauge
	ClientConfigRequestDuration metric.Float64Histogram
	ClientNamingRequestDuration metric.Float64Histogram
	GlobalMeter                 metric.Meter
)

type nacosEnabler struct{}

func (n nacosEnabler) Enable() bool {
	return os.Getenv("OTEL_INSTRUMENTATION_NACOS_EXPERIMENTAL_METRICS_ENABLE") == "true"
}

var NacosEnabler nacosEnabler

func InitNacosExperimentalMetrics(m metric.Meter) {
	GlobalMeter = m
	if GlobalMeter == nil {
		return
	}
	var err error
	ClientServiceInfoMapSize, err = GlobalMeter.Int64ObservableGauge("nacos.client.serviceinfo.size", metric.WithDescription("Size of service info map"))
	if err != nil {
		log.Printf("failed to init ClientServiceInfoMapSize metrics")
	}
	ClientConfigCacheMapSize, err = GlobalMeter.Int64ObservableGauge("nacos.client.configinfo.size", metric.WithDescription("Size of config cache map"))
	if err != nil {
		log.Printf("failed to init ClientConfigCacheMapSize metrics")
	}
	ClientDomBeatMapSize, err = GlobalMeter.Int64ObservableGauge("nacos.client.dombeat.size", metric.WithDescription("Size of dom beat map"))
	if err != nil {
		log.Printf("failed to init ClientDomBeatMapSize metrics")
	}
	ClientConfigRequestDuration, err = GlobalMeter.Float64Histogram("nacos.client.config.request.duration", metric.WithDescription("Duration of config request"))
	if err != nil {
		log.Printf("failed to init ClientConfigRequestDuration metrics")
	}
	ClientNamingRequestDuration, err = GlobalMeter.Float64Histogram("nacos.client.naming.request.duration", metric.WithDescription("Duration of naming request"))
	if err != nil {
		log.Printf("failed to init ClientNamingRequestDuration metrics")
	}
}
