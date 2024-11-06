package main

import (
	"context"
	"net/http"
)

func main() {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://www.baidu.com", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("otelbuild", "true")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

}
