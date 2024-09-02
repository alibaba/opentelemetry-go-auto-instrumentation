/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	//"github.com/dubbogo/gost/log/logger"
)

var grpcGreeterImpl = new(GreeterClientImpl)

type GreeterProvider struct {
	UnimplementedGreeterServer
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *HelloRequest) (*User, error) {
	//logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &User{Name: "Hello " + in.Name, Id: "1234578910", Age: 21}, nil
}

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/helloworld/go-server/conf/dubbogo.yaml
func RunServer() {
	config.SetProviderService(&GreeterProvider{})
	config.SetConsumerService(grpcGreeterImpl)
	if err := config.Load(config.WithPath("./dubbogo.yaml")); err != nil {
		panic(err)
	}
	select {}
}

func RunClient() {
	//config.SetConsumerService(grpcGreeterImpl)
	/*	if err := config.Load(config.WithPath("../v3.1.1/conf/dubbogo.yaml")); err != nil {
		panic(err)
	}*/

	//logger.Info("start to test dubbo")
	req := &HelloRequest{
		Name: "laurence",
	}
	retry, err := grpcGreeterImpl.SayHello(context.Background(), req)
	if err != nil {
		//	logger.Error(err)
	}
	xx, _ := json.Marshal(retry)
	fmt.Println(string(xx))
	//	logger.Infof("client response result: %v\n", reply)
}
