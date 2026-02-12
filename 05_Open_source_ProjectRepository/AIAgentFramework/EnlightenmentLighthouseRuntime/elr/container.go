// Package elr implements the container management for Enlightenment Lighthouse Runtime
package elr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ContainerStatus represents the status of a container
type ContainerStatus string

const (
	// ContainerStatusCreated represents a container that has been created but not started
	ContainerStatusCreated ContainerStatus = "created"
	// ContainerStatusRunning represents a container that is currently running
	ContainerStatusRunning ContainerStatus = "running"
	// ContainerStatusStopped represents a container that has been stopped
	ContainerStatusStopped ContainerStatus = "stopped"
	// ContainerStatusPaused represents a container that is paused
	ContainerStatusPaused ContainerStatus = "paused"
	// ContainerStatusError represents a container that has encountered an error
	ContainerStatusError ContainerStatus = "error"
)

// ContainerConfig represents the configuration for a container
type ContainerConfig struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

// Container represents a container instance
type Container struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	Dir       string            `json:"dir"`
	Runtime   *Runtime          `json:"-"`
	Status    ContainerStatus   `json:"status"`
	Created   time.Time         `json:"created"`
	Started   *time.Time        `json:"started,omitempty"`
	Stopped   *time.Time        `json:"stopped,omitempty"`
	PID       int               `json:"pid,omitempty"`
	ExitCode  int               `json:"exit_code,omitempty"`
	Error     string            `json:"error,omitempty"`
}

// Start starts the container
func (c *Container) Start() error {
	if c.Status == ContainerStatusRunning {
		return fmt.Errorf("container is already running")
	}

	if c.Status == ContainerStatusError {
		return fmt.Errorf("container is in error state")
	}

	// Update status
	c.Status = ContainerStatusRunning
	startTime := time.Now()
	c.Started = &startTime

	// Save container config
	if err := c.saveConfig(); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to save container config: %v", err)
		c.saveConfig()
		return c.Error
	}

	fmt.Printf("Starting container: %s (%s)\n", c.ID, c.Name)

	// TODO: Implement container start logic
	// This is a placeholder for now

	// Simulate container start
	c.PID = os.Getpid() + 1000

	fmt.Printf("Started container: %s (%s) with PID %d\n", c.ID, c.Name, c.PID)
	return nil
}

// Stop stops the container
func (c *Container) Stop() error {
	if c.Status != ContainerStatusRunning {
		return fmt.Errorf("container is not running")
	}

	// Update status
	c.Status = ContainerStatusStopped
	stopTime := time.Now()
	c.Stopped = &stopTime
	c.ExitCode = 0

	// Save container config
	if err := c.saveConfig(); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to save container config: %v", err)
		c.saveConfig()
		return c.Error
	}

	fmt.Printf("Stopping container: %s (%s)\n", c.ID, c.Name)

	// TODO: Implement container stop logic
	// This is a placeholder for now

	fmt.Printf("Stopped container: %s (%s) with exit code %d\n", c.ID, c.Name, c.ExitCode)
	return nil
}

// Pause pauses the container
func (c *Container) Pause() error {
	if c.Status != ContainerStatusRunning {
		return fmt.Errorf("container is not running")
	}

	// Update status
	c.Status = ContainerStatusPaused

	// Save container config
	if err := c.saveConfig(); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to save container config: %v", err)
		c.saveConfig()
		return c.Error
	}

	fmt.Printf("Paused container: %s (%s)\n", c.ID, c.Name)

	// TODO: Implement container pause logic
	// This is a placeholder for now

	return nil
}

// Unpause unpauses the container
func (c *Container) Unpause() error {
	if c.Status != ContainerStatusPaused {
		return fmt.Errorf("container is not paused")
	}

	// Update status
	c.Status = ContainerStatusRunning

	// Save container config
	if err := c.saveConfig(); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to save container config: %v", err)
		c.saveConfig()
		return c.Error
	}

	fmt.Printf("Unpaused container: %s (%s)\n", c.ID, c.Name)

	// TODO: Implement container unpause logic
	// This is a placeholder for now

	return nil
}

// Restart restarts the container
func (c *Container) Restart() error {
	// Stop container if running
	if c.Status == ContainerStatusRunning {
		if err := c.Stop(); err != nil {
			return fmt.Errorf("failed to stop container: %v", err)
		}
	}

	// Start container
	return c.Start()
}

// Remove removes the container
func (c *Container) Remove() error {
	return c.Runtime.DeleteContainer(c.ID)
}

// saveConfig saves the container configuration to disk
func (c *Container) saveConfig() error {
	configPath := filepath.Join(c.Dir, "config.json")

	// Create directory if not exists
	if err := os.MkdirAll(c.Dir, 0755); err != nil {
		return fmt.Errorf("failed to create container directory: %v", err)
	}

	// Marshal container to JSON
	configJSON, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal container config: %v", err)
	}

	// Write config to file
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write container config: %v", err)
	}

	return nil
}

// loadContainer loads a container from disk
func loadContainer(containerDir string, runtime *Runtime) (*Container, error) {
	configPath := filepath.Join(containerDir, "config.json")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("container config file not found")
	}

	// Read config file
	configJSON, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read container config: %v", err)
	}

	// Unmarshal config
	container := &Container{}
	if err := json.Unmarshal(configJSON, container); err != nil {
		return nil, fmt.Errorf("failed to unmarshal container config: %v", err)
	}

	// Set runtime
	container.Runtime = runtime

	// Set directory
	container.Dir = containerDir

	return container, nil
}

// String returns a string representation of the container
func (c *Container) String() string {
	return fmt.Sprintf("%s (%s) - %s", c.ID, c.Name, c.Status)
}
