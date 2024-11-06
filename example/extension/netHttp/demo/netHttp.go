package main

import (
	"context"
	"net/http"
)

func main() {
	// 定义请求的URL
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://www.baidu.com", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("otelbuild", "true")
	client := &http.Client{}
	resp, err := client.Do(req)

	// 确保在函数结束时关闭响应的主体
	defer resp.Body.Close()

}
