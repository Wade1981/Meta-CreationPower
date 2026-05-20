package types

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"micro_model/project"
)

// NodeJSAdapter Node.js 项目适配器
type NodeJSAdapter struct {
	project.BaseAdapter
}

// NewNodeJSAdapter 创建 Node.js 项目适配器
func NewNodeJSAdapter() project.ProjectAdapter {
	return &NodeJSAdapter{}
}

// Deploy 部署 Node.js 项目
func (a *NodeJSAdapter) Deploy(project *project.Project, sandboxID string) error {
	// 检查项目目录是否存在
	if _, err := os.Stat(project.Path); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", project.Path)
	}

	// 检查 package.json 文件是否存在
	packageJSONPath := filepath.Join(project.Path, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found in project directory")
	}

	// 安装 npm 依赖
	if err := a.installDependencies(project); err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	fmt.Printf("Deployed Node.js project: %s\n", project.Name)

	return nil
}

// Undeploy 卸载 Node.js 项目
func (a *NodeJSAdapter) Undeploy(project *project.Project, sandboxID string) error {
	// 清理项目依赖（可选）
	nodeModulesPath := filepath.Join(project.Path, "node_modules")
	if _, err := os.Stat(nodeModulesPath); err == nil {
		if err := os.RemoveAll(nodeModulesPath); err != nil {
			return fmt.Errorf("failed to remove node_modules: %v", err)
		}
	}

	fmt.Printf("Undeployed Node.js project: %s\n", project.Name)

	return nil
}

// Start 启动 Node.js 项目
func (a *NodeJSAdapter) Start(project *project.Project, sandboxID string) error {
	// 检查项目目录是否存在
	if _, err := os.Stat(project.Path); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", project.Path)
	}

	// 检查 package.json 文件是否存在
	packageJSONPath := filepath.Join(project.Path, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found in project directory")
	}

	// 启动项目
	cmd := exec.Command("npm", "start")
	cmd.Dir = project.Path

	// 启动项目在后台运行
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Node.js project: %v", err)
	}

	// 记录进程 ID
	// 这里可以将进程 ID 保存到项目状态中

	fmt.Printf("Started Node.js project: %s (PID: %d)\n", project.Name, cmd.Process.Pid)

	return nil
}

// Stop 停止 Node.js 项目
func (a *NodeJSAdapter) Stop(project *project.Project, sandboxID string) error {
	// 这里需要根据实际情况实现停止项目的逻辑
	// 例如，通过进程 ID 停止项目

	// 简化实现，暂时返回 nil
	fmt.Printf("Stopped Node.js project: %s\n", project.Name)

	return nil
}

// Monitor 监控 Node.js 项目
func (a *NodeJSAdapter) Monitor(project *project.Project, sandboxID string) (project.Resources, error) {
	// 这里需要根据实际情况实现监控项目的逻辑
	// 例如，通过进程 ID 监控项目的资源使用情况

	// 简化实现，返回默认资源使用情况
	return project.Resources{
		CPU:     0.1,
		Memory:  100 * 1024 * 1024, // 100MB
		Disk:    50 * 1024 * 1024,  // 50MB
		Network: 1 * 1024 * 1024,   // 1MB
	}, nil
}

// installDependencies 安装 Node.js 项目依赖
func (a *NodeJSAdapter) installDependencies(project *project.Project) error {
	// 检查 npm 是否可用
	_, err := exec.LookPath("npm")
	if err != nil {
		return fmt.Errorf("npm is not available: %v", err)
	}

	// 安装依赖
	cmd := exec.Command("npm", "install")
	cmd.Dir = project.Path

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm install failed: %v, output: %s", err, string(output))
	}

	// 检查是否有警告
	if strings.Contains(string(output), "warning") {
		fmt.Printf("npm install completed with warnings: %s\n", string(output))
	} else {
		fmt.Printf("npm install completed successfully\n")
	}

	return nil
}
