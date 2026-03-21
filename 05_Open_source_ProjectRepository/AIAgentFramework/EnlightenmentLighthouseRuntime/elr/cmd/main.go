package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"elr"
)

func main() {
	// 打印启动信息
	fmt.Println("Starting Enlightenment Lighthouse Runtime...")

	// 创建默认配置
	config := &elr.Config{
		LogLevel: "info",
	}

	// 创建运行时
	runtime, err := elr.NewRuntime(config)
	if err != nil {
		fmt.Printf("Failed to create runtime: %v\n", err)
		os.Exit(1)
	}

	// 启动运行时
	if err := runtime.Start(); err != nil {
		fmt.Printf("Failed to start runtime: %v\n", err)
		os.Exit(1)
	}

	// 等待信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 停止运行时
	fmt.Println("Shutting down Enlightenment Lighthouse Runtime...")
	if err := runtime.Stop(); err != nil {
		fmt.Printf("Failed to stop runtime: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Enlightenment Lighthouse Runtime stopped successfully!")
}
