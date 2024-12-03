// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"time"

	_ "go.opentelemetry.io/otel"
	"golang.org/x/time/rate"
)

func main() {
	n, _ := fmt.Printf("helloworld%s", "ingodwetrust")
	println(n)

	println(rate.Every(time.Duration(1) * time.Second))
}
