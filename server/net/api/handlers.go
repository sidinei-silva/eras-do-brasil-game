package api

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func Handlers(mux *http.ServeMux) {
	adminClientFileServer := http.FileServer(http.Dir("../client/adminClient"))

	mux.HandleFunc("/", helloHandler)
	mux.Handle("/admin/", http.StripPrefix("/admin/", adminClientFileServer))
}
