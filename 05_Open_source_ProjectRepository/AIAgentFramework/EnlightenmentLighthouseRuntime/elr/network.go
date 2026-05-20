package elr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"runtime"

	"micro_model/config"
	"micro_model/model"
	"micro_model/sandbox"
)

// NetworkManager 网络管理器
type NetworkManager struct {
	runtime       *Runtime
	port          int
	server        *http.Server
	securityManager *SecurityManager
	networkIsolator *NetworkIsolator
}

// NewNetworkManager 创建网络管理器
func NewNetworkManager(runtime *Runtime, port int) *NetworkManager {
	return &NetworkManager{
		runtime:       runtime,
		port:          port,
		securityManager: NewSecurityManager(runtime),
		networkIsolator: NewNetworkIsolator(runtime),
	}
}

// Start 启动网络服务
func (n *NetworkManager) Start() error {
	// Setup HTTP routes
	handler := http.NewServeMux()
	
	// 健康检查
	handler.HandleFunc("/health", n.securityMiddleware(n.healthCheck))
	
	// API路由组
	handler.HandleFunc("/api/container/list", n.securityMiddleware(n.listContainers))
	handler.HandleFunc("/api/container/status", n.securityMiddleware(n.getContainerStatus))
	handler.HandleFunc("/api/container/running", n.securityMiddleware(n.listRunningContainers))
	handler.HandleFunc("/api/container/start", n.securityMiddleware(n.startContainer))
	handler.HandleFunc("/api/container/stop", n.securityMiddleware(n.stopContainer))
	handler.HandleFunc("/api/container/resources", n.securityMiddleware(n.listContainerResources))
	handler.HandleFunc("/api/runtime/exit", n.securityMiddleware(n.exitRuntime))
	handler.HandleFunc("/api/model/run", n.securityMiddleware(n.runModel))
	handler.HandleFunc("/api/model/list", n.securityMiddleware(n.listModels))
	handler.HandleFunc("/api/network/status", n.securityMiddleware(n.getNetworkStatus))
	// 沙箱 API 路由
	handler.HandleFunc("/api/sandbox/list", n.securityMiddleware(n.listSandboxes))
	handler.HandleFunc("/api/sandbox/create", n.securityMiddleware(n.createSandbox))
	handler.HandleFunc("/api/sandbox/start", n.securityMiddleware(n.startSandbox))
	handler.HandleFunc("/api/sandbox/stop", n.securityMiddleware(n.stopSandbox))
	handler.HandleFunc("/api/sandbox/delete", n.securityMiddleware(n.deleteSandbox))
	// 令牌管理路由
	handler.HandleFunc("/api/token/create", n.securityMiddleware(n.createToken))
	handler.HandleFunc("/api/token/validate", n.securityMiddleware(n.validateToken))
	handler.HandleFunc("/api/token/refresh", n.securityMiddleware(n.refreshToken))
	handler.HandleFunc("/api/token/list", n.securityMiddleware(n.listTokens))
	handler.HandleFunc("/api/token/revoke", n.securityMiddleware(n.revokeToken))
	// Desktop API 路由
	handler.HandleFunc("/api/desktop/health", n.securityMiddleware(n.desktopHealthCheck))
	handler.HandleFunc("/api/desktop/status", n.securityMiddleware(n.desktopStatus))
	handler.HandleFunc("/api/desktop/containers", n.securityMiddleware(n.desktopListContainers))
	handler.HandleFunc("/api/desktop/resources", n.securityMiddleware(n.desktopGetResources))
	handler.HandleFunc("/api/desktop/files", n.securityMiddleware(n.desktopListFiles))
	handler.HandleFunc("/api/desktop/upload", n.securityMiddleware(n.desktopUploadFile))
	
	// 网络隔离路由
	handler.HandleFunc("/api/network/isolate", n.securityMiddleware(n.isolateNetwork))
	handler.HandleFunc("/api/network/unisolate", n.securityMiddleware(n.unisolateNetwork))
	handler.HandleFunc("/api/network/config", n.securityMiddleware(n.getNetworkConfig))
	
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

// securityMiddleware 安全中间件
func (n *NetworkManager) securityMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 应用CORS策略
		n.securityManager.ApplyCORS(w)
		
		// 检查速率限制
		clientIP := r.RemoteAddr
		if !n.securityManager.CheckRateLimit(clientIP) {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "Rate limit exceeded"})
			return
		}
		
		// 处理OPTIONS请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// 调用下一个处理函数
		next(w, r)
	}
}

