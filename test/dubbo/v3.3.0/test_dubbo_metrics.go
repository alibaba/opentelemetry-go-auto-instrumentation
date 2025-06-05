package main

import (
	"context"
	"strconv"
	"time"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func main() {
	// starter server
	go setupDubbo()
	time.Sleep(3 * time.Second)
	// use a client to request to the server
	sendBasicDubboReq(context.Background())
	// verify metrics
	verifier.WaitAndAssertMetrics(map[string]func(metricdata.ResourceMetrics){
		"rpc.server.duration": func(mrs metricdata.ResourceMetrics) {
			if len(mrs.ScopeMetrics) <= 0 {
				panic("No rpc.server.duration metrics received!")
			}
			point := mrs.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
			if point.DataPoints[0].Count <= 0 {
				panic("rpc.server.duration metrics count is not positive, actually " + strconv.Itoa(int(point.DataPoints[0].Count)))
			}
			verifier.VerifyRpcServerMetricsAttributes(point.DataPoints[0].Attributes.ToSlice(), "Greet", "greet.GreetService", "apache_dubbo", ":20000")
		},
		"rpc.client.duration": func(rm metricdata.ResourceMetrics) {
			if len(rm.ScopeMetrics) <= 0 {
				panic("No rpc.client.duration metrics received!")
			}
			point := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
			if point.DataPoints[0].Count <= 0 {
				panic("rpc.client.duration metrics count is not positive, actually " + strconv.Itoa(int(point.DataPoints[0].Count)))
			}
			verifier.VerifyRpcClientMetricsAttributes(point.DataPoints[0].Attributes.ToSlice(), "Greet", "greet.GreetService", "apache_dubbo", "127.0.0.1:20000")
		},
	})
}
