package api

import (
	"fmt"
	"sync"
)

// APIService API服务接口
type APIService interface {
	Start() error
	Stop() error
}

// APIServiceObserver API服务观察者接口
type APIServiceObserver interface {
	OnServiceStarted(serviceType string, port int)
	OnServiceStopped(serviceType string)
}

// APIServiceManager API服务管理器
type APIServiceManager struct {
	services  map[string]APIService
	observers []APIServiceObserver
	mu        sync.Mutex
}

// NewAPIServiceManager 创建API服务管理器
func NewAPIServiceManager() *APIServiceManager {
	return &APIServiceManager{
		services:  make(map[string]APIService),
		observers: make([]APIServiceObserver, 0),
	}
}

// RegisterService 注册API服务
func (m *APIServiceManager) RegisterService(serviceType string, service APIService) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.services[serviceType] = service
}

// UnregisterService 注销API服务
func (m *APIServiceManager) UnregisterService(serviceType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.services, serviceType)
}

// StartService 启动指定的API服务
func (m *APIServiceManager) StartService(serviceType string) error {
	m.mu.Lock()
	service, exists := m.services[serviceType]
	m.mu.Unlock()
	
	if !exists {
		return fmt.Errorf("service type %s not registered", serviceType)
	}
	
	if err := service.Start(); err != nil {
		return err
	}
	
	// 通知观察者服务已启动
	m.notifyObserversOnStart(serviceType, 0) // 这里需要根据实际情况获取端口
	
	return nil
}

// StopService 停止指定的API服务
func (m *APIServiceManager) StopService(serviceType string) error {
	m.mu.Lock()
	service, exists := m.services[serviceType]
	m.mu.Unlock()
	
	if !exists {
		return fmt.Errorf("service type %s not registered", serviceType)
	}
	
	if err := service.Stop(); err != nil {
		return err
	}
	
	// 通知观察者服务已停止
	m.notifyObserversOnStop(serviceType)
	
	return nil
}

// StartAllServices 启动所有API服务
func (m *APIServiceManager) StartAllServices() error {
	m.mu.Lock()
	services := make(map[string]APIService, len(m.services))
	for k, v := range m.services {
		services[k] = v
	}
	m.mu.Unlock()
	
	for serviceType, service := range services {
		if err := service.Start(); err != nil {
			return err
		}
		// 通知观察者服务已启动
		m.notifyObserversOnStart(serviceType, 0) // 这里需要根据实际情况获取端口
	}
	
	return nil
}

// StopAllServices 停止所有API服务
func (m *APIServiceManager) StopAllServices() error {
	m.mu.Lock()
	services := make(map[string]APIService, len(m.services))
	for k, v := range m.services {
		services[k] = v
	}
	m.mu.Unlock()
	
	for serviceType, service := range services {
		if err := service.Stop(); err != nil {
			return err
		}
		// 通知观察者服务已停止
		m.notifyObserversOnStop(serviceType)
	}
	
	return nil
}

// AddObserver 添加观察者
func (m *APIServiceManager) AddObserver(observer APIServiceObserver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.observers = append(m.observers, observer)
}

// RemoveObserver 移除观察者
func (m *APIServiceManager) RemoveObserver(observer APIServiceObserver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, obs := range m.observers {
		if obs == observer {
			m.observers = append(m.observers[:i], m.observers[i+1:]...)
			break
		}
	}
}

// notifyObserversOnStart 通知观察者服务已启动
func (m *APIServiceManager) notifyObserversOnStart(serviceType string, port int) {
	m.mu.Lock()
	observers := make([]APIServiceObserver, len(m.observers))
	copy(observers, m.observers)
	m.mu.Unlock()
	
	for _, observer := range observers {
		observer.OnServiceStarted(serviceType, port)
	}
}

// notifyObserversOnStop 通知观察者服务已停止
func (m *APIServiceManager) notifyObserversOnStop(serviceType string) {
	m.mu.Lock()
	observers := make([]APIServiceObserver, len(m.observers))
	copy(observers, m.observers)
	m.mu.Unlock()
	
	for _, observer := range observers {
		observer.OnServiceStopped(serviceType)
	}
}

// GetService 获取指定的API服务
func (m *APIServiceManager) GetService(serviceType string) (APIService, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	service, exists := m.services[serviceType]
	return service, exists
}

// ListServices 列出所有注册的API服务
func (m *APIServiceManager) ListServices() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	serviceTypes := make([]string, 0, len(m.services))
	for serviceType := range m.services {
		serviceTypes = append(serviceTypes, serviceType)
	}
	return serviceTypes
}
