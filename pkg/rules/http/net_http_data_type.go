// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package http

import (
	"net/http"
	"net/url"
	"strconv"
)

type netHttpRequest struct {
	method  string
	url     *url.URL
	host    string
	isTls   bool
	header  http.Header
	version string
}

type netHttpResponse struct {
	statusCode int
	header     http.Header
}

func getProtocolVersion(majorVersion, minorVersion int) string {
	if majorVersion == 1 && minorVersion == 0 {
		return "1.0"
	} else if majorVersion == 1 && minorVersion == 1 {
		return "1.1"
	} else if majorVersion == 1 && minorVersion == 2 {
		return "1.2"
	} else if majorVersion == 2 {
		return "2"
	} else if majorVersion == 3 {
		return "3"
	}
	return strconv.Itoa(majorVersion) + "." + strconv.Itoa(minorVersion)
}
