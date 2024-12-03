// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

type service struct {
	HelloGrpcServer
}

func (s *service) Hello(ctx context.Context, req *Req) (*Resp, error) {
	if req.Error {
		return &Resp{Message: "Error Grpc"}, errors.New("error Grpc")
	}
	return &Resp{Message: "Hello Gprc"}, nil
}

func setupGRPC() {
	lis, err := net.Listen("tcp", "0.0.0.0:9003")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterHelloGrpcServer(s, &service{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func sendReq(ctx context.Context) string {
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

func sendErrReq(ctx context.Context) string {
	conn, err := grpc.NewClient("localhost:9003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	c := NewHelloGrpcClient(conn)
	resp, err := c.Hello(ctx, &Req{Error: true})
	if err != nil {
		log.Printf("error %v\n", err)
	}
	if resp != nil {
		return resp.Message
	}
	return "error resp"
}
