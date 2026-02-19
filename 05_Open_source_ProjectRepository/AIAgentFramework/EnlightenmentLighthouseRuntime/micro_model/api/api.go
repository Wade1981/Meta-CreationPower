package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"micro-model/config"
	"micro-model/container"
	"micro-model/model"
	"micro-model/monitor"
	"micro-model/sandbox"
)

// APIServer API服务器
type APIServer struct {
	config         *config.ServerConfig
	modelManager   *model.ModelManager
	containerManager *container.ContainerManager
	sandboxRuntime *sandbox.SandboxRuntime
	monitorService *monitor.MonitorService
	server         *http.Server
}

// NewAPIServer 创建API服务器
func NewAPIServer(config *config.ServerConfig, modelManager *model.ModelManager, containerManager *container.ContainerManager, sandboxRuntime *sandbox.SandboxRuntime, monitorService *monitor.MonitorService) *APIServer {
	return &APIServer{
		config:         config,
		modelManager:   modelManager,
		containerManager: containerManager,
		sandboxRuntime: sandboxRuntime,
		monitorService: monitorService,
	}
}

// Start 启动API服务器
func (a *APIServer) Start() error {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	router := gin.Default()

	// 设置API路由
	a.setupRoutes(router)

	// 创建HTTP服务器
	a.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.config.Host, a.config.Port),
		Handler: router,
	}

	// 启动服务器
	fmt.Printf("API server started on %s:%d\n", a.config.Host, a.config.Port)
	return a.server.ListenAndServe()
}

// Stop 停止API服务器
func (a *APIServer) Stop() error {
	if a.server != nil {
		return a.server.Shutdown(nil)
	}
	return nil
}

// setupRoutes 设置API路由
func (a *APIServer) setupRoutes(router *gin.Engine) {
	// 健康检查
	router.GET("/health", a.healthCheck)

	// API路由组
	api := router.Group("/api")
	{
		// 模型管理
		models := api.Group("/models")
		{
			models.GET("/", a.listModels)
			models.GET("/:id", a.getModel)
			models.POST("/download", a.downloadModel)
			models.DELETE("/:id", a.deleteModel)
			models.PUT("/:id", a.updateModel)
			models.POST("/run", a.runModel)
		}

		// 容器管理
		containers := api.Group("/containers")
		{
			containers.GET("/", a.listContainers)
			containers.GET("/:name", a.getContainer)
			containers.POST("/create", a.createContainer)
			containers.POST("/start", a.startContainer)
			containers.POST("/stop", a.stopContainer)
			containers.DELETE("/:name", a.removeContainer)
		}

		// 沙箱运行时
		sandbox := api.Group("/sandbox")
		{
			sandbox.GET("/status/:container", a.getRuntimeStatus)
			sandbox.POST("/execute", a.executeCommand)
			sandbox.POST("/stop", a.stopRuntime)
		}

		// 监控服务
		monitor := api.Group("/monitor")
		{
			monitor.GET("/metrics", a.getMetrics)
			monitor.GET("/status", a.getMonitorStatus)
		}
	}
}

// healthCheck 健康检查
func (a *APIServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "micro-model-server",
	})
}

// listModels 列出所有模型
func (a *APIServer) listModels(c *gin.Context) {
	models, err := a.modelManager.ListModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models)
}

// getModel 获取模型信息
func (a *APIServer) getModel(c *gin.Context) {
	modelID := c.Param("id")
	model, err := a.modelManager.GetModel(modelID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model)
}

// downloadModel 下载模型
func (a *APIServer) downloadModel(c *gin.Context) {
	var req struct {
		ModelID     string `json:"model_id" binding:"required"`
		ModelType   string `json:"model_type"`
		DownloadURL string `json:"download_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.modelManager.DownloadModel(req.ModelID, req.ModelType, req.DownloadURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Model %s downloaded successfully", req.ModelID)})
}

// deleteModel 删除模型
func (a *APIServer) deleteModel(c *gin.Context) {
	modelID := c.Param("id")
	err := a.modelManager.DeleteModel(modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Model %s deleted successfully", modelID)})
}

// updateModel 更新模型
func (a *APIServer) updateModel(c *gin.Context) {
	modelID := c.Param("id")
	var req struct {
		DownloadURL string `json:"download_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.modelManager.UpdateModel(modelID, req.DownloadURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Model %s updated successfully", modelID)})
}

// runModel 运行模型
func (a *APIServer) runModel(c *gin.Context) {
	var req struct {
		ContainerName string `json:"container_name" binding:"required"`
		ModelID       string `json:"model_id" binding:"required"`
		Input         string `json:"input" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录模型运行开始时间
	startTime := time.Now()

	// 运行模型
	output, err := a.sandboxRuntime.RunModel(req.ContainerName, req.ModelID, req.Input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 记录模型运行结束时间并更新监控指标
	runDuration := time.Since(startTime)
	a.monitorService.RecordModelRun(runDuration)

	c.JSON(http.StatusOK, gin.H{
		"output":   output,
		"duration": runDuration.Seconds(),
	})
}

// listContainers 列出所有容器
func (a *APIServer) listContainers(c *gin.Context) {
	containers, err := a.containerManager.ListContainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, containers)
}

// getContainer 获取容器信息
func (a *APIServer) getContainer(c *gin.Context) {
	containerName := c.Param("name")
	container, err := a.containerManager.GetContainer(containerName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, container)
}

// createContainer 创建容器
func (a *APIServer) createContainer(c *gin.Context) {
	var req struct {
		ContainerName string                 `json:"container_name" binding:"required"`
		ModelID       string                 `json:"model_id" binding:"required"`
		Resources     map[string]interface{} `json:"resources"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.containerManager.CreateContainer(req.ContainerName, req.ModelID, req.Resources)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Container %s created successfully", req.ContainerName)})
}

// startContainer 启动容器
func (a *APIServer) startContainer(c *gin.Context) {
	var req struct {
		ContainerName string `json:"container_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.containerManager.StartContainer(req.ContainerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Container %s started successfully", req.ContainerName)})
}

// stopContainer 停止容器
func (a *APIServer) stopContainer(c *gin.Context) {
	var req struct {
		ContainerName string `json:"container_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.containerManager.StopContainer(req.ContainerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Container %s stopped successfully", req.ContainerName)})
}

// removeContainer 删除容器
func (a *APIServer) removeContainer(c *gin.Context) {
	containerName := c.Param("name")
	err := a.containerManager.RemoveContainer(containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Container %s removed successfully", containerName)})
}

// getRuntimeStatus 获取运行时状态
func (a *APIServer) getRuntimeStatus(c *gin.Context) {
	containerName := c.Param("container")
	status, err := a.sandboxRuntime.GetRuntimeStatus(containerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// executeCommand 执行命令
func (a *APIServer) executeCommand(c *gin.Context) {
	var req struct {
		ContainerName string   `json:"container_name" binding:"required"`
		Command       []string `json:"command" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := a.sandboxRuntime.ExecuteCommand(req.ContainerName, req.Command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": output})
}

// stopRuntime 停止运行时
func (a *APIServer) stopRuntime(c *gin.Context) {
	var req struct {
		ContainerName string `json:"container_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := a.sandboxRuntime.StopRuntime(req.ContainerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Runtime for container %s stopped successfully", req.ContainerName)})
}

// getMetrics 获取监控指标
func (a *APIServer) getMetrics(c *gin.Context) {
	// 重定向到Prometheus监控端点
	c.Redirect(http.StatusMovedPermanently, "/metrics")
}

// getMonitorStatus 获取监控状态
func (a *APIServer) getMonitorStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "monitoring",
	})
}
