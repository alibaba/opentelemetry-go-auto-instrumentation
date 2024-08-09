package http

type HttpClientSpanNameExtractor[REQUEST any, RESPONSE any] struct {
	Getter HttpClientAttrsGetter[REQUEST, RESPONSE]
}

func (h *HttpClientSpanNameExtractor[REQUEST, RESPONSE]) Extract(request REQUEST) string {
	method := h.Getter.GetRequestMethod(request)
	if method == "" {
		return "HTTP"
	}
	return method
}

type HttpServerSpanNameExtractor[REQUEST any, RESPONSE any] struct {
	Getter HttpServerAttrsGetter[REQUEST, RESPONSE]
}

func (h *HttpServerSpanNameExtractor[REQUEST, RESPONSE]) Extract(request REQUEST) string {
	method := h.Getter.GetRequestMethod(request)
	route := h.Getter.GetHttpRoute(request)
	if method == "" {
		return "HTTP"
	}
	if route == "" {
		return method
	}
	return method + " " + route
}
