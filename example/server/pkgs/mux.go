package pkgs

import (
	"github.com/gorilla/mux"
	"net/http"
)

func SetMux() {
	r := mux.NewRouter()

	r.HandleFunc("/mux-service1", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello Mux"))
	})

	http.Handle("/", r)
	http.ListenAndServe(":9004", nil)
}
