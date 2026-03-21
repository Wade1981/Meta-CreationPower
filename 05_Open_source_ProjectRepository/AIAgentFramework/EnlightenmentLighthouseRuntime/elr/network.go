package elr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NetworkManager 网络管理器
type NetworkManager struct {
	runtime *Runtime
	port    int
	server  *http.Server
}

// NewNetworkManager 创建网络管理器
func NewNetworkManager(runtime *Runtime, port int) *NetworkManager {
	return &NetworkManager{
		runtime: runtime,
		port:    port,
	}
}

// Start 启动网络服务
func (n *NetworkManager) Start() error {
	// 设置HTTP路由
	handler := http.NewServeMux()
	
	// 健康检查
	handler.HandleFunc("/health", n.healthCheck)
	
	// API路由组
	handler.HandleFunc("/api/container/list", n.listContainers)
	handler.HandleFunc("/api/container/status", n.getContainerStatus)
	handler.HandleFunc("/api/model/run", n.runModel)
	handler.HandleFunc("/api/model/list", n.listModels)
	handler.HandleFunc("/api/network/status", n.getNetworkStatus)
	// 令牌管理路由
	handler.HandleFunc("/api/token/create", n.createToken)
	handler.HandleFunc("/api/token/validate", n.validateToken)
	handler.HandleFunc("/api/token/refresh", n.refreshToken)
	handler.HandleFunc("/api/token/list", n.listTokens)
	handler.HandleFunc("/api/token/revoke", n.revokeToken)
	
	// 创建HTTP服务器
	serverAddr := fmt.Sprintf(":%d", n.port)
	n.server = &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}
	
	// 启动服务器
	fmt.Printf("ELR network service starting on port %d\n", n.port)
	go func() {
		if err := n.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Network service error: %v\n", err)
		}
	}()
	
	return nil
}

// Stop 停止网络服务
func (n *NetworkManager) Stop() error {
	if n.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return n.server.Shutdown(ctx)
	}
	return nil
}

// healthCheck 健康检查
func (n *NetworkManager) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "elr-network",
		"version":   "1.0.0",
	})
}

// listContainers 列出所有容器
func (n *NetworkManager) listContainers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	containers := n.runtime.ListContainers()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(containers)
}

// getContainerStatus 获取容器状态
func (n *NetworkManager) getContainerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	containerID := r.URL.Query().Get("id")
	if containerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container ID is required"})
		return
	}
	
	container, err := n.runtime.GetContainer(containerID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container not found"})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(container)
}

// runModel 运行模型
func (n *NetworkManager) runModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		ContainerID string `json:"container_id"`
		ModelID     string `json:"model_id"`
		Input       string `json:"input"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	if req.ContainerID == "" || req.ModelID == "" || req.Input == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container ID, model ID, and input are required"})
		return
	}
	
	// 运行模型（模拟实现）
	output := fmt.Sprintf("Model %s in container %s processed input: %s", req.ModelID, req.ContainerID, req.Input)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"output": output,
		"timestamp": time.Now().Unix(),
	})
}

// listModels 列出所有模型
func (n *NetworkManager) listModels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// 模拟模型列表
	models := []map[string]interface{}{
		{
			"id": "elr-chat",
			"name": "ELR Chat Model",
			"version": "1.0",
			"type": "text",
		},
		{
			"id": "fish-speech",
			"name": "Fish Speech Model",
			"version": "1.0",
			"type": "speech",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models)
}

// getNetworkStatus 获取网络状态
func (n *NetworkManager) getNetworkStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 构建网络状态响应
	networkStatus := map[string]interface{}{
		"desktop_api": map[string]interface{}{
			"address": "http://localhost:8081",
			"port":    8081,
			"status":  "running",
		},
		"public_api": map[string]interface{}{
			"address": "http://localhost:8080",
			"port":    8080,
			"status":  "running",
		},
		"model_api": map[string]interface{}{
			"address": "http://localhost:8080/api/model",
			"port":    8080,
			"status":  "running",
		},
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(networkStatus)
}

// createToken 创建新令牌
func (n *NetworkManager) createToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Description == "" {
		req.Description = "ELR Container Token"
	}

	token, err := n.runtime.TokenManager.GenerateToken(req.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"message": "Token generated successfully",
		"timestamp": time.Now().Unix(),
	})
}

// validateToken 验证令牌
func (n *NetworkManager) validateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token is required"})
		return
	}

	valid, message := n.runtime.TokenManager.ValidateToken(req.Token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid": valid,
		"message": message,
		"timestamp": time.Now().Unix(),
	})
}

// refreshToken 刷新令牌
func (n *NetworkManager) refreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Token       string `json:"token"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token is required"})
		return
	}

	if req.Description == "" {
		req.Description = "Refreshed ELR Container Token"
	}

	newToken, err := n.runtime.TokenManager.RefreshToken(req.Token, req.Description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": newToken,
		"message": "Token refreshed successfully",
		"timestamp": time.Now().Unix(),
	})
}

// listTokens 列出所有令牌
func (n *NetworkManager) listTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tokens := n.runtime.TokenManager.ListTokens()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tokens": tokens,
		"count": len(tokens),
		"timestamp": time.Now().Unix(),
	})
}

// revokeToken 撤销令牌
func (n *NetworkManager) revokeToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TokenID string `json:"token_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.TokenID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token ID is required"})
		return
	}

	if err := n.runtime.TokenManager.RevokeToken(req.TokenID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Token revoked successfully",
		"timestamp": time.Now().Unix(),
	})
}
