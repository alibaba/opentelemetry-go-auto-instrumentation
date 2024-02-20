package pkgs

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type service struct {
	HelloGrpcServer
}

func (s *service) Hello(context.Context, *Req) (*Resp, error) {
	return &Resp{Message: "Hello Gprc"}, nil
}

func sendReq(ctx context.Context) string {
	conn, err := grpc.Dial("127.0.0.1:9003", grpc.WithTransportCredentials(insecure.NewCredentials()))

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

func SetupGRPC() {
	go func() {
		http.Handle("/grpc-service1", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte(sendReq(request.Context())))
		}))
		err := http.ListenAndServe(":9002", nil)
		if err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen("tcp", "127.0.0.1:9003")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	RegisterHelloGrpcServer(s, &service{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
