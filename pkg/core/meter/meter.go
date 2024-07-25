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
