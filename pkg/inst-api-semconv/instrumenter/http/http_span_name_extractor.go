package http

type HttpClientSpanNameExtractor[REQUEST any, RESPONSE any] struct {
	getter HttpClientAttrsGetter[REQUEST, RESPONSE]
}

func (h *HttpClientSpanNameExtractor[REQUEST, RESPONSE]) Extract(request REQUEST) string {
	method := h.getter.GetRequestMethod(request)
	if method == "" {
		return "HTTP"
	}
	return method
}

type HttpServerSpanNameExtractor[REQUEST any, RESPONSE any] struct {
	getter HttpServerAttrsGetter[REQUEST, RESPONSE]
}

func (h *HttpServerSpanNameExtractor[REQUEST, RESPONSE]) Extract(request REQUEST) string {
	method := h.getter.GetRequestMethod(request)
	route := h.getter.GetHttpRoute(request)
	if method == "" {
		return "HTTP"
	}
	if route == "" {
		return method
	}
	return method + " " + route
}
