package test

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"testing"
	"time"
)

const rabbitmq_dependency_name = "https://github.com/rabbitmq/amqp091-go"
const rabbitmq_module_name = "amqp091go"

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("rabbitmq_cascading-1.10.0-test", rabbitmq_module_name, "1.10.0", "1.10.0", "1.22.0", "", TestRabbitMQCascading),
		NewGeneralTestCase("rabbitmq_no_cascading-1.10.0-test", rabbitmq_module_name, "1.10.0", "1.10.0", "1.22.0", "", TestRabbitMQNOCascading),
	)

}
func TestRabbitMQCascading(t *testing.T, env ...string) {
	rabbitC, port := initRabbitMQContainer()
	defer testcontainers.CleanupContainer(t, rabbitC)
	UseApp("amqp091go/v1.10.0")
	RunGoBuild(t, "go", "build", "test_mq_cascading.go", "base.go")
	env = append(env, "RabbitMQ_PORT="+port.Port())
	RunApp(t, "test_mq_cascading", env...)
}
func TestRabbitMQNOCascading(t *testing.T, env ...string) {
	rabbitC, port := initRabbitMQContainer()
	defer testcontainers.CleanupContainer(t, rabbitC)
	UseApp("amqp091go/v1.10.0")
	RunGoBuild(t, "go", "build", "test_mq_no_cascading.go", "base.go")
	env = append(env, "RabbitMQ_PORT="+port.Port())
	RunApp(t, "test_mq_no_cascading", env...)
}
func initRabbitMQContainer() (testcontainers.Container, nat.Port) {
	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:4.0.7-alpine",
		ExposedPorts: []string{"5672/tcp"},
	}
	rabbitC, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	port, err := rabbitC.MappedPort(context.Background(), "5672")
	if err != nil {
		panic(err)
	}
	return rabbitC, port
}
