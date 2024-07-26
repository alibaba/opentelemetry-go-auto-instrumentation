package verifier

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"net"
	"net/http"
	"os"
)

const IS_IN_TEST = "IN_OTEL_TEST"

func GetAttribute(attrs []attribute.KeyValue, name string) attribute.Value {
	for _, attr := range attrs {
		if string(attr.Key) == name {
			return attr.Value
		}
	}
	return attribute.Value{}
}

func Assert(cond bool, format string, args ...interface{}) {
	if !cond {
		panic(fmt.Sprintf(format, args...))
	}
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		panic("Failed to create a free port: " + err.Error())
	}
	cli, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic("Failed to create a free port: " + err.Error())
	}
	defer cli.Close()
	return cli.Addr().(*net.TCPAddr).Port, nil
}

func GetServer(ctx context.Context, url string) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func IsInTest() bool {
	return os.Getenv(IS_IN_TEST) == "true"
}
