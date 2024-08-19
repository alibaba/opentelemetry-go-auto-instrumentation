//go:build ignore

package rule

import (
	"net/http"
	"net/url"
)

type ginRequest struct {
	method     string
	url        url.URL
	host       string
	isTls      bool
	header     http.Header
	version    string
	handleName string
}

type ginResponse struct {
	statusCode int
	header     http.Header
}

var netGinServerInstrument = BuildGinServerOtelInstrumenter()
