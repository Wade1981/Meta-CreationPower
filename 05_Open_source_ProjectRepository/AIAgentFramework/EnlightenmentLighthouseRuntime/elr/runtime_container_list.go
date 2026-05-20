package elr

import (
	"encoding/json"
	"os"
	"sync"
)

// RuntimeContainerList manages the list of running containers
type RuntimeContainerList struct {
	containers map[string]*RunningContainerInfo
	mutex      sync.RWMutex
	filePath   string
}

// RunningContainerInfo represents information about a running container
type RunningContainerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	PID  int    `json:"pid"`
}

// Global runtime container list instance
var (
	runtimeContainerList *RuntimeContainerList
	runtimeContainerListOnce sync.Once
)

// GetRuntimeContainerList returns the global runtime container list instance
func GetRuntimeContainerList() *RuntimeContainerList {
	runtimeContainerListOnce.Do(func() {
		// Create runtime container list without file path (in-memory only)
		runtimeContainerList = &RuntimeContainerList{
			containers: make(map[string]*RunningContainerInfo),
			filePath:   "",
		}
		// Don't load from file - runtime container list should be in-memory only
	})
	return runtimeContainerList
}

// AddContainer adds a container to the runtime container list
func (rcl *RuntimeContainerList) AddContainer(id, name string, pid int) {
	rcl.mutex.Lock()
	defer rcl.mutex.Unlock()
	rcl.containers[id] = &RunningContainerInfo{
		ID:   id,
		Name: name,
		PID:  pid,
	}
	rcl.save()
}

// RemoveContainer removes a container from the runtime container list
func (rcl *RuntimeContainerList) RemoveContainer(id string) {
	rcl.mutex.Lock()
	defer rcl.mutex.Unlock()
	delete(rcl.containers, id)
	rcl.save()
}

// GetContainer gets a container from the runtime container list
func (rcl *RuntimeContainerList) GetContainer(id string) (*RunningContainerInfo, bool) {
	rcl.mutex.RLock()
	defer rcl.mutex.RUnlock()
	container, exists := rcl.containers[id]
	return container, exists
}

// ListContainers returns all containers in the runtime container list
func (rcl *RuntimeContainerList) ListContainers() []*RunningContainerInfo {
	rcl.mutex.RLock()
	defer rcl.mutex.RUnlock()
	containers := make([]*RunningContainerInfo, 0, len(rcl.containers))
	for _, container := range rcl.containers {
		containers = append(containers, container)
	}
	return containers
}

// Clear clears the runtime container list
func (rcl *RuntimeContainerList) Clear() {
	rcl.mutex.Lock()
	defer rcl.mutex.Unlock()
	clear(rcl.containers)
	rcl.save()
}

// Size returns the size of the runtime container list
func (rcl *RuntimeContainerList) Size() int {
	rcl.mutex.RLock()
	defer rcl.mutex.RUnlock()
	return len(rcl.containers)
}

// save saves the runtime container list to file
func (rcl *RuntimeContainerList) save() {
	// Don't save to file if filePath is empty (in-memory only)
	if rcl.filePath == "" {
		return
	}
	
	data, err := json.MarshalIndent(rcl.containers, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(rcl.filePath, data, 0644)
}

// load loads the runtime container list from file
func (rcl *RuntimeContainerList) load() {
	// Don't load from file if filePath is empty (in-memory only)
	if rcl.filePath == "" {
		return
	}
	
	data, err := os.ReadFile(rcl.filePath)
	if err != nil {
		return
	}
	json.Unmarshal(data, &rcl.containers)
}
