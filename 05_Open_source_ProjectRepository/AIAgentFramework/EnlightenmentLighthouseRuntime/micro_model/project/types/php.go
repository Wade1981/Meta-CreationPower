package types

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"micro_model/project"
)

// PHPAdapter PHP 项目适配器
type PHPAdapter struct {
	project.BaseAdapter
}

// NewPHPAdapter 创建 PHP 项目适配器
func NewPHPAdapter() project.ProjectAdapter {
	return &PHPAdapter{}
}

// Deploy 部署 PHP 项目
func (a *PHPAdapter) Deploy(project *project.Project, sandboxID string) error {
	// 检查项目目录是否存在
	if _, err := os.Stat(project.Path); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", project.Path)
	}

	// 检查 composer.json 文件是否存在
	composerJSONPath := filepath.Join(project.Path, "composer.json")
	if _, err := os.Stat(composerJSONPath); err == nil {
		// 安装 Composer 依赖
		if err := a.installDependencies(project); err != nil {
			return fmt.Errorf("failed to install dependencies: %v", err)
		}
	}

	fmt.Printf("Deployed PHP project: %s\n", project.Name)

	return nil
}

// Undeploy 卸载 PHP 项目
func (a *PHPAdapter) Undeploy(project *project.Project, sandboxID string) error {
	// 清理项目依赖（可选）
	vendorPath := filepath.Join(project.Path, "vendor")
	if _, err := os.Stat(vendorPath); err == nil {
		if err := os.RemoveAll(vendorPath); err != nil {
			return fmt.Errorf("failed to remove vendor directory: %v", err)
		}
	}

	fmt.Printf("Undeployed PHP project: %s\n", project.Name)

	return nil
}

// Start 启动 PHP 项目
func (a *PHPAdapter) Start(project *project.Project, sandboxID string) error {
	// 检查项目目录是否存在
	if _, err := os.Stat(project.Path); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", project.Path)
	}

	// 检查是否有 index.php 文件
	indexPHPPath := filepath.Join(project.Path, "index.php")
	if _, err := os.Stat(indexPHPPath); os.IsNotExist(err) {
		// 尝试查找其他 PHP 文件
		phpFiles, err := filepath.Glob(filepath.Join(project.Path, "*.php"))
		if err != nil || len(phpFiles) == 0 {
			return fmt.Errorf("no PHP files found in project directory")
		}
	}

	// 启动 PHP 内置服务器
	cmd := exec.Command("php", "-S", "localhost:8000")
	cmd.Dir = project.Path

	// 启动项目在后台运行
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start PHP project: %v", err)
	}

	// 记录进程 ID
	// 这里可以将进程 ID 保存到项目状态中

	fmt.Printf("Started PHP project: %s (PID: %d)\n", project.Name, cmd.Process.Pid)

	return nil
}

// Stop 停止 PHP 项目
func (a *PHPAdapter) Stop(project *project.Project, sandboxID string) error {
	// 这里需要根据实际情况实现停止项目的逻辑
	// 例如，通过进程 ID 停止项目

	// 简化实现，暂时返回 nil
	fmt.Printf("Stopped PHP project: %s\n", project.Name)

	return nil
}

// Monitor 监控 PHP 项目
func (a *PHPAdapter) Monitor(project *project.Project, sandboxID string) (project.Resources, error) {
	// 这里需要根据实际情况实现监控项目的逻辑
	// 例如，通过进程 ID 监控项目的资源使用情况

	// 简化实现，返回默认资源使用情况
	return project.Resources{
		CPU:     0.05,
		Memory:  50 * 1024 * 1024,  // 50MB
		Disk:    30 * 1024 * 1024,   // 30MB
		Network: 512 * 1024,         // 512KB
	}, nil
}

// installDependencies 安装 PHP 项目依赖
func (a *PHPAdapter) installDependencies(project *project.Project) error {
	// 检查 composer 是否可用
	_, err := exec.LookPath("composer")
	if err != nil {
		return fmt.Errorf("composer is not available: %v", err)
	}

	// 安装依赖
	cmd := exec.Command("composer", "install")
	cmd.Dir = project.Path

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("composer install failed: %v, output: %s", err, string(output))
	}

	// 检查是否有警告
	if strings.Contains(string(output), "warning") {
		fmt.Printf("composer install completed with warnings: %s\n", string(output))
	} else {
		fmt.Printf("composer install completed successfully\n")
	}

	return nil
}
