// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"example/demo/pkg"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go func() {
		pkg.InitDB()
		pkg.SetupHttp()
	}()

	http.ListenAndServe("0.0.0.0:6060", nil)

	signalCh := make(chan os.Signal, 1)

	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	<-signalCh

	os.Exit(0)
}
