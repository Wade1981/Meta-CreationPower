package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PublicAPIServer 公共API服务器
type PublicAPIServer struct {
	port   int
	server *http.Server
}

// NewPublicAPIServer 创建公共API服务器
func NewPublicAPIServer(port int) *PublicAPIServer {
	return &PublicAPIServer{
		port: port,
	}
}

// Start 启动公共API服务器
func (p *PublicAPIServer) Start() error {
	handler := http.NewServeMux()
	
	handler.HandleFunc("/health", p.healthCheck)
	handler.HandleFunc("/api/status", p.getStatus)
	handler.HandleFunc("/api/network/status", p.getNetworkStatus)
	handler.HandleFunc("/api/container/list", p.listContainers)
	handler.HandleFunc("/api/model/list", p.listModels)
	
	serverAddr := fmt.Sprintf(":%d", p.port)
	p.server = &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}
	
	fmt.Println("Public API is starting, please wait...")
	fmt.Printf("Public API service starting on port %d\n", p.port)
	
	// 启动一个goroutine来运行服务器
	go func() {
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Public API service error: %v\n", err)
		}
	}()
	
	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)
	
	fmt.Printf("Public API: http://localhost:%d\n", p.port)
	fmt.Println("Public API started successfully!")
	
	return nil
}

// Stop 停止公共API服务器
func (p *PublicAPIServer) Stop() error {
	if p.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return p.server.Shutdown(ctx)
	}
	return nil
}

// healthCheck 健康检查
func (p *PublicAPIServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "elr-public-api",
		"version":   "1.0.0",
	})
}

// getStatus 获取状态
func (p *PublicAPIServer) getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "running",
		"message": "Public API Service is running",
	})
}

// getNetworkStatus 获取网络状态
func (p *PublicAPIServer) getNetworkStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"desktop_api": map[string]interface{}{
			"address": "http://localhost:8081",
			"port":    8081,
			"status":  "check",
		},
		"public_api": map[string]interface{}{
			"address": fmt.Sprintf("http://localhost:%d", p.port),
			"port":    p.port,
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

// listContainers 列出容器
func (p *PublicAPIServer) listContainers(w http.ResponseWriter, r *http.Request) {
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

// listModels 列出模型
func (p *PublicAPIServer) listModels(w http.ResponseWriter, r *http.Request) {
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
