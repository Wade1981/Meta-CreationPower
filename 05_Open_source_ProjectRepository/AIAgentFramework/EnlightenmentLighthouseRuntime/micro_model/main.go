package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"micro-model/api"
	"micro-model/config"
	"micro-model/container"
	"micro-model/model"
	"micro-model/monitor"
	"micro-model/sandbox"
)

func main() {
	port := flag.Int("port", 0, "Server port (overrides config)")
	ip := flag.String("ip", "", "Server IP address (overrides config)")
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *port > 0 {
		cfg.Server.Port = *port
	}
	if *ip != "" {
		cfg.Server.Host = *ip
	}

	if len(flag.Args()) > 0 {
		if p, err := strconv.Atoi(flag.Arg(0)); err == nil && p > 0 {
			cfg.Server.Port = p
		}
	}

	setupLogger()

	modelManager, err := model.NewModelManager(&cfg.Model)
	if err != nil {
		log.Fatalf("Failed to initialize model manager: %v", err)
	}

	containerManager, err := container.NewContainerManager(&cfg.Container)
	if err != nil {
		log.Fatalf("Failed to initialize container manager: %v", err)
	}

	sandboxRuntime, err := sandbox.NewSandboxRuntime(&cfg.Sandbox)
	if err != nil {
		log.Fatalf("Failed to initialize sandbox runtime: %v", err)
	}

	monitorService, err := monitor.NewMonitorService(&cfg.Monitoring)
	if err != nil {
		log.Fatalf("Failed to initialize monitor service: %v", err)
	}

	apiServer := api.NewAPIServer(&cfg.Server, modelManager, containerManager, sandboxRuntime, monitorService)

	go monitorService.Start()

	go func() {
		if err := apiServer.Start(); err != nil {
			log.Printf("Failed to start API server: %v", err)
		}
	}()

	fmt.Printf("Micro Model Server started on %s:%d\n", cfg.Server.Host, cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down services...")

	if err := apiServer.Stop(); err != nil {
		log.Printf("Error stopping API server: %v", err)
	}

	monitorService.Stop()

	containerManager.Cleanup()

	log.Println("Services stopped successfully")
}

func setupLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logDir := "./logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
		}
	}

	logFile, err := os.OpenFile(logDir+"/micro_model.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}

	log.SetOutput(logFile)
}
