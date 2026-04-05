package main

import (
	"fmt"
	"net/http"
)


func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}


func handlers(mux *http.ServeMux) {
	mux.HandleFunc("/", helloHandler)
	adminClientFileServer := http.FileServer(http.Dir("../client/adminClient"))
	mux.Handle("/admin/", http.StripPrefix("/admin/", adminClientFileServer))
}

func main() {
	mux := http.NewServeMux()
	handlers(mux)
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}