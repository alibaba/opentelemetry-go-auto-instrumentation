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

package test

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

const rocketmq_dependency_name = "github.com/apache/rocketmq-client-go/v2"
const rocketmq_module_name = "rocketmq"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("rocketmq_basic-2.1.2-test", rocketmq_module_name, "2.0.0", "", "1.18", "", TestRocketMQBasic),
		NewGeneralTestCase("rocketmq_producer-2.1.2-test", rocketmq_module_name, "2.0.0", "", "1.18", "", TestRocketMQProducer),
		NewGeneralTestCase("rocketmq_consumer-2.1.2-test", rocketmq_module_name, "2.0.0", "", "1.18", "", TestRocketMQConsumer),
	)
}

func TestRocketMQBasic(t *testing.T, env ...string) {
	containers := initRocketMQContainer(t)
	defer containers.CleanupContainers(context.Background())
	UseApp("rocketmq/v2.1.2")
	RunGoBuild(t, "go", "build", "test_mq_basic.go", "base.go")
	env = append(env, "NAMESRV_ADDR="+containers.NameSrvAddr)
	env = append(env, "BROKER_ADDR="+containers.BrokerAddr)
	RunApp(t, "test_mq_basic", env...)
}

func TestRocketMQProducer(t *testing.T, env ...string) {
	containers := initRocketMQContainer(t)
	defer containers.CleanupContainers(context.Background())

	UseApp("rocketmq/v2.1.2")
	RunGoBuild(t, "go", "build", "test_mq_producer.go", "base.go")
	env = append(env, "NAMESRV_ADDR="+containers.NameSrvAddr)
	env = append(env, "BROKER_ADDR="+containers.BrokerAddr)
	RunApp(t, "test_mq_producer", env...)
}

func TestRocketMQConsumer(t *testing.T, env ...string) {
	containers := initRocketMQContainer(t)
	defer containers.CleanupContainers(context.Background())

	UseApp("rocketmq/v2.1.2")
	RunGoBuild(t, "go", "build", "test_mq_consumer.go", "base.go")
	env = append(env, "NAMESRV_ADDR="+containers.NameSrvAddr)
	env = append(env, "BROKER_ADDR="+containers.BrokerAddr)
	RunApp(t, "test_mq_consumer", env...)
}

// RocketMQContainers 封装RocketMQ相关容器，便于统一管理
type RocketMQContainers struct {
	NameSrvContainer testcontainers.Container
	BrokerContainer  testcontainers.Container
	NameSrvAddr      string
	BrokerAddr       string
}

// CleanupContainers 清理所有RocketMQ容器资源
func (r *RocketMQContainers) CleanupContainers(ctx context.Context) {
	if r.BrokerContainer != nil {
		_ = r.BrokerContainer.Terminate(ctx)
	}
	if r.NameSrvContainer != nil {
		_ = r.NameSrvContainer.Terminate(ctx)
	}
}

// initRocketMQContainer 初始化RocketMQ测试环境
func initRocketMQContainer(t *testing.T) *RocketMQContainers {
	ctx := context.Background()
	testNetwork, err := network.New(ctx)
	// 启动NameServer容器
	nameServerReq := testcontainers.ContainerRequest{
		Image:        "apache/rocketmq:4.9.4",
		ExposedPorts: []string{"9876/tcp"},
		Cmd:          []string{"sh", "-c", "/home/rocketmq/rocketmq-4.9.4/bin/mqnamesrv"},
		WaitingFor:   wait.ForLog("Name Server boot success.").WithStartupTimeout(30 * time.Second),
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				"9876/tcp": []nat.PortBinding{{
					HostIP:   "0.0.0.0",
					HostPort: "9876",
				}},
			}
		},
		Networks: []string{testNetwork.Name},
		NetworkAliases: map[string][]string{
			testNetwork.Name: {"rocketmq-nameserver"},
		},
	}

	nameServerC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: nameServerReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start RocketMQ NameServer container: %v", err)
	}

	// 获取NameServer地址
	host := "127.0.0.1"
	//port, err := nameServerC.MappedPort(ctx, "9876")
	//if err != nil {
	//	nameServerC.Terminate(ctx)
	//	t.Fatalf("Failed to get NameServer port: %v", err)
	//}

	nameSrvAddr := host + ":" + "9876"

	// 创建broker配置
	brokerConf := `brokerClusterName=DefaultCluster
brokerName=broker-a
brokerId=0
autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true
deleteWhen=04
fileReservedTime=48
brokerRole=ASYNC_MASTER
flushDiskType=ASYNC_FLUSH
aclEnable=false
brokerIP1=` + host
	// 启动Broker容器
	brokerReq := testcontainers.ContainerRequest{
		Image:        "apache/rocketmq:4.9.4",
		ExposedPorts: []string{"10911/tcp", "10909/tcp"},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				"10911/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "10911"}},
				"10909/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "10909"}},
				"10910/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "10910"}},
			}
		},
		Env: map[string]string{
			"NAMESRV_ADDR": "rocketmq-nameserver:9876",
		},
		Networks: []string{testNetwork.Name},
		NetworkAliases: map[string][]string{
			testNetwork.Name: {"rocketmq-broker"},
		},
		Cmd: []string{"sh", "-c",
			"echo '" + brokerConf + "' > /home/rocketmq/rocketmq-4.9.4/conf/broker.conf && " +
				"/home/rocketmq/rocketmq-4.9.4/bin/mqbroker -n rocketmq-nameserver:9876 -c /home/rocketmq/rocketmq-4.9.4/conf/broker.conf"},
		WaitingFor: wait.ForLog("The broker[broker-a").WithStartupTimeout(60 * time.Second),
	}

	brokerC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: brokerReq,
		Started:          true,
	})
	if err != nil {
		nameServerC.Terminate(ctx)
		t.Fatalf("Failed to start Broker container: %v", err)
	}
	// 获取Broker地址
	//brokerPort, err := brokerC.MappedPort(ctx, "10911")
	//if err != nil {
	//	nameServerC.Terminate(ctx)
	//	brokerC.Terminate(ctx)
	//	t.Fatalf("获取Broker端口失败: %v", err)
	//}
	brokerAddr := host + ":" + "10911"

	// 等待Broker完全初始化
	time.Sleep(5 * time.Second)

	containers := &RocketMQContainers{
		NameSrvContainer: nameServerC,
		BrokerContainer:  brokerC,
		NameSrvAddr:      nameSrvAddr,
		BrokerAddr:       brokerAddr,
	}

	return containers
}
