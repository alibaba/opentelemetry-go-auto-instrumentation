// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
