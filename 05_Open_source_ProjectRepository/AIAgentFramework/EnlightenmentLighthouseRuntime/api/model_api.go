package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"micro_model/config"
	"micro_model/container"
	"micro_model/model"
	"micro_model/monitor"
	"micro_model/sandbox"
)

// ModelAPIServer 模型API服务器
type ModelAPIServer struct {
	config          *config.ServerConfig
	modelManager    *model.ModelManager
	containerManager *container.ContainerManager
	sandboxRuntime  *sandbox.SandboxRuntime
	monitorService  *monitor.MonitorService
	server          *http.Server
}

// NewModelAPIServer 创建模型API服务器
func NewModelAPIServer(config *config.ServerConfig, modelManager *model.ModelManager, containerManager *container.ContainerManager, sandboxRuntime *sandbox.SandboxRuntime, monitorService *monitor.MonitorService) *ModelAPIServer {
	return &ModelAPIServer{
		config:          config,
		modelManager:    modelManager,
		containerManager: containerManager,
		sandboxRuntime:  sandboxRuntime,
		monitorService:  monitorService,
	}
}

// Start 启动模型API服务器
func (m *ModelAPIServer) Start() error {
	// 创建HTTP服务器
	handler := http.NewServeMux()
	
	// 设置API路由
	m.setupRoutes(handler)

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
func (m *ModelAPIServer) Stop() error {
	if m.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return m.server.Shutdown(ctx)
	}
	return nil
}

// setupRoutes 设置API路由
func (m *ModelAPIServer) setupRoutes(handler *http.ServeMux) {
	// 健康检查
	handler.HandleFunc("/health", m.healthCheck)

	// API路由
	handler.HandleFunc("/api/models", m.listModels)
	handler.HandleFunc("/api/models/", m.getModel)
	handler.HandleFunc("/api/models/download", m.downloadModel)
	handler.HandleFunc("/api/models/delete", m.deleteModel)
	handler.HandleFunc("/api/models/update", m.updateModel)
	handler.HandleFunc("/api/models/run", m.runModel)

	handler.HandleFunc("/api/containers", m.listContainers)
	handler.HandleFunc("/api/containers/", m.getContainer)
	handler.HandleFunc("/api/containers/create", m.createContainer)
	handler.HandleFunc("/api/containers/start", m.startContainer)
	handler.HandleFunc("/api/containers/stop", m.stopContainer)
	handler.HandleFunc("/api/containers/remove", m.removeContainer)

	handler.HandleFunc("/api/sandbox/status", m.getRuntimeStatus)
	handler.HandleFunc("/api/sandbox/execute", m.executeCommand)
	handler.HandleFunc("/api/sandbox/stop", m.stopRuntime)

	handler.HandleFunc("/api/monitor/metrics", m.getMetrics)
	handler.HandleFunc("/api/monitor/status", m.getMonitorStatus)
}

// healthCheck 健康检查
func (m *ModelAPIServer) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "elr-model-api",
	})
}

// listModels 列出所有模型
func (m *ModelAPIServer) listModels(w http.ResponseWriter, r *http.Request) {
	models, err := m.modelManager.ListModels()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models)
}

// getModel 获取模型信息
func (m *ModelAPIServer) getModel(w http.ResponseWriter, r *http.Request) {
	modelID := r.URL.Path[len("/api/models/"):]
	model, err := m.modelManager.GetModel(modelID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model)
}

// downloadModel 下载模型
func (m *ModelAPIServer) downloadModel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ModelID     string `json:"model_id"`
		ModelType   string `json:"model_type"`
		DownloadURL string `json:"download_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	if req.ModelID == "" || req.DownloadURL == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "model_id and download_url are required"})
		return
	}

	err := m.modelManager.DownloadModel(req.ModelID, req.ModelType, req.DownloadURL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("Model %s downloaded successfully", req.ModelID)})
}

// deleteModel 删除模型
func (m *ModelAPIServer) deleteModel(w http.ResponseWriter, r *http.Request) {
	modelID := r.URL.Path[len("/api/models/"):]
	err := m.modelManager.DeleteModel(modelID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("Model %s deleted successfully", modelID)})
}

// updateModel 更新模型
func (m *ModelAPIServer) updateModel(w http.ResponseWriter, r *http.Request) {
	modelID := r.URL.Path[len("/api/models/"):]
	var req struct {
		DownloadURL string `json:"download_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	if req.DownloadURL == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "download_url is required"})
		return
	}

	err := m.modelManager.UpdateModel(modelID, req.DownloadURL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("Model %s updated successfully", modelID)})
}

