package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://127.0.0.1:" + strconv.Itoa(port) + "/b")
	if err != nil {
		log.Printf("request provider error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("success"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
}
