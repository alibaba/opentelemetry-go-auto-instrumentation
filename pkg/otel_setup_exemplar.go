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

//go:build exemplar
// +build exemplar

package pkg

import (
	"context"
	"os"
	"strconv"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
)

// Exemplar configuration environment variables
const (
	exemplar_enabled        = "OTEL_GO_AUTO_EXEMPLARS_ENABLED"
	exemplar_filter         = "OTEL_METRICS_EXEMPLAR_FILTER"
	exemplar_reservoir_size = "OTEL_METRICS_EXEMPLAR_RESERVOIR_SIZE"
)

// initMetricsWithExemplars initializes metrics with exemplar support
func initMetricsWithExemplars(ctx context.Context, isGRPC bool) (metric.Exporter, []metric.Option, error) {
	var opts []metric.Option
	var metricExporter metric.Exporter
	var err error

	// Create OTLP exporter
	if isGRPC {
		metricExporter, err = otlpmetricgrpc.New(ctx)
	} else {
		metricExporter, err = otlpmetrichttp.New(ctx)
	}

	if err != nil {
		return nil, nil, err
	}

	opts = append(opts, metric.WithReader(metric.NewPeriodicReader(metricExporter)))

	// Configure exemplar support if enabled
	if os.Getenv(exemplar_enabled) != "false" {
		// Set exemplar filter
		var exemplarFilter exemplar.Filter
		switch os.Getenv(exemplar_filter) {
		case "always_on":
			exemplarFilter = exemplar.AlwaysOnFilter
		case "always_off":
			exemplarFilter = exemplar.AlwaysOffFilter
		case "trace_based", "":
			exemplarFilter = exemplar.TraceBasedFilter
		default:
			exemplarFilter = exemplar.TraceBasedFilter
		}
		opts = append(opts, metric.WithExemplarFilter(exemplarFilter))

		// Configure exemplar reservoir size
		reservoirSize := 5 // default
		if size := os.Getenv(exemplar_reservoir_size); size != "" {
			if parsed, err := strconv.Atoi(size); err == nil && parsed > 0 {
				reservoirSize = parsed
			}
		}

		// Add view for histogram exemplars
		histogramView := metric.NewView(
			metric.Instrument{Kind: metric.InstrumentKindHistogram},
			metric.Stream{
				ExemplarReservoirProviderSelector: func(agg metric.Aggregation) exemplar.ReservoirProvider {
					return exemplar.FixedSizeReservoirProvider(reservoirSize)
				},
			},
		)
		opts = append(opts, metric.WithView(histogramView))
	}

	return metricExporter, opts, nil
}