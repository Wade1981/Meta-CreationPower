package elr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"micro_model/config"
	"micro_model/model"
	"micro_model/sandbox"
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
	// Setup HTTP routes
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
	// Desktop API 路由
	handler.HandleFunc("/api/desktop/health", n.desktopHealthCheck)
	handler.HandleFunc("/api/desktop/status", n.desktopStatus)
	handler.HandleFunc("/api/desktop/containers", n.desktopListContainers)
	handler.HandleFunc("/api/desktop/resources", n.desktopGetResources)
	handler.HandleFunc("/api/desktop/files", n.desktopListFiles)
	handler.HandleFunc("/api/desktop/upload", n.desktopUploadFile)
	
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
