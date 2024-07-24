package http

type HttpStatusCodeConverter interface {
	IsError(statusCode int) bool
}

type ClientHttpStatusCodeConverter struct{}

func (c *ClientHttpStatusCodeConverter) IsError(statusCode int) bool {
	return statusCode >= 500 || statusCode < 100
}

type ServerHttpStatusCodeConverter struct{}

func (s *ServerHttpStatusCodeConverter) IsError(statusCode int) bool {
	return statusCode >= 400 || statusCode < 100
}
