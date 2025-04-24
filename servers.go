package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// startServer starts a simple HTTP server on a given port with multiple routes.
func startServer(port string) {
	mux := http.NewServeMux() // multiplex router for different paths

	// Define multiple endpoints
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server running on port %s\n", port)
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server %s is healthy!\n", port)
	})

	// log.Printf("Server started on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func MyServers() {
	ports := []string{"5001", "5002", "5003"}
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			startServer(p)
		}(port)
	}

	wg.Wait() // Keep the main function alive
}
