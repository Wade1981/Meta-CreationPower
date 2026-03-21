// Test script for network.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Mock NetworkManager for testing
type NetworkManager struct {
	port   int
	server *http.Server
}

// NewNetworkManager creates a new network manager
func NewNetworkManager(port int) *NetworkManager {
	return &NetworkManager{
		port: port,
	}
}

// Start starts the network service
func (n *NetworkManager) Start() error {
	// Set HTTP route
	handler := http.NewServeMux()
	
	// Health check
	handler.HandleFunc("/health", n.healthCheck)
	
	// API routes
	handler.HandleFunc("/api/container/list", n.listContainers)
	handler.HandleFunc("/api/model/run", n.runModel)
	
	// Create HTTP server
	serverAddr := fmt.Sprintf(":%d", n.port)
	n.server = &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}
	
	// Start server
	fmt.Printf("Network service starting on port %d\n", n.port)
	go func() {
		if err := n.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Network service error: %v\n", err)
		}
	}()
	
	return nil
}

// Stop stops the network service
func (n *NetworkManager) Stop() error {
	if n.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return n.server.Shutdown(ctx)
	}
	return nil
}

// healthCheck handles health check requests
func (n *NetworkManager) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","timestamp":%d,"service":"elr-network"}`, time.Now().Unix())
}

// listContainers handles container list requests
func (n *NetworkManager) listContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `[{"id":"test-container","name":"Test Container","status":"running"}]`)
}

// runModel handles model run requests
func (n *NetworkManager) runModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"output":"Model run successfully","timestamp":%d}`)
}

func main() {
	// Create network manager
	networkManager := NewNetworkManager(8080)

	// Start network service
	if err := networkManager.Start(); err != nil {
		fmt.Printf("Error starting network service: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Network service started successfully")
	fmt.Println("Network service is running on port 8080")
	fmt.Println("You can test it by visiting: http://localhost:8080/health")
	fmt.Println("Press Ctrl+C to stop")

	// Wait for signal to stop
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	// Stop network service
	if err := networkManager.Stop(); err != nil {
		fmt.Printf("Error stopping network service: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Network service stopped successfully")
}
