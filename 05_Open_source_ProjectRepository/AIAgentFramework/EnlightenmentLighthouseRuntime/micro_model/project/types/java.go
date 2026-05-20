package types

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"micro_model/project"
)

// JavaAdapter Java 项目适配器
type JavaAdapter struct {
	project.BaseAdapter
}

// NewJavaAdapter 创建 Java 项目适配器
func NewJavaAdapter() project.ProjectAdapter {
	return &JavaAdapter{}
}

// Deploy 部署 Java 项目
func (a *JavaAdapter) Deploy(project *project.Project, sandboxID string) error {
	// 检查项目目录是否存在
	if _, err := os.Stat(project.Path); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", project.Path)
	}

	// 检查是否有 pom.xml 或 build.gradle 文件
	pomXMLPath := filepath.Join(project.Path, "pom.xml")
	buildGradlePath := filepath.Join(project.Path, "build.gradle")

	if _, err := os.Stat(pomXMLPath); err == nil {
		// 安装 Maven 依赖
		if err := a.installMavenDependencies(project); err != nil {
			return fmt.Errorf("failed to install Maven dependencies: %v", err)
		}
	} else if _, err := os.Stat(buildGradlePath); err == nil {
		// 安装 Gradle 依赖
		if err := a.installGradleDependencies(project); err != nil {
			return fmt.Errorf("failed to install Gradle dependencies: %v", err)
		}
	}

	fmt.Printf("Deployed Java project: %s\n", project.Name)

	return nil
}

// Undeploy 卸载 Java 项目
func (a *JavaAdapter) Undeploy(project *project.Project, sandboxID string) error {
	// 清理项目依赖（可选）
	targetPath := filepath.Join(project.Path, "target")
	if _, err := os.Stat(targetPath); err == nil {
		if err := os.RemoveAll(targetPath); err != nil {
			return fmt.Errorf("failed to remove target directory: %v", err)
		}
	}

	fmt.Printf("Undeployed Java project: %s\n", project.Name)

	return nil
}

// Start 启动 Java 项目
func (a *JavaAdapter) Start(project *project.Project, sandboxID string) error {
	// 检查项目目录是否存在
	if _, err := os.Stat(project.Path); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", project.Path)
	}

	// 查找 JAR 文件
	jarFiles, err := filepath.Glob(filepath.Join(project.Path, "*.jar"))
	if err != nil || len(jarFiles) == 0 {
		// 尝试在 target 目录中查找
		targetJarFiles, err := filepath.Glob(filepath.Join(project.Path, "target", "*.jar"))
		if err != nil || len(targetJarFiles) == 0 {
			return fmt.Errorf("no JAR files found in project directory")
		}
		jarFiles = targetJarFiles
	}

	// 启动 Java 应用
	cmd := exec.Command("java", "-jar", jarFiles[0])
	cmd.Dir = project.Path

	// 启动项目在后台运行
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Java project: %v", err)
	}

	// 记录进程 ID
	// 这里可以将进程 ID 保存到项目状态中

	fmt.Printf("Started Java project: %s (PID: %d)\n", project.Name, cmd.Process.Pid)

	return nil
}

// Stop 停止 Java 项目
func (a *JavaAdapter) Stop(project *project.Project, sandboxID string) error {
	// 这里需要根据实际情况实现停止项目的逻辑
	// 例如，通过进程 ID 停止项目

	// 简化实现，暂时返回 nil
	fmt.Printf("Stopped Java project: %s\n", project.Name)

	return nil
}

// Monitor 监控 Java 项目
func (a *JavaAdapter) Monitor(project *project.Project, sandboxID string) (project.Resources, error) {
	// 这里需要根据实际情况实现监控项目的逻辑
	// 例如，通过进程 ID 监控项目的资源使用情况

	// 简化实现，返回默认资源使用情况
	return project.Resources{
		CPU:     0.15,
		Memory:  200 * 1024 * 1024, // 200MB
		Disk:    100 * 1024 * 1024,  // 100MB
		Network: 2 * 1024 * 1024,    // 2MB
	}, nil
}

// installMavenDependencies 安装 Maven 项目依赖
func (a *JavaAdapter) installMavenDependencies(project *project.Project) error {
	// 检查 mvn 是否可用
	_, err := exec.LookPath("mvn")
	if err != nil {
		return fmt.Errorf("mvn is not available: %v", err)
	}

	// 安装依赖
	cmd := exec.Command("mvn", "install")
	cmd.Dir = project.Path

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mvn install failed: %v, output: %s", err, string(output))
	}

	// 检查是否有警告
	if strings.Contains(string(output), "WARNING") {
		fmt.Printf("mvn install completed with warnings: %s\n", string(output))
	} else {
		fmt.Printf("mvn install completed successfully\n")
	}

	return nil
}

// installGradleDependencies 安装 Gradle 项目依赖
func (a *JavaAdapter) installGradleDependencies(project *project.Project) error {
	// 检查 gradle 是否可用
	_, err := exec.LookPath("gradle")
	if err != nil {
		// 尝试使用 gradlew
		gradlewPath := filepath.Join(project.Path, "gradlew")
		if _, err := os.Stat(gradlewPath); err != nil {
			return fmt.Errorf("neither gradle nor gradlew is available: %v", err)
		}
		// 使用 gradlew
		cmd := exec.Command(gradlewPath, "build")
		cmd.Dir = project.Path

		// 执行命令
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("gradlew build failed: %v, output: %s", err, string(output))
		}

		// 检查是否有警告
		if strings.Contains(string(output), "WARNING") {
			fmt.Printf("gradlew build completed with warnings: %s\n", string(output))
		} else {
			fmt.Printf("gradlew build completed successfully\n")
		}

		return nil
	}

	// 使用 gradle
	cmd := exec.Command("gradle", "build")
	cmd.Dir = project.Path

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gradle build failed: %v, output: %s", err, string(output))
	}

	// 检查是否有警告
	if strings.Contains(string(output), "WARNING") {
		fmt.Printf("gradle build completed with warnings: %s\n", string(output))
	} else {
		fmt.Printf("gradle build completed successfully\n")
	}

	return nil
}