// runModel 运行模型
func (m *ModelAPIServer) runModel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContainerName string `json:"container_name"`
		ModelID       string `json:"model_id"`
		Input         string `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	if req.ContainerName == "" || req.ModelID == "" || req.Input == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "container_name, model_id, and input are required"})
		return
	}

	// 记录模型运行开始时间
	startTime := time.Now()

	// 运行模型
	output, err := m.sandboxRuntime.RunModel(req.ContainerName, req.ModelID, req.Input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	// 记录模型运行结束时间并更新监控指标
	runDuration := time.Since(startTime)
	m.monitorService.RecordModelRun(runDuration)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"output":   output,
		"duration": runDuration.Seconds(),
	})
}

// listContainers 列出所有容器
func (m *ModelAPIServer) listContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := m.containerManager.ListContainers()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(containers)
}

// getContainer 获取容器信息
func (m *ModelAPIServer) getContainer(w http.ResponseWriter, r *http.Request) {
	containerName := r.URL.Path[len("/api/containers/"):]
	container, err := m.containerManager.GetContainer(containerName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(container)
}

// createContainer 创建容器
func (m *ModelAPIServer) createContainer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContainerName string                 `json:"container_name"`
		ModelID       string                 `json:"model_id"`
		Resources     map[string]interface{} `json:"resources"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	if req.ContainerName == "" || req.ModelID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "container_name and model_id are required"})
		return
	}

	err := m.containerManager.CreateContainer(req.ContainerName, req.ModelID, req.Resources)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("Container %s created successfully", req.ContainerName)})
}

// startContainer 启动容器
func (m *ModelAPIServer) startContainer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContainerName string `json:"container_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	if req.ContainerName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "container_name is required"})
		return
	}

	err := m.containerManager.StartContainer(req.ContainerName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("Container %s started successfully", req.ContainerName)})
}

// stopContainer 停止容器
func (m *ModelAPIServer) stopContainer(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContainerName string `json:"container_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	if req.ContainerName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "container_name is required"})
		return
	}

	err := m.containerManager.StopContainer(req.ContainerName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("Container %s stopped successfully", req.ContainerName)})
}

// removeContainer 删除容器
func (m *ModelAPIServer) removeContainer(w http.ResponseWriter, r *http.Request) {
	containerName := r.URL.Path[len("/api/containers/"):]
	err := m.containerManager.RemoveContainer(containerName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("Container %s removed successfully", containerName)})
}

// getRuntimeStatus 获取运行时状态
func (m *ModelAPIServer) getRuntimeStatus(w http.ResponseWriter, r *http.Request) {
	containerName := r.URL.Path[len("/api/sandbox/status/"):]
	status, err := m.sandboxRuntime.GetRuntimeStatus(containerName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// executeCommand 执行命令
func (m *ModelAPIServer) executeCommand(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContainerName string   `json:"container_name"`
		Command       []string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	if req.ContainerName == "" || len(req.Command) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "container_name and command are required"})
		return
	}

	output, err := m.sandboxRuntime.ExecuteCommand(req.ContainerName, req.Command)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"output": output})
}

// stopRuntime 停止运行时
func (m *ModelAPIServer) stopRuntime(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ContainerName string `json:"container_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	if req.ContainerName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "container_name is required"})
		return
	}

	err := m.sandboxRuntime.StopRuntime(req.ContainerName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("Runtime for container %s stopped successfully", req.ContainerName)})
}

// getMetrics 获取监控指标
func (m *ModelAPIServer) getMetrics(w http.ResponseWriter, r *http.Request) {
	// 重定向到Prometheus监控端点
	http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
}

// getMonitorStatus 获取监控状态
func (m *ModelAPIServer) getMonitorStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "monitoring",
	})
}
