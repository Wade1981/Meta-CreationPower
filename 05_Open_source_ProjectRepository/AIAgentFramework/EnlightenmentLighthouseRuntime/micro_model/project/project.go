package project

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"micro_model/config"
)

// ProjectType 项目类型
type ProjectType string

const (
	// ProjectTypeNodeJS Node.js 项目
	ProjectTypeNodeJS ProjectType = "nodejs"
	// ProjectTypePHP PHP 项目
	ProjectTypePHP ProjectType = "php"
	// ProjectTypeJava Java 项目
	ProjectTypeJava ProjectType = "java"
)

// ProjectStatus 项目状态
type ProjectStatus string

const (
	// ProjectStatusCreated 创建状态
	ProjectStatusCreated ProjectStatus = "created"
	// ProjectStatusDeployed 已部署状态
	ProjectStatusDeployed ProjectStatus = "deployed"
	// ProjectStatusRunning 运行中状态
	ProjectStatusRunning ProjectStatus = "running"
	// ProjectStatusStopped 已停止状态
	ProjectStatusStopped ProjectStatus = "stopped"
	// ProjectStatusError 错误状态
	ProjectStatusError ProjectStatus = "error"
)

// Project 项目信息
type Project struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Type        ProjectType   `json:"type"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Path        string        `json:"path"`
	Status      ProjectStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Resources   Resources     `json:"resources"`
	Properties  *Properties   `json:"properties"`
}

// Resources 资源使用情况
type Resources struct {
	CPU    float64 `json:"cpu"`
	Memory int64   `json:"memory"`
	Disk   int64   `json:"disk"`
	Network int64  `json:"network"`
}

// Properties 项目属性
type Properties struct {
	Dependencies Dependencies `json:"dependencies"`
	Config       map[string]interface{} `json:"config"`
}

// Dependencies 依赖管理
type Dependencies struct {
	Npm     []string `json:"npm"`
	Composer []string `json:"composer"`
	Maven    []string `json:"maven"`
	Gradle   []string `json:"gradle"`
}

// ProjectConfig 项目配置
type ProjectConfig struct {
	Name        string        `json:"name"`
	Type        ProjectType   `json:"type"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Path        string        `json:"path"`
	Properties  *Properties   `json:"properties"`
}

// ProjectManager 项目管理器
type ProjectManager struct {
	config   *config.Config
	projects map[string]*Project
	adapters map[ProjectType]ProjectAdapter
	mutex    sync.RWMutex
}

// NewProjectManager 创建项目管理器
func NewProjectManager(config *config.Config) *ProjectManager {
	manager := &ProjectManager{
		config:   config,
		projects: make(map[string]*Project),
		adapters: make(map[ProjectType]ProjectAdapter),
	}

	// 注册项目适配器
	manager.registerAdapters()

	// 从磁盘加载项目
	if err := manager.LoadProjects(); err != nil {
		fmt.Printf("Warning: failed to load projects: %v\n", err)
	}

	return manager
}

// registerAdapters 注册项目适配器
func (pm *ProjectManager) registerAdapters() {
	// 注册 Node.js 项目适配器
	pm.adapters[ProjectTypeNodeJS] = &NodeJSAdapter{}

	// 注册 PHP 项目适配器
	pm.adapters[ProjectTypePHP] = &PHPAdapter{}

	// 注册 Java 项目适配器
	pm.adapters[ProjectTypeJava] = &JavaAdapter{}
}

// CreateProject 创建项目
func (pm *ProjectManager) CreateProject(config ProjectConfig) (*Project, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// 生成项目 ID
	projectID := fmt.Sprintf("project-%d", time.Now().UnixNano()/1000000)

	// 创建项目实例
	project := &Project{
		ID:          projectID,
		Name:        config.Name,
		Type:        config.Type,
		Version:     config.Version,
		Description: config.Description,
		Path:        config.Path,
		Status:      ProjectStatusCreated,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Resources: Resources{
			CPU:     0,
			Memory:  0,
			Disk:    0,
			Network: 0,
		},
		Properties: config.Properties,
	}

	// 保存项目
	pm.projects[projectID] = project

	// 确保项目目录存在
	if err := os.MkdirAll(project.Path, 0755); err != nil {
		project.Status = ProjectStatusError
		return project, fmt.Errorf("failed to create project directory: %v", err)
	}

	// 保存项目到磁盘
	if err := pm.SaveProjects(); err != nil {
		fmt.Printf("Warning: failed to save projects: %v\n", err)
	}

	fmt.Printf("Created project: %s (%s)\n", project.ID, project.Name)

	return project, nil
}

