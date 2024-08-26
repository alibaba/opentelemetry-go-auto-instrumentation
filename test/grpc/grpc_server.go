package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendReq(ctx context.Context) string {
	conn, err := grpc.NewClient("localhost:9003", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	c := NewHelloGrpcClient(conn)

	resp, err := c.Hello(ctx, &Req{})
	if err != nil {
		panic(err)
	}
	return resp.Message
}
