package main

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	pb "kratos/api/helloworld/v1"
	"log"
)

func callHTTP() {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint("127.0.0.1:8777"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "client"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)

	/*// returns error
	_, err = client.SayHello(context.Background(), &pb.HelloRequest{Name: "error"})
	if err != nil {
		log.Printf("[http] SayHello error: %v\n", err)
	}
	if errors.IsBadRequest(err) {
		log.Printf("[http] SayHello error is invalid argument: %v\n", err)
	}*/
}

func callGRPC() {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9777"),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "client"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)

	// returns error
	/*_, err = client.SayHello(context.Background(), &pb.HelloRequest{Name: "error"})
	if err != nil {
		log.Printf("[grpc] SayHello error: %v\n", err)
	}
	if errors.IsBadRequest(err) {
		log.Printf("[grpc] SayHello error is invalid argument: %v\n", err)
	}*/
}

func RunClient() {
	callHTTP()
	//time.Sleep(15 * time.Second)
}

func RunGrpcClient() {
	callGRPC()
	//time.Sleep(15 * time.Second)
}
