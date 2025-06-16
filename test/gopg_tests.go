package test

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

const gopg_dependency_name = "github.com/go-pg/pg/v10"
const gopg_module_name = "gopg"

func init() {
	TestCases = append(TestCases, NewGeneralTestCase("test_gopg_crud", gopg_module_name, "v10.10.0", "v10.14.0", "1.19", "", TestGopgCrud10100))
}

func TestGopgCrud10100(t *testing.T, env ...string) {
	_, postgresPort := initPostgresContainer()
	UseApp("gopg/v10.10.0")
	RunGoBuild(t, "go", "build", "test_gopg_crud.go")
	env = append(env, "POSTGRES_PORT="+postgresPort.Port())
	RunApp(t, "test_gopg_crud", env...)
}

func initPostgresContainer() (testcontainers.Container, nat.Port) {
	containerReqeust := testcontainers.ContainerRequest{
		Image:        "postgres:4.0",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("waiting for connections")}
	postgresC, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{ContainerRequest: containerReqeust, Started: true})
	if err != nil {
		panic(err)
	}
	port, err := postgresC.MappedPort(context.Background(), "5432")
	if err != nil {
		panic(err)
	}
	return postgresC, port
}