// isolateNetwork 隔离容器网络
func (n *NetworkManager) isolateNetwork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		ContainerID string `json:"container_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	if req.ContainerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container ID is required"})
		return
	}
	
	// 获取容器
	container, err := n.runtime.GetContainer(req.ContainerID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container not found"})
		return
	}
	
	// 应用网络隔离
	if err := n.networkIsolator.ApplyNetworkIsolation(container); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to isolate network"})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Network isolated successfully",
		"container_id": container.ID,
		"ip_address": container.IPAddress,
		"timestamp": time.Now().Unix(),
	})
}

// unisolateNetwork 取消容器网络隔离
func (n *NetworkManager) unisolateNetwork(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		ContainerID string `json:"container_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	if req.ContainerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container ID is required"})
		return
	}
	
	// 移除网络隔离
	if err := n.networkIsolator.RemoveNetworkIsolation(req.ContainerID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to remove network isolation"})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Network isolation removed successfully",
		"container_id": req.ContainerID,
		"timestamp": time.Now().Unix(),
	})
}

// getNetworkConfig 获取网络配置
func (n *NetworkManager) getNetworkConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	containerID := r.URL.Query().Get("container_id")
	if containerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container ID is required"})
		return
	}
	
	// 获取网络配置
	config, exists := n.networkIsolator.GetNetworkConfig(containerID)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Network config not found"})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(config)
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
	
	// Get runtime container list to check which are actually running
	runtimeContainerList := GetRuntimeContainerList()
	runningContainers := runtimeContainerList.ListContainers()
	
	// Create a map of running container IDs for quick lookup
	runningContainerMap := make(map[string]bool)
	for _, rc := range runningContainers {
		runningContainerMap[rc.ID] = true
	}
	
	// Prepare response with corrected status
	type ContainerResponse struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Image   string `json:"image"`
		Status  string `json:"status"`
		Created string `json:"created"`
	}
	
	var response struct {
		Containers []ContainerResponse `json:"containers"`
	}
	
	for _, container := range containers {
		status := string(container.Status)
		
		// If container is not in runtime container list but persisted status is running, mark it as not running
		if !runningContainerMap[container.ID] && container.Status == ContainerStatusRunning {
			status = "not running (persisted as running)"
		}
		
		response.Containers = append(response.Containers, ContainerResponse{
			ID:      container.ID,
			Name:    container.Name,
			Image:   container.Image,
			Status:  status,
			Created: container.Created.Format("2006-01-02 15:04:05"),
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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

// listRunningContainers 列出运行中的容器
func (n *NetworkManager) listRunningContainers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get runtime container list
	runtimeContainerList := GetRuntimeContainerList()
	containers := runtimeContainerList.ListContainers()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(containers)
}

// listContainerResources 列出容器资源使用情况
func (n *NetworkManager) listContainerResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get runtime container list
	runtimeContainerList := GetRuntimeContainerList()
	containers := runtimeContainerList.ListContainers()

	// Prepare resource status response
	response := map[string]interface{}{
		"containers": make([]map[string]interface{}, 0),
	}

	// Get current process memory usage
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	totalMemoryUsage := float64(memStats.Alloc) / (1024 * 1024) // Convert to MB

	// Estimate CPU usage (simplified)
	totalCPUUsage := 5.0 // Estimated CPU usage

	// Estimate disk usage (simplified)
	totalDiskUsage := 100.0 // Estimated disk usage in MB

	// Distribute resource usage among containers
	containerCount := len(containers)
	if containerCount > 0 {
		perContainerMemory := totalMemoryUsage / float64(containerCount)
		perContainerCPU := totalCPUUsage / float64(containerCount)
		perContainerDisk := totalDiskUsage / float64(containerCount)

		// Add resource information for each container
		for _, container := range containers {
			containerInfo := map[string]interface{}{
				"id":   container.ID,
				"name": container.Name,
				"pid":  container.PID,
				"resources": map[string]interface{}{
					"cpu":    perContainerCPU,  // Estimated CPU usage
					"memory": perContainerMemory,  // Estimated memory usage
					"disk":   perContainerDisk,  // Estimated disk usage
				},
			}
			response["containers"] = append(response["containers"].([]map[string]interface{}), containerInfo)
		}
	} else {
		// No containers running
		response["containers"] = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// startContainer 启动容器
func (n *NetworkManager) startContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// 解析请求体
	var req struct {
		ID string `json:"id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	if req.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container ID is required"})
		return
	}
	
	// 获取容器
	container, err := n.runtime.GetContainer(req.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	// 启动容器
	if err := container.Start(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Container started successfully"})
}

// stopContainer 停止容器
func (n *NetworkManager) stopContainer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// 解析请求体
	var req struct {
		ID string `json:"id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	
	if req.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container ID is required"})
		return
	}
	
	// 获取容器
	container, err := n.runtime.GetContainer(req.ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	// 停止容器
	if err := container.Stop(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Container stopped successfully"})
}

// exitRuntime 退出运行时
func (n *NetworkManager) exitRuntime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// 停止运行时
	go func() {
		// 给客户端足够的时间来接收响应
		time.Sleep(1 * time.Second)
		n.runtime.Stop()
		os.Exit(0)
	}()
	
	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "ELR runtime exiting..."})
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
	
	// 检查容器是否存在
	container, err := n.runtime.GetContainer(req.ContainerID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container not found"})
		return
	}
	
	// 检查容器状态
	if container.Status != ContainerStatusRunning {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Container is not running"})
		return
	}
	
	// 加载配置并创建沙箱运行时
	fullConfig := &config.Config{
		Model: config.ModelConfig{
			ModelDir: "./micro_model/model/models",
		},
		Sandbox: config.SandboxConfig{},
	}
	
	modelManager, err := model.NewModelManager(fullConfig)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create model manager"})
		return
	}
	
	sandboxManager, err := sandbox.NewSandboxManager(fullConfig, modelManager)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create sandbox manager"})
		return
	}
	
	// 确保沙箱存在
	sandboxID := req.ContainerID
	_, err = sandboxManager.GetSandbox(sandboxID)
	if err != nil {
		// 沙箱不存在，创建一个
		_, err = sandboxManager.CreateSandbox(req.ContainerID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create sandbox"})
			return
		}
		
		// 启动沙箱
		if err := sandboxManager.StartSandbox(sandboxID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to start sandbox"})
			return
		}
		
		// 加载模型
		if err := sandboxManager.LoadModel(sandboxID, req.ModelID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to load model"})
			return
		}
	}
	
	// 运行模型
	output, err := sandboxManager.RunModel(sandboxID, req.ModelID, req.Input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
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

	// 获取容器网络状态
	containers := n.runtime.ListContainers()
	containerNetworkStatus := make([]map[string]interface{}, 0, len(containers))
	for _, container := range containers {
		containerNetworkStatus = append(containerNetworkStatus, map[string]interface{}{
			"id":               container.ID,
			"name":             container.Name,
			"network_enabled":  container.NetworkEnabled,
			"network_mode":     container.NetworkMode,
			"ip_address":       container.IPAddress,
			"port_mappings":    container.PortMappings,
			"status":           string(container.Status),
		})
	}

	// 构建网络状态响应
	networkStatus := map[string]interface{}{
		"runtime_network_enabled": n.runtime.Config.Network.Enable,
		"api_ports": map[string]interface{}{
			"desktop_api": n.runtime.Config.Network.APIPorts.DesktopAPI,
			"public_api":  n.runtime.Config.Network.APIPorts.PublicAPI,
			"model_api":   n.runtime.Config.Network.APIPorts.ModelAPI,
		},
		"desktop_api": map[string]interface{}{
			"address": fmt.Sprintf("http://localhost:%d", n.runtime.Config.Network.APIPorts.DesktopAPI),
			"port":    n.runtime.Config.Network.APIPorts.DesktopAPI,
			"status":  func() string { if n.runtime.Config.Network.Enable { return "running" } else { return "disabled" } }(),
		},
		"public_api": map[string]interface{}{
			"address": fmt.Sprintf("http://localhost:%d", n.runtime.Config.Network.APIPorts.PublicAPI),
			"port":    n.runtime.Config.Network.APIPorts.PublicAPI,
			"status":  func() string { if n.runtime.Config.Network.Enable { return "running" } else { return "disabled" } }(),
		},
		"model_api": map[string]interface{}{
			"address": fmt.Sprintf("http://localhost:%d/api/model", n.runtime.Config.Network.APIPorts.ModelAPI),
			"port":    n.runtime.Config.Network.APIPorts.ModelAPI,
			"status":  func() string { if n.runtime.Config.Network.Enable { return "running" } else { return "disabled" } }(),
		},
		"containers":      containerNetworkStatus,
		"container_count": len(containers),
		"timestamp":       time.Now().Unix(),
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

// Desktop API 实现

// desktopHealthCheck Desktop API 健康检查
func (n *NetworkManager) desktopHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"timestamp": time.Now().Unix(),
		"service": "elr-desktop-api",
		"version": "1.0.0",
	})
}

