package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"test-server"}`))
	})

	fmt.Println("Test server started on 127.0.0.1:9002")
	if err := http.ListenAndServe("127.0.0.1:9002", nil); err != nil {
		fmt.Printf("Failed to start test server: %v\n", err)
	}
}
