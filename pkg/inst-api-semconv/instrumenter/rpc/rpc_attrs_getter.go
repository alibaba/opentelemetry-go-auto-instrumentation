package rpc

type RpcAttrsGetter[REQUEST any] interface {
	GetSystem(request REQUEST) string
	GetService(request REQUEST) string
	GetMethod(request REQUEST) string
}
