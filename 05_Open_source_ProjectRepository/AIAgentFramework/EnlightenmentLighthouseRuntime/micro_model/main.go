package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"micro-model/api"
	"micro-model/config"
	"micro-model/container"
	"micro-model/model"
	"micro-model/monitor"
	"micro-model/sandbox"
)

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	setupLogger()

	// 初始化模型管理器
	modelManager, err := model.NewModelManager(cfg.Model)
	if err != nil {
		log.Fatalf("Failed to initialize model manager: %v", err)
	}

	// 初始化容器管理器
	containerManager, err := container.NewContainerManager(cfg.Container)
	if err != nil {
		log.Fatalf("Failed to initialize container manager: %v", err)
	}

	// 初始化沙箱运行时
	sandboxRuntime, err := sandbox.NewSandboxRuntime(cfg.Sandbox)
	if err != nil {
		log.Fatalf("Failed to initialize sandbox runtime: %v", err)
	}

	// 初始化监控服务
	monitorService, err := monitor.NewMonitorService(cfg.Monitoring)
	if err != nil {
		log.Fatalf("Failed to initialize monitor service: %v", err)
	}

	// 初始化API服务
	apiServer := api.NewAPIServer(cfg.Server, modelManager, containerManager, sandboxRuntime, monitorService)

	// 启动监控服务
	go monitorService.Start()

	// 启动API服务
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭服务
	log.Println("Shutting down services...")

	// 停止API服务
	if err := apiServer.Stop(); err != nil {
		log.Printf("Error stopping API server: %v", err)
	}

	// 停止监控服务
	monitorService.Stop()

	// 清理资源
	containerManager.Cleanup()

	log.Println("Services stopped successfully")
}

func setupLogger() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 创建日志目录
	logDir := "./logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
		}
	}

	// 打开日志文件
	logFile, err := os.OpenFile(logDir+"/micro_model.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}

	// 设置日志输出
	log.SetOutput(logFile)
	defer logFile.Close()
}
