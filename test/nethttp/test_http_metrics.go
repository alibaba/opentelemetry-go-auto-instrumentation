// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"net/http"
	"strconv"
	"time"
)

func setupMetricHttp() {
	http.HandleFunc("/a", helloHandler)
	var err error
	port, err = verifier.GetFreePort()
	if err != nil {
		panic(err)
	}
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	go setupMetricHttp()
	time.Sleep(1 * time.Second)
	_, err := http.Get("http://127.0.0.1:" + strconv.Itoa(port) + "/a")
	if err != nil {
		panic(err)
	}
	time.Sleep(3 * time.Second)
	verifier.WaitAndAssertMetrics(map[string]func(metricdata.ResourceMetrics){
		"http.server.request.duration": func(mrs metricdata.ResourceMetrics) {
			if len(mrs.ScopeMetrics) <= 0 {
				panic("No http.server.request.duration metrics received!")
			}
			point := mrs.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
			if point.DataPoints[0].Count <= 0 {
				panic("http.server.request.duration metrics count is not positive, actually " + strconv.Itoa(int(point.DataPoints[0].Count)))
			}
			verifier.VerifyHttpServerMetricsAttributes(point.DataPoints[0].Attributes.ToSlice(), "GET", "/a", "", "http", "1.1", "http", 200)
		},
		"http.client.request.duration": func(mrs metricdata.ResourceMetrics) {
			if len(mrs.ScopeMetrics) <= 0 {
				panic("No http.client.request.duration metrics received!")
			}
			point := mrs.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
			if point.DataPoints[0].Count <= 0 {
				panic("http.client.request.duration metrics count is not positive, actually " + strconv.Itoa(int(point.DataPoints[0].Count)))
			}
			verifier.VerifyHttpClientMetricsAttributes(point.DataPoints[0].Attributes.ToSlice(), "GET", "127.0.0.1:"+strconv.Itoa(port), "", "http", "1.1", port, 200)
		},
	})
}
