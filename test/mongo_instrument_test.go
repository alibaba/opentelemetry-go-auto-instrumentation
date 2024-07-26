package test

import (
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
)

const dependency_name = "go.mongodb.org/mongo-driver"
const module_name = "mongo"

var min_version *Version
var max_version *Version

var min_go_version *Version
var max_go_version *Version

var mongoC testcontainers.Container

var mongoPort nat.Port

func TestMain(m *testing.M) {
	runtime.Version()
	min_version, _ = NewVersion("v1.11.1")
	max_version, _ = NewVersion("v1.15.1")
	min_go_version, _ = NewVersion("1.18")
	max_go_version, _ = NewVersion("1.21.12")
	goVersion, _ := NewVersion(strings.ReplaceAll(runtime.Version(), "go", ""))
	if goVersion.LessThan(min_go_version) || goVersion.GreaterThan(max_go_version) {
		log.Printf("This test does not suppport go " + goVersion.String())
		return
	}
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:4.0",
		ExposedPorts: []string{"27017/tcp"},
	}
	var err error
	mongoC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start mongo: %s", err)
	}
	mongoPort, err = mongoC.MappedPort(ctx, "27017")
	defer func() {
		if err := mongoC.Terminate(context.Background()); err != nil {
			log.Fatalf("Could not stop mongo: %s", err)
		}
	}()
	r := m.Run()
	os.Exit(r)
}

func TestGoMongo111(t *testing.T) {
	testMongo(t, "v1.11.1")
}

func TestMongoMuzzle(t *testing.T) {
	ExecMuzzle(t, dependency_name, module_name, min_version, max_version)
}

func TestMongoLatest(t *testing.T) {
	ExecLatestTest(t, dependency_name, module_name, min_version, max_version, func(t *testing.T, v *Version) {
		RunInstrument(t, "-debuglog")
		RunApp(t, v.Original(), "MONGO_PORT="+mongoPort.Port())
	})
}

func testMongo(t *testing.T, version string) {
	UseApp("mongo/" + version)
	RunInstrument(t, "-debuglog")
	RunApp(t, version, "MONGO_PORT="+mongoPort.Port())
}
