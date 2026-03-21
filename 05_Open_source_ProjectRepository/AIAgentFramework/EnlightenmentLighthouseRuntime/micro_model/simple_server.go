package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"simple-server"}`))
	})

	serverAddr := "127.0.0.1:9003"
	fmt.Printf("Simple server starting on %s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		fmt.Printf("Failed to start simple server: %v\n", err)
	}
}
