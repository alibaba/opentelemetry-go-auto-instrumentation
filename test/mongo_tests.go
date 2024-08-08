package test

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/version"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"testing"
)

const dependency_name = "go.mongodb.org/mongo-driver"
const module_name = "mongo"

func init() {
	TestCases = append(TestCases, NewGeneralTestCase("mongo-1.11.1-test", dependency_name, module_name, "v1.11.1", "v1.15.1", "1.18", "", TestGoMongo111),
		NewMuzzleTestCase("mongo-1.11.1-muzzle", dependency_name, module_name, "v1.11.1", "v1.15.1", "1.18", ""),
		NewLatestDepthTestCase("mongo-1.11.1-latestDepth", dependency_name, module_name, "v1.11.1", "v1.15.1", "1.18", "", TestMongoLatest))
}

func TestGoMongo111(t *testing.T, env ...string) {
	mongoC, mongoPort := initMongoContainer()
	defer clearMongoContainer(mongoC)
	UseApp("mongo/v1.11.1")
	RunInstrument(t, "-debuglog")
	env = append(env, "MONGO_PORT="+mongoPort.Port())
	RunApp(t, "v1.11.1", env...)
}

func TestMongoLatest(t *testing.T, v *version.Version, env ...string) {
	mongoC, mongoPort := initMongoContainer()
	defer clearMongoContainer(mongoC)
	RunInstrument(t, "-debuglog")
	env = append(env, "MONGO_PORT="+mongoPort.Port())
	RunApp(t, v.Original(), env...)
}

func initMongoContainer() (testcontainers.Container, nat.Port) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:4.0",
		ExposedPorts: []string{"27017/tcp"},
	}
	mongoC, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	port, err := mongoC.MappedPort(context.Background(), "27017")
	if err != nil {
		panic(err)
	}
	return mongoC, port
}

func clearMongoContainer(mongoC testcontainers.Container) {
	if err := mongoC.Terminate(context.Background()); err != nil {
		log.Fatal(err)
	}
}