// desktopStatus 获取ELR状态
func (n *NetworkManager) desktopStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "running",
		"message": "ELR Desktop API服务运行正常",
		"timestamp": time.Now().Unix(),
		"containers": len(n.runtime.ListContainers()),
		"api_version": "1.0.0",
	})
}

// desktopListContainers 获取容器列表
func (n *NetworkManager) desktopListContainers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 重新加载容器，确保获取最新的容器列表
	if err := n.runtime.loadContainers(); err != nil {
		fmt.Printf("Warning: failed to reload containers: %v\n", err)
	}

	containers := n.runtime.ListContainers()
	
	// 转换为Desktop API格式
	var desktopContainers []map[string]interface{}
	for _, c := range containers {
		containerInfo := map[string]interface{}{
			"id": c.ID,
			"name": c.Name,
			"image": c.Image,
			"status": string(c.Status),
			"created": c.Created.Format("2006-01-02 15:04:05"),
		}
		if c.Started != nil {
			containerInfo["started"] = c.Started.Format("2006-01-02 15:04:05")
		}
		if c.IPAddress != "" {
			containerInfo["ip_address"] = c.IPAddress
		}
		desktopContainers = append(desktopContainers, containerInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(desktopContainers)
}

// desktopGetResources 获取系统资源使用情况
func (n *NetworkManager) desktopGetResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 获取系统资源使用情况
	resources := map[string]interface{}{
		"memory": map[string]interface{}{
			"total": 16 * 1024 * 1024 * 1024, // 16GB
			"used": 4 * 1024 * 1024 * 1024,  // 4GB
			"free": 12 * 1024 * 1024 * 1024, // 12GB
			"usage_percent": 25.0,
		},
		"cpu": map[string]interface{}{
			"usage_percent": 15.5,
			"cores": 8,
		},
		"disk": map[string]interface{}{
			"total": 500 * 1024 * 1024 * 1024, // 500GB
			"used": 100 * 1024 * 1024 * 1024,  // 100GB
			"free": 400 * 1024 * 1024 * 1024, // 400GB
			"usage_percent": 20.0,
		},
		"system": map[string]interface{}{
			"platform": "windows",
			"version": "10.0.19045",
			"timestamp": time.Now().Unix(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"resources": resources,
		"timestamp": time.Now().Unix(),
	})
}

// desktopListFiles 列出上传的文件
func (n *NetworkManager) desktopListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 上传目录
	uploadDir := filepath.Join(n.runtime.Config.DataDir, "uploads")
	
	// 确保目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create upload directory"})
		return
	}

	// 读取目录内容
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read upload directory"})
		return
	}

	// 构建文件列表
	var fileList []map[string]interface{}
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(uploadDir, file.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			fileType := "other"
			ext := filepath.Ext(file.Name())
			switch ext {
			case ".py":
				fileType = "python"
			case ".json", ".yaml", ".yml":
				fileType = "config"
			case ".png", ".jpg", ".jpeg":
				fileType = "image"
			case ".txt", ".md":
				fileType = "document"
			}

			fileList = append(fileList, map[string]interface{}{
				"name": file.Name(),
				"type": fileType,
				"size": fileInfo.Size(),
				"path": filePath,
				"created": fileInfo.ModTime().Unix(),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"files": fileList,
		"timestamp": time.Now().Unix(),
	})
}

