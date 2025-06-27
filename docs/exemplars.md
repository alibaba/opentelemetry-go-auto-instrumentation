# Exemplar Support in OpenTelemetry Go Auto-Instrumentation

This document describes how to enable and use exemplar support in the Alibaba OpenTelemetry Go auto-instrumentation project.

## What are Exemplars?

Exemplars are recorded values that associate OpenTelemetry context to metric events. They create a powerful bridge between aggregated metrics and individual traces, allowing you to:

- Jump from a metric spike directly to example traces that contributed to that spike
- Debug performance issues by correlating high-level metrics with detailed execution traces
- Understand the context behind metric anomalies

## Configuration

Exemplar support can be configured using environment variables:

### Enable/Disable Exemplars

```bash
# Enable exemplars (default: true)
export OTEL_GO_AUTO_EXEMPLARS_ENABLED=true

# Disable exemplars
export OTEL_GO_AUTO_EXEMPLARS_ENABLED=false
```

### Exemplar Filter

Controls when exemplars are recorded:

```bash
# Only record exemplars when there's an active trace (default)
export OTEL_METRICS_EXEMPLAR_FILTER=trace_based

# Always record exemplars
export OTEL_METRICS_EXEMPLAR_FILTER=always_on

# Never record exemplars
export OTEL_METRICS_EXEMPLAR_FILTER=always_off
```

### Reservoir Size

Controls how many exemplars are kept per metric:

```bash
# Keep up to 10 exemplars per metric (default: 5)
export OTEL_METRICS_EXEMPLAR_RESERVOIR_SIZE=10
```

## Usage

### Building with Exemplar Support

```bash
# Build your application with auto-instrumentation
otelbuild ./...

# Or use with go build
go build -toolexec="otelbuild" .
```

### Runtime Configuration Example

```bash
# Set up OpenTelemetry exporter
export OTEL_SERVICE_NAME="my-service"
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4317"

# Configure exemplars
export OTEL_GO_AUTO_EXEMPLARS_ENABLED=true
export OTEL_METRICS_EXEMPLAR_FILTER=trace_based
export OTEL_METRICS_EXEMPLAR_RESERVOIR_SIZE=10

# Run your application
./myapp
```

## Requirements

1. **OTLP Exporter**: Exemplars only work with OTLP exporters (gRPC or HTTP). Prometheus text format does not support exemplars.

2. **Trace Sampling**: When using `trace_based` filter, exemplars are only recorded for sampled traces.

3. **Backend Support**: Your observability backend must support exemplars (e.g., Prometheus with exemplar storage enabled, Grafana Tempo, etc.)

## How It Works

1. The auto-instrumentation captures trace context at function entry points
2. When metrics are recorded, the current trace context is associated with the measurement
3. If the trace is sampled and exemplar filters pass, an exemplar is created
4. The exemplar contains the metric value, timestamp, and trace/span IDs
5. Backend systems can then link from metric visualizations to the associated traces

## Example: HTTP Server Metrics with Exemplars

When your HTTP server handles requests, the auto-instrumentation will:

1. Start a trace span for the request
2. Capture the trace context
3. Record metrics (e.g., request duration) with the trace context
4. Create exemplars linking the metric to the specific trace

In your monitoring dashboard, you can then click on a latency spike and jump directly to the trace that caused it.

## Performance Considerations

- **Memory Usage**: Each exemplar consumes memory. Adjust `OTEL_METRICS_EXEMPLAR_RESERVOIR_SIZE` based on your needs.
- **CPU Overhead**: Use `trace_based` filter to minimize overhead by only recording exemplars for sampled traces.
- **Network**: Exemplars increase the size of metric exports. Monitor your network usage.

## Troubleshooting

### No Exemplars Appearing

1. Verify OTLP exporter is being used
2. Check if traces are being sampled
3. Ensure backend supports exemplars
4. Verify exemplar filter is not set to `always_off`

### High Memory Usage

1. Reduce `OTEL_METRICS_EXEMPLAR_RESERVOIR_SIZE`
2. Use `trace_based` filter instead of `always_on`
3. Reduce trace sampling rate

### Debugging

Enable debug logging to see exemplar activity:

```bash
export OTEL_LOG_LEVEL=debug
```

## Limitations

1. Only works with OTLP exporters
2. Prometheus text format does not support exemplars (use Remote Write instead)
3. Requires manual SDK implementation for custom metrics
4. Performance overhead when using `always_on` filter

## Future Improvements

- Support for custom exemplar filters
- Adaptive reservoir sizing
- Better integration with custom instrumentation
- Performance optimizations for high-throughput scenarios