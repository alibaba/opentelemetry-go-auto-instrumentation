package pkgs

import (
	"io"
	"net/http"
)

func doService(w http.ResponseWriter, r *http.Request, downstream string) {
	c := r.Context()
	req, err := http.NewRequestWithContext(c, "GET", downstream, nil)
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
	w.Write(b)
}

func service1(w http.ResponseWriter, r *http.Request) {
	doService(w, r, "http://localhost:9000/http-service2")
}

func service2(w http.ResponseWriter, r *http.Request) {
	doService(w, r, "http://localhost:9000/http-service3")
}

func service3(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello Http!"))
	if err != nil {
		panic(err)
	}
}

func SetupHttp() {
	http.Handle("/http-service1", http.HandlerFunc(service1))
	http.Handle("/http-service2", http.HandlerFunc(service2))
	http.Handle("/http-service3", http.HandlerFunc(service3))

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}
}
