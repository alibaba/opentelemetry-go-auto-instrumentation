package main

import (
	"encoding/json"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func requestServer() {
	client := &fasthttp.Client{}

	reqURL := "http://localhost:3000/fiber"

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	req.SetRequestURI(reqURL)
	req.Header.SetMethod(fasthttp.MethodGet)

	err := client.Do(req, resp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Response body is:\n%s", resp.Body())
}

func setupHttp() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/fiber", func(c *fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}

func main() {
	// starter server
	go setupHttp()
	time.Sleep(3 * time.Second)
	// use a http client to request to the server
	requestServer()
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		xx, _ := json.Marshal(stubs)
		fmt.Println(string(xx))
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://localhost:3000/fiber", "http", "", "tcp", "ipv4", "", "localhost:3000", 200, 0, 3000)
		verifier.VerifyHttpServerAttributes(stubs[0][1], "GET /fiber", "GET", "http", "tcp", "ipv4", "", "localhost:3000", "fasthttp", "http", "/fiber", "", "/fiber", 200)
	}, 1)
}
