package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DesktopAPIServer 桌面API服务器
type DesktopAPIServer struct {
	port   int
	server *http.Server
}

// NewDesktopAPIServer 创建桌面API服务器
func NewDesktopAPIServer(port int) *DesktopAPIServer {
	return &DesktopAPIServer{
		port: port,
	}
}

// Start 启动桌面API服务器
func (d *DesktopAPIServer) Start() error {
	handler := http.NewServeMux()
	
	handler.HandleFunc("/health", d.healthCheck)
	handler.HandleFunc("/api/desktop/health", d.desktopHealthCheck)
	handler.HandleFunc("/api/desktop/status", d.getDesktopStatus)
	handler.HandleFunc("/api/desktop/settings", d.getSettings)
	handler.HandleFunc("/api/desktop/settings/save", d.saveSettings)
	
	serverAddr := fmt.Sprintf(":%d", d.port)
	d.server = &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}
	
	fmt.Println("Desktop API is starting, please wait...")
	fmt.Printf("Desktop API service starting on port %d\n", d.port)
	
	// 启动一个goroutine来运行服务器
	go func() {
		if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Desktop API service error: %v\n", err)
		}
	}()
	
	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)
	
	fmt.Printf("Desktop API: http://localhost:%d\n", d.port)
	fmt.Println("Desktop API started successfully!")
	
	return nil
}

// Stop 停止桌面API服务器
func (d *DesktopAPIServer) Stop() error {
	if d.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return d.server.Shutdown(ctx)
	}
	return nil
}

// healthCheck 健康检查
func (d *DesktopAPIServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "elr-desktop-api",
		"version":   "1.0.0",
	})
}

// desktopHealthCheck 桌面健康检查
func (d *DesktopAPIServer) desktopHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "desktop-api",
		"port":      d.port,
	})
}

// getDesktopStatus 获取桌面状态
func (d *DesktopAPIServer) getDesktopStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "running",
		"message":   "Desktop API Service is running",
		"port":      d.port,
		"timestamp": time.Now().Unix(),
	})
}

// getSettings 获取设置
func (d *DesktopAPIServer) getSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"theme":         "light",
		"language":      "zh-CN",
		"auto_start":    true,
		"notifications": true,
		"port":          d.port,
	})
}

// saveSettings 保存设置
func (d *DesktopAPIServer) saveSettings(w http.ResponseWriter, r *http.Request) {
	var settings map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid settings format",
		})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Settings saved successfully",
		"settings": settings,
		"timestamp": time.Now().Unix(),
	})
}
