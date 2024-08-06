//go:build ignore

package rule

import (
	"context"
	"net/http"
	"strconv"
)

var netHttpClientInstrumenter = BuildNetHttpClientOtelInstrumenter()

func clientOnEnter(call http.CallContext, t *http.Transport, req *http.Request) {
	ctx := netHttpClientInstrumenter.Start(req.Context(), netHttpRequest{
		method:  req.Method,
		url:     *req.URL,
		header:  req.Header,
		version: strconv.Itoa(req.ProtoMajor) + "." + strconv.Itoa(req.ProtoMinor),
	})
	req = req.WithContext(ctx)
	call.SetParam(1, req)
	data := make(map[string]interface{}, 1)
	data["ctx"] = ctx
	call.SetData(data)
	return
}

func clientOnExit(call http.CallContext, res *http.Response, err error) {
	data, ok := call.GetData().(map[string]interface{})
	if !ok || data == nil || data["ctx"] == nil {
		return
	}
	ctx := data["ctx"].(context.Context)
	if res != nil {
		netHttpClientInstrumenter.End(ctx, netHttpRequest{
			method: res.Request.Method,
			url:    *res.Request.URL,
			header: res.Request.Header,
		}, netHttpResponse{
			statusCode: res.StatusCode,
			header:     res.Header,
		}, err)
	} else {
		netHttpClientInstrumenter.End(ctx, netHttpRequest{}, netHttpResponse{
			statusCode: 500,
		}, err)
	}
}