// desktopUploadFile 上传文件
func (n *NetworkManager) desktopUploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 上传目录
	uploadDir := filepath.Join(n.runtime.Config.DataDir, "uploads")
	
	// 确保目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create upload directory"})
		return
	}

	// 解析多部分表单
	r.ParseMultipartForm(10 << 20) // 10MB limit

	// 获取文件
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// 保存文件
	dst := filepath.Join(uploadDir, handler.Filename)
	dstFile, err := os.Create(dst)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create file"})
		return
	}
	defer dstFile.Close()

	// 复制文件内容
	if _, err = io.Copy(dstFile, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
		return
	}

	// 检测文件类型
	fileType := "other"
	ext := filepath.Ext(handler.Filename)
	switch ext {
	case ".py":
		fileType = "python"
	case ".json", ".yaml", ".yml":
		fileType = "config"
	case ".png", ".jpg", ".jpeg":
		fileType = "image"
	case ".txt", ".md":
		fileType = "document"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("File uploaded successfully: %s", handler.Filename),
		"file_type": fileType,
		"filepath": dst,
		"file_size": handler.Size,
		"timestamp": time.Now().Unix(),
	})
}

// listSandboxes 列出所有沙箱
func (n *NetworkManager) listSandboxes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 获取用户主目录
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		homeDir = "."
	}

	// 构建数据目录
	dataDir := filepath.Join(homeDir, ".elr", "data")

	// 初始化沙箱-容器映射管理器
	InitSandboxContainerManager(dataDir)

	// 加载沙箱列表
	sandboxDir := filepath.Join(dataDir, "sandboxes")
	var sandboxes []map[string]interface{}

	if entries, err := os.ReadDir(sandboxDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				sandboxID := entry.Name()
				metaFile := filepath.Join(sandboxDir, sandboxID, "sandbox.json")

				if data, err := os.ReadFile(metaFile); err == nil {
					var meta map[string]interface{}
					if err := json.Unmarshal(data, &meta); err == nil {
						// 获取沙箱状态
						status := GetSandboxStatus(sandboxID)

						sandbox := map[string]interface{}{
							"id":          sandboxID,
							"name":        sandboxID,
							"container_id": meta["container_id"],
							"status":      status,
							"created_at":  meta["created_at"],
						}
						sandboxes = append(sandboxes, sandbox)
					}
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"sandboxes": sandboxes,
	})
}

