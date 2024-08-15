package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func RunClient() {
	logrus.Warn("Warn msg 12123123") // not be written
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:10138/http-service1", nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	fmt.Println(resp.Header)

}

func SetupHttp() {
	http.Handle("/http-service1", http.HandlerFunc(service))

	err := http.ListenAndServe(":10138", nil)
	if err != nil {
		panic(err)
	}
}

func service(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	for key, value := range header {
		values := strings.Join(value, ",")
		fmt.Printf("[Headers] Key is %s\tValue is %s\n", key, values)
	}

	_, err := w.Write([]byte("Hello Http!"))
	if err != nil {
		panic(err)
	}
}

func main() {
	go func() {
		SetupHttp()
	}()
	RunClient()
	logrus.Warn("Warn msg")   // written in general.log
	logrus.Error("Error msg") // written in error.log
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel) // default

	//logrus.AddHook(&emailHook{})
}
