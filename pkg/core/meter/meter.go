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

package meter

import "go.opentelemetry.io/otel/metric"

var globalMeterProvider MeterProvider

type MeterMeta map[string]string

type MeterWithMeta struct {
	meter metric.Meter
	meta  MeterMeta
}

func NewMeterWithMeta(meter metric.Meter) *MeterWithMeta {
	return &MeterWithMeta{meter: meter, meta: make(map[string]string)}
}

func NewMeterWithMetaWithKVs(meter metric.Meter, tag MeterMeta) *MeterWithMeta {
	return &MeterWithMeta{meter: meter, meta: tag}
}

func (m *MeterWithMeta) Meter() metric.Meter {
	return m.meter
}

func (m *MeterWithMeta) Metas() MeterMeta {
	return m.meta
}

func (m *MeterWithMeta) Meta(key string) (string, bool) {
	val, ok := m.meta[key]
	return val, ok
}

func (m *MeterWithMeta) SetMeta(key string, val string) {
	m.meta[key] = val
}

type MeterProvider interface {
	GetMeters() []MeterWithMeta
}

type OtelMeterProvider struct {
	metricsMeter MeterWithMeta
}

func NewOtelMeterProvider(metricsMeter MeterWithMeta) *OtelMeterProvider {
	return &OtelMeterProvider{metricsMeter: metricsMeter}
}

func (o *OtelMeterProvider) GetMeters() []MeterWithMeta {
	return []MeterWithMeta{o.metricsMeter}
}

func GetMeterProvider() MeterProvider {
	return globalMeterProvider
}

func SetMeterProvider(meterProvider MeterProvider) {
	globalMeterProvider = meterProvider
}
