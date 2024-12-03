// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
)

var Placeholder = ""

func main() {
	n, _ := fmt.Printf("helloworld%s", "ingodwetrust")
	println(n)
	println(fmt.Sprintf("placeholder:%s", Placeholder))
}