// GetProject 获取项目
func (pm *ProjectManager) GetProject(projectID string) (*Project, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	project, exists := pm.projects[projectID]
	if !exists {
		return nil, fmt.Errorf("project %s not found", projectID)
	}

	return project, nil
}

// ListProjects 列出所有项目
func (pm *ProjectManager) ListProjects() []*Project {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	projects := make([]*Project, 0, len(pm.projects))
	for _, project := range pm.projects {
		projects = append(projects, project)
	}

	return projects
}

// DeleteProject 删除项目
func (pm *ProjectManager) DeleteProject(projectID string) error {
	pm.mutex.Lock()
	project, exists := pm.projects[projectID]
	if !exists {
		pm.mutex.Unlock()
		return fmt.Errorf("project %s not found", projectID)
	}

	// 检查项目是否正在运行
	if project.Status == ProjectStatusRunning {
		// 项目正在运行，先停止
		project.Status = ProjectStatusCreated
		fmt.Printf("Stopped running project: %s\n", projectID)
	}

	// 删除项目目录
	if err := os.RemoveAll(project.Path); err != nil {
		pm.mutex.Unlock()
		return fmt.Errorf("failed to delete project directory: %v", err)
	}

	// 从项目列表中移除
	delete(pm.projects, projectID)
	pm.mutex.Unlock()

	// 保存项目到磁盘
	if err := pm.SaveProjects(); err != nil {
		fmt.Printf("Warning: failed to save projects: %v\n", err)
	}

	fmt.Printf("Deleted project: %s\n", projectID)
	return nil
}

// DeployProject 部署项目
func (pm *ProjectManager) DeployProject(projectID string, sandboxID string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	project, exists := pm.projects[projectID]
	if !exists {
		return fmt.Errorf("project %s not found", projectID)
	}

	// 获取项目适配器
	adapter, exists := pm.adapters[project.Type]
	if !exists {
		return fmt.Errorf("no adapter found for project type: %s", project.Type)
	}

	// 部署项目
	if err := adapter.Deploy(project, sandboxID); err != nil {
		project.Status = ProjectStatusError
		return fmt.Errorf("failed to deploy project: %v", err)
	}

	// 更新项目状态
	project.Status = ProjectStatusDeployed
	project.UpdatedAt = time.Now()

	// 保存项目到磁盘
	if err := pm.SaveProjects(); err != nil {
		fmt.Printf("Warning: failed to save projects: %v\n", err)
	}

	fmt.Printf("Deployed project: %s to sandbox: %s\n", projectID, sandboxID)

	return nil
}

// UndeployProject 卸载项目
func (pm *ProjectManager) UndeployProject(projectID string, sandboxID string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	project, exists := pm.projects[projectID]
	if !exists {
		return fmt.Errorf("project %s not found", projectID)
	}

	// 获取项目适配器
	adapter, exists := pm.adapters[project.Type]
	if !exists {
		return fmt.Errorf("no adapter found for project type: %s", project.Type)
	}

	// 卸载项目
	if err := adapter.Undeploy(project, sandboxID); err != nil {
		project.Status = ProjectStatusError
		return fmt.Errorf("failed to undeploy project: %v", err)
	}

	// 更新项目状态
	project.Status = ProjectStatusCreated
	project.UpdatedAt = time.Now()

	// 保存项目到磁盘
	if err := pm.SaveProjects(); err != nil {
		fmt.Printf("Warning: failed to save projects: %v\n", err)
	}

	fmt.Printf("Undeployed project: %s from sandbox: %s\n", projectID, sandboxID)

	return nil
}

