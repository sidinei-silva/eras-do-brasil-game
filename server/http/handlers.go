package handlers

import (
	"fmt"
	"net/http"
)


func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}


func Handlers(mux *http.ServeMux) {
	mux.HandleFunc("/", helloHandler)
	adminClientFileServer := http.FileServer(http.Dir("../client/adminClient"))
	mux.Handle("/admin/", http.StripPrefix("/admin/", adminClientFileServer))
}

