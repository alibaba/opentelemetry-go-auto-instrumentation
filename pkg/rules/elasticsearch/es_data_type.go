// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package elasticsearch

import "net/http"

type esRequest struct {
	request *http.Request
	address string
	op      string
	params  []any
}