// StartProject 启动项目
func (pm *ProjectManager) StartProject(projectID string, sandboxID string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	project, exists := pm.projects[projectID]
	if !exists {
		return fmt.Errorf("project %s not found", projectID)
	}

	// 检查项目状态
	if project.Status != ProjectStatusDeployed {
		return fmt.Errorf("project %s is not deployed", projectID)
	}

	// 获取项目适配器
	adapter, exists := pm.adapters[project.Type]
	if !exists {
		return fmt.Errorf("no adapter found for project type: %s", project.Type)
	}

	// 启动项目
	if err := adapter.Start(project, sandboxID); err != nil {
		project.Status = ProjectStatusError
		return fmt.Errorf("failed to start project: %v", err)
	}

	// 更新项目状态
	project.Status = ProjectStatusRunning
	project.UpdatedAt = time.Now()

	// 保存项目到磁盘
	if err := pm.SaveProjects(); err != nil {
		fmt.Printf("Warning: failed to save projects: %v\n", err)
	}

	fmt.Printf("Started project: %s in sandbox: %s\n", projectID, sandboxID)

	return nil
}

// StopProject 停止项目
func (pm *ProjectManager) StopProject(projectID string, sandboxID string) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	project, exists := pm.projects[projectID]
	if !exists {
		return fmt.Errorf("project %s not found", projectID)
	}

	// 检查项目状态
	if project.Status != ProjectStatusRunning {
		return fmt.Errorf("project %s is not running", projectID)
	}

	// 获取项目适配器
	adapter, exists := pm.adapters[project.Type]
	if !exists {
		return fmt.Errorf("no adapter found for project type: %s", project.Type)
	}

	// 停止项目
	if err := adapter.Stop(project, sandboxID); err != nil {
		project.Status = ProjectStatusError
		return fmt.Errorf("failed to stop project: %v", err)
	}

	// 更新项目状态
	project.Status = ProjectStatusDeployed
	project.UpdatedAt = time.Now()

	// 保存项目到磁盘
	if err := pm.SaveProjects(); err != nil {
		fmt.Printf("Warning: failed to save projects: %v\n", err)
	}

	fmt.Printf("Stopped project: %s in sandbox: %s\n", projectID, sandboxID)

	return nil
}

// GetProjectStatus 获取项目状态
func (pm *ProjectManager) GetProjectStatus(projectID string) (ProjectStatus, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	project, exists := pm.projects[projectID]
	if !exists {
		return ProjectStatusError, fmt.Errorf("project %s not found", projectID)
	}

	return project.Status, nil
}

// MonitorProject 监控项目
func (pm *ProjectManager) MonitorProject(projectID string, sandboxID string) (*Project, error) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	project, exists := pm.projects[projectID]
	if !exists {
		return nil, fmt.Errorf("project %s not found", projectID)
	}

	// 获取项目适配器
	adapter, exists := pm.adapters[project.Type]
	if !exists {
		return nil, fmt.Errorf("no adapter found for project type: %s", project.Type)
	}

	// 监控项目
	resources, err := adapter.Monitor(project, sandboxID)
	if err != nil {
		return project, fmt.Errorf("failed to monitor project: %v", err)
	}

	// 更新项目资源使用情况
	project.Resources = resources

	return project, nil
}

// LoadProjects 从磁盘加载项目
func (pm *ProjectManager) LoadProjects() error {
	// 检查存储文件是否存在
	storagePath := "./projects-state.json"
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		return nil
	}

	// 读取存储文件
	data, err := os.ReadFile(storagePath)
	if err != nil {
		return fmt.Errorf("failed to read projects state file: %v", err)
	}

	// 解析项目信息
	var projects []*Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return fmt.Errorf("failed to unmarshal projects state: %v", err)
	}

	// 加载项目信息到内存
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	for _, project := range projects {
		pm.projects[project.ID] = project
	}

	fmt.Printf("Loaded %d projects from persistent storage\n", len(projects))
	return nil
}

// SaveProjects 保存项目到磁盘
func (pm *ProjectManager) SaveProjects() error {
	// 直接获取项目列表，不获取锁
	// 因为这个方法通常在已经持有锁的情况下调用
	projects := make([]*Project, 0, len(pm.projects))
	for _, project := range pm.projects {
		projects = append(projects, project)
	}

	// 序列化项目信息
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal projects state: %v", err)
	}

	// 写入存储文件
	storagePath := "./projects-state.json"
	if err := os.WriteFile(storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write projects state file: %v", err)
	}

	return nil
}
