// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package dep

func BadDep() string {
	return "baddep"
}

func init() {
	println(BadDep())
}
