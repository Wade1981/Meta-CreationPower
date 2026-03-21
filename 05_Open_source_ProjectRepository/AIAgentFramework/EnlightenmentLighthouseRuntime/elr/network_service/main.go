package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type NetworkManager struct {
	port   int
	server *http.Server
}

func NewNetworkManager(port int) *NetworkManager {
	return &NetworkManager{
		port: port,
	}
}

func (n *NetworkManager) Start() error {
	handler := http.NewServeMux()
	
	handler.HandleFunc("/health", n.healthCheck)
	handler.HandleFunc("/api/status", n.getStatus)
	handler.HandleFunc("/api/network/status", n.getNetworkStatus)
	handler.HandleFunc("/api/container/list", n.listContainers)
	handler.HandleFunc("/api/model/list", n.listModels)
	
	serverAddr := fmt.Sprintf(":%d", n.port)
	n.server = &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}
	
	fmt.Printf("ELR Network Service starting on port %d\n", n.port)
	fmt.Printf("Address: http://localhost:%d\n", n.port)
	
	go func() {
		if err := n.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Network service error: %v\n", err)
		}
	}()
	
	return nil
}

func (n *NetworkManager) Stop() error {
	if n.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return n.server.Shutdown(ctx)
	}
	return nil
}

func (n *NetworkManager) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "elr-network-service",
		"version":   "1.0.0",
	})
}

func (n *NetworkManager) getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "running",
		"message": "ELR Network Service is running",
	})
}

func (n *NetworkManager) getNetworkStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"desktop_api": map[string]interface{}{
			"address": "http://localhost:8081",
			"port":    8081,
			"status":  "check",
		},
		"public_api": map[string]interface{}{
			"address": "http://localhost:8080",
			"port":    8080,
			"status":  "running",
		},
		"model_service": map[string]interface{}{
			"address": "http://localhost:8082",
			"port":    8082,
			"status":  "check",
		},
		"micro_model_server": map[string]interface{}{
			"address": "http://localhost:8083",
			"port":    8083,
			"status":  "check",
		},
		"timestamp": time.Now().Unix(),
	})
}

func (n *NetworkManager) listContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{
			"id":     "elr-1234567890",
			"name":   "test-container",
			"image":  "ubuntu:latest",
			"status": "created",
		},
		{
			"id":     "elr-0987654321",
			"name":   "python-app",
			"image":  "python:3.9",
			"status": "running",
		},
	})
}

func (n *NetworkManager) listModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{
			"id":      "elr-chat",
			"name":    "ELR Chat Model",
			"version": "1.0",
			"type":    "text",
		},
		{
			"id":      "fish-speech",
			"name":    "Fish Speech Model",
			"version": "1.0",
			"type":    "speech",
		},
	})
}

func main() {
	port := 8080
	
	if len(os.Args) > 1 {
		fmt.Sscanf(os.Args[1], "%d", &port)
	}
	
	nm := NewNetworkManager(port)
	
	if err := nm.Start(); err != nil {
		fmt.Printf("Failed to start network service: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("ELR Network Service started successfully!")
	fmt.Println("Press Ctrl+C to stop")
	
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	
	fmt.Println("\nShutting down ELR Network Service...")
	if err := nm.Stop(); err != nil {
		fmt.Printf("Failed to stop network service: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("ELR Network Service stopped successfully!")
}