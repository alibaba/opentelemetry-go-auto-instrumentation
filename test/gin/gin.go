package main

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"net/http"
	"strconv"
	"time"
)

func RequestInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.FullPath()
		method := context.Request.Method
		fmt.Println("请求路径为:", path, "请求方法:", method)
		context.Next()
	}
}

var port int

func setupHttp() {
	engine := gin.Default()
	engine.LoadHTMLGlob("templates/*")
	engine.Use(RequestInfo())
	engine.GET("/query", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Gin HTML模板示例",
		})
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 1,
			"msg":  c.FullPath(),
		})

	})
	var err error
	port, err = verifier.GetFreePort()
	if err != nil {
		panic(err)
	}
	engine.Run(":" + strconv.Itoa(port))
}

func requestServer() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:"+strconv.Itoa(port)+"/query", nil)
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

func main() {
	// starter server
	go setupHttp()
	time.Sleep(3 * time.Second)
	// use a http client to request to the server
	requestServer()
	// verify trace
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyHttpClientAttributes(stubs[0][0], "GET", "GET", "http://127.0.0.1:"+strconv.Itoa(port)+"/query", "http", "1.1", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), 200, 0, int64(port))
		verifier.VerifyHttpServerAttributes(stubs[0][1], "GET /query", "GET", "http", "tcp", "ipv4", "", "127.0.0.1:"+strconv.Itoa(port), "Go-http-client/1.1", "", "/query", "", "/query", 200)
		if stubs[0][1].Parent.TraceID().String() != stubs[0][0].SpanContext.TraceID().String() {
			log.Fatal("span 1 should be child of span 0")
		}
	}, 1)
}
