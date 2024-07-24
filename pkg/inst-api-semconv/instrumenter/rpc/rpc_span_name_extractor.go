package rpc

type RpcSpanNameExtractor[REQUEST any] struct {
	getter RpcAttrsGetter[REQUEST]
}

func (r *RpcSpanNameExtractor[REQUEST]) Extract(request REQUEST) string {
	service := r.getter.GetService(request)
	method := r.getter.GetMethod(request)
	if service == "" || method == "" {
		return "RPC request"
	}
	return service + "/" + method
}
