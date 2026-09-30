// main.go - a tiny web server we will dockerize
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Read the port from an environment variable, fallback to 8080.
	// This is a 12-factor best practice: config comes from the environment.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// The Vietnamese message below keeps full diacritics.
		fmt.Fprintln(w, "Xin chào từ ứng dụng Go chạy trong Docker!")

		// Get hostname (which is the container ID)
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
		fmt.Fprintf(w, "Hostname: %s\n", hostname)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	addr := ":" + port
	log.Printf("server listening on %s", addr)
	// log.Fatal stops the program if the server fails to start.
	log.Fatal(http.ListenAndServe(addr, nil))
}