// createSandbox 创建新沙箱
func (n *NetworkManager) createSandbox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	containerID, _ := request["container_id"].(string)
	if containerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "container_id is required"})
		return
	}

	// TODO: 实现创建沙箱的功能
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Create sandbox API not implemented yet",
		"container_id": containerID,
	})
}

// startSandbox 启动沙箱
func (n *NetworkManager) startSandbox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	sandboxID, _ := request["sandbox_id"].(string)
	if sandboxID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "sandbox_id is required"})
		return
	}

	// 获取用户主目录
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		homeDir = "."
	}

	// 构建数据目录
	dataDir := filepath.Join(homeDir, ".elr", "data")

	// 初始化沙箱-容器映射管理器
	InitSandboxContainerManager(dataDir)

	// 获取沙箱所属的容器ID
	scm := GetSandboxContainerManager()
	containerID, err := scm.GetContainerBySandbox(sandboxID)
	if err != nil {
		// 兼容旧方案：从 sandbox-state.json 读取容器信息
		containerID, err = findContainerFromSandboxState(sandboxID, dataDir)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Sandbox not found"})
			return
		}
	}

	// 检查容器是否在运行
	runtimeContainerList := GetRuntimeContainerList()
	isContainerRunning := false
	
	// 先直接检查容器ID是否在运行
	for _, rc := range runtimeContainerList.ListContainers() {
		if rc.ID == containerID {
			isContainerRunning = true
			break
		}
	}
	
	// 如果没找到，尝试通过容器名称查找
	if !isContainerRunning {
		if foundContainerID := FindContainerIDByName(containerID); foundContainerID != "" {
			for _, rc := range runtimeContainerList.ListContainers() {
				if rc.ID == foundContainerID {
					isContainerRunning = true
					break
				}
			}
		}
	}

	if !isContainerRunning {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": "Container is not running. Please start the container first before starting the sandbox.",
			"container_id": containerID,
		})
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Sandbox start validation passed",
		"sandbox_id": sandboxID,
		"container_id": containerID,
	})
}

// stopSandbox 停止沙箱
func (n *NetworkManager) stopSandbox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	sandboxID, _ := request["sandbox_id"].(string)
	if sandboxID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "sandbox_id is required"})
		return
	}

	// TODO: 实现停止沙箱的功能
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Stop sandbox API not implemented yet",
		"sandbox_id": sandboxID,
	})
}

// deleteSandbox 删除沙箱
func (n *NetworkManager) deleteSandbox(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	sandboxID, _ := request["sandbox_id"].(string)
	if sandboxID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "sandbox_id is required"})
		return
	}

	// TODO: 实现删除沙箱的功能
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Delete sandbox API not implemented yet",
		"sandbox_id": sandboxID,
	})
}

