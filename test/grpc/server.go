package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type service struct {
	HelloGrpcServer
}

func (s *service) Hello(context.Context, *Req) (*Resp, error) {
	return &Resp{Message: "Hello Gprc"}, nil
}

func SetupGRPC() {

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
