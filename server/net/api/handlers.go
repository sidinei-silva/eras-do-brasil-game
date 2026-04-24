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
	adminClient2FileServer := http.FileServer(http.Dir("../client/adminClient2"))

	mux.HandleFunc("/", helloHandler)
	mux.Handle("/admin/", http.StripPrefix("/admin/", adminClientFileServer))
	mux.Handle("/admin2/", http.StripPrefix("/admin2/", adminClient2FileServer))
}