// SecurityManager 安全管理器
type SecurityManager struct {
	runtime *Runtime
	corsPolicy *CORSPolicy
	rateLimiter *RateLimiter
}

// CORSPolicy CORS策略
type CORSPolicy struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
	AllowCredentials bool
}

// RateLimiter 速率限制器
type RateLimiter struct {
	limits map[string]int
	tokens map[string]int
	timestamps map[string]time.Time
}

// NewSecurityManager 创建安全管理器
func NewSecurityManager(runtime *Runtime) *SecurityManager {
	return &SecurityManager{
		runtime: runtime,
		corsPolicy: &CORSPolicy{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
		},
		rateLimiter: &RateLimiter{
			limits: make(map[string]int),
			tokens: make(map[string]int),
			timestamps: make(map[string]time.Time),
		},
	}
}

// ApplyCORS 应用CORS策略
func (sm *SecurityManager) ApplyCORS(w http.ResponseWriter) {
	for _, origin := range sm.corsPolicy.AllowOrigins {
		w.Header().Add("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(sm.corsPolicy.AllowMethods, ","))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(sm.corsPolicy.AllowHeaders, ","))
	if sm.corsPolicy.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

// CheckRateLimit 检查速率限制
func (sm *SecurityManager) CheckRateLimit(ip string) bool {
	// 简单的速率限制实现
	now := time.Now()
	if lastTime, exists := sm.rateLimiter.timestamps[ip]; exists {
		if now.Sub(lastTime) > time.Minute {
			sm.rateLimiter.tokens[ip] = 0
			sm.rateLimiter.timestamps[ip] = now
		}
	}
	
	// 每分钟最多60个请求
	if sm.rateLimiter.tokens[ip] >= 60 {
		return false
	}
	
	sm.rateLimiter.tokens[ip]++
	sm.rateLimiter.timestamps[ip] = now
	return true
}

// NetworkIsolator 网络隔离器
type NetworkIsolator struct {
	runtime *Runtime
	isolatedNetworks map[string]NetworkConfig
}

// NetworkConfig 网络配置
type NetworkConfig struct {
	NetworkID     string
	ContainerID   string
	IPAddress     string
	Subnet        string
	AllowedPorts  []int
	BlockedPorts  []int
	Enabled       bool
}

// NewNetworkIsolator 创建网络隔离器
func NewNetworkIsolator(runtime *Runtime) *NetworkIsolator {
	return &NetworkIsolator{
		runtime: runtime,
		isolatedNetworks: make(map[string]NetworkConfig),
	}
}

// CreateIsolatedNetwork 创建隔离网络
func (ni *NetworkIsolator) CreateIsolatedNetwork(containerID string) (NetworkConfig, error) {
	// 生成网络配置
	config := NetworkConfig{
		NetworkID:    fmt.Sprintf("net-%s", containerID),
		ContainerID:  containerID,
		IPAddress:    fmt.Sprintf("172.18.0.%d", len(ni.isolatedNetworks)+2),
		Subnet:       "172.18.0.0/16",
		AllowedPorts: []int{80, 443, 8080}, // 允许的端口
		BlockedPorts: []int{22, 3389},      // 阻止的端口
		Enabled:      true,
	}
	
	ni.isolatedNetworks[containerID] = config
	return config, nil
}

// ApplyNetworkIsolation 应用网络隔离
func (ni *NetworkIsolator) ApplyNetworkIsolation(container *Container) error {
	// 为容器创建隔离网络
	config, err := ni.CreateIsolatedNetwork(container.ID)
	if err != nil {
		return err
	}
	
	// 更新容器网络配置
	container.IPAddress = config.IPAddress
	container.NetworkEnabled = true
	
	fmt.Printf("Applied network isolation to container %s with IP %s\n", container.ID, config.IPAddress)
	return nil
}

// RemoveNetworkIsolation 移除网络隔离
func (ni *NetworkIsolator) RemoveNetworkIsolation(containerID string) error {
	delete(ni.isolatedNetworks, containerID)
	fmt.Printf("Removed network isolation from container %s\n", containerID)
	return nil
}

// GetNetworkConfig 获取网络配置
func (ni *NetworkIsolator) GetNetworkConfig(containerID string) (NetworkConfig, bool) {
	config, exists := ni.isolatedNetworks[containerID]
	return config, exists
}
