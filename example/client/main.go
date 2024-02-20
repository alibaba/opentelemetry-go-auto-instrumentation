package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net/http"
	"time"
)

func main() {
	{
		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:9000/http-service1", nil)
		if err != nil {
			panic(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}

	{
		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:9001/gin-service1", nil)
		if err != nil {
			panic(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}

	{
		ctx := context.Background()
		conn, err := grpc.Dial("127.0.0.1:9003", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		client := NewHelloGrpcClient(conn)
		resp, err := client.Hello(ctx, &Req{})
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.Message)
	}

	time.Sleep(time.Second * 3)
}
