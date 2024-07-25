package http

type HttpClientSpanNameExtractor[REQUEST any] struct {
	getter HttpClientAttrsGetter[REQUEST, any]
}

func (h *HttpClientSpanNameExtractor[REQUEST]) Extract(request REQUEST) string {
	method := h.getter.GetRequestMethod(request)
	if method == "" {
		return "HTTP"
	}
	return method
}

type HttpServerSpanNameExtractor[REQUEST any] struct {
	getter HttpServerAttrsGetter[REQUEST, any]
}

func (h *HttpServerSpanNameExtractor[REQUEST]) Extract(request REQUEST) string {
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
