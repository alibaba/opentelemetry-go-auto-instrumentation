//go:build ignore

package rule

import (
	"net/http"
	"net/url"
)

type netHttpRequest struct {
	method  string
	url     url.URL
	header  http.Header
	version string
}

type netHttpResponse struct {
	statusCode int
	header     http.Header
}
