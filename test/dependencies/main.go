package main

import (
	"fmt"
	_ "github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> helloHandler called")
	fmt.Fprintf(w, "Hello, World!\n")
}

func main() {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()

	resp, _ := http.Get(ts.URL)
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
}
