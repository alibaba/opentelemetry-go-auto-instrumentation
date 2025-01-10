// Copyright (c) 2025 Alibaba Group Holding Ltd.
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
