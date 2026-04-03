package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MockModelAPIConfig 模拟Model API配置
type MockModelAPIConfig struct {
	Host string
	Port int
}

// GetHost 获取主机
func (c *MockModelAPIConfig) GetHost() string {
	return c.Host
}

// GetPort 获取端口
func (c *MockModelAPIConfig) GetPort() int {
	return c.Port
}

// MockModelAPIServer 模拟Model API服务器
type MockModelAPIServer struct {
	config *MockModelAPIConfig
	server *http.Server
}

// NewMockModelAPIServer 创建模型API服务器（模拟实现）
func NewMockModelAPIServer(config *MockModelAPIConfig) *MockModelAPIServer {
	return &MockModelAPIServer{
		config: config,
	}
}

// Start 启动模型API服务器
func (m *MockModelAPIServer) Start() error {
	// 创建HTTP服务器
	handler := http.NewServeMux()
	
	// 设置API路由
	handler.HandleFunc("/health", m.healthCheck)
	handler.HandleFunc("/api/models", m.listModels)
	handler.HandleFunc("/api/models/elr-chat", m.getModel)
	handler.HandleFunc("/api/containers", m.listContainers)
	handler.HandleFunc("/api/sandbox/status", m.getRuntimeStatus)
	handler.HandleFunc("/api/monitor/metrics", m.getMetrics)
	handler.HandleFunc("/api/monitor/status", m.getMonitorStatus)

	serverAddr := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	m.server = &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}

	// 启动服务器
	fmt.Println("Model API is starting, please wait...")
	fmt.Printf("Model API service starting on port %d\n", m.config.Port)
	
	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Model API service error: %v\n", err)
		}
	}()
	
	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)
	
	fmt.Printf("Model API: http://localhost:%d\n", m.config.Port)
	fmt.Println("Model API started successfully!")
	
	return nil
}

// Stop 停止模型API服务器
func (m *MockModelAPIServer) Stop() error {
	if m.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return m.server.Shutdown(ctx)
	}
	return nil
}

// healthCheck 健康检查
func (m *MockModelAPIServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "elr-model-api",
	})
}

// listModels 列出所有模型
func (m *MockModelAPIServer) listModels(w http.ResponseWriter, r *http.Request) {
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

// getModel 获取模型信息
func (m *MockModelAPIServer) getModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      "elr-chat",
		"name":    "ELR Chat Model",
		"version": "1.0",
		"type":    "text",
		"description": "A chat model for ELR",
	})
}

// listContainers 列出所有容器
func (m *MockModelAPIServer) listContainers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{
			"name":   "test-container",
			"status": "running",
			"model":  "elr-chat",
		},
	})
}

// getRuntimeStatus 获取运行时状态
func (m *MockModelAPIServer) getRuntimeStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "running",
		"container": "test-container",
		"model":     "elr-chat",
	})
}

// getMetrics 获取监控指标
func (m *MockModelAPIServer) getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("# HELP model_runs_total Total number of model runs\n# TYPE model_runs_total counter\nmodel_runs_total 0\n"))
}

// getMonitorStatus 获取监控状态
func (m *MockModelAPIServer) getMonitorStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "monitoring",
	})
}
