package main

import (
	"github.com/go-micro/examples/server/handler"
	"github.com/go-micro/examples/server/subscriber"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/server"
	"log"
)

func setupHttp() {
	// optionally setup command line usage
	cmd.Init()

	// Initialise Server
	server.Init(
		server.Name("go.micro.srv.example"),
	)

	// Register Handlers
	server.Handle(
		server.NewHandler(
			new(handler.Example),
		),
	)

	// Register Subscribers
	if err := server.Subscribe(
		server.NewSubscriber(
			"topic.example",
			new(subscriber.Example),
		),
	); err != nil {
		log.Fatal(err)
	}

	/*if err := server.Subscribe(
		server.NewSubscriber(
			"topic.example",
			subscriber.Handler,
		),
	); err != nil {
		log.Fatal(err)
	}*/

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
