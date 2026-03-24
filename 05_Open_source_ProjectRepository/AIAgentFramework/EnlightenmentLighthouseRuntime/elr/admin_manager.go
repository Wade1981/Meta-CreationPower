// Package elr implements admin management for Enlightenment Lighthouse Runtime
package elr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// AdminRole represents the role of an admin
type AdminRole string

const (
	// AdminRoleSuper represents a super admin with full permissions
	AdminRoleSuper AdminRole = "super"
	// AdminRoleRegular represents a regular admin with limited permissions
	AdminRoleRegular AdminRole = "regular"
)

// AdminStatus represents the status of an admin
type AdminStatus string

const (
	// AdminStatusActive represents an active admin
	AdminStatusActive AdminStatus = "active"
	// AdminStatusInactive represents an inactive admin
	AdminStatusInactive AdminStatus = "inactive"
)

// AdminPermission represents a permission for an admin
type AdminPermission struct {
	ContainerID string `json:"container_id"`
	CanManage   bool   `json:"can_manage"`
	CanRead     bool   `json:"can_read"`
	CanWrite    bool   `json:"can_write"`
}

// Admin represents an admin user
type Admin struct {
	ID           string             `json:"id"`
	Username     string             `json:"username"`
	Role         AdminRole          `json:"role"`
	Status       AdminStatus        `json:"status"`
	Token        string             `json:"token"`
	Permissions  []AdminPermission  `json:"permissions"`
	CreatedAt    time.Time          `json:"created_at"`
	LastLoginAt  *time.Time         `json:"last_login_at,omitempty"`
}

// AdminManager represents the admin manager
type AdminManager struct {
	adminFile string
	admins    []Admin
}

// NewAdminManager creates a new admin manager
func NewAdminManager(adminFile string) *AdminManager {
	return &AdminManager{
		adminFile: adminFile,
		admins:    []Admin{},
	}
}

// LoadAdmins loads admins from file
func (am *AdminManager) LoadAdmins() error {
	if _, err := os.Stat(am.adminFile); os.IsNotExist(err) {
		// Admin file doesn't exist, initialize with empty admins
		am.admins = []Admin{}
		return nil
	}

	// Read encrypted data
	encryptedData, err := os.ReadFile(am.adminFile)
	if err != nil {
		return fmt.Errorf("failed to read admin file: %w", err)
	}

	// Decrypt data
	decryptedData, err := Decrypt(encryptedData)
	if err != nil {
		return fmt.Errorf("failed to decrypt admin data: %w", err)
	}

	var adminData struct {
		Admins       []Admin `json:"admins"`
		LastUpdated  int64   `json:"last_updated"`
	}

	if err := json.Unmarshal(decryptedData, &adminData); err != nil {
		return fmt.Errorf("failed to unmarshal admin data: %w", err)
	}

	am.admins = adminData.Admins
	return nil
}

// SaveAdmins saves admins to file
func (am *AdminManager) SaveAdmins() error {
	adminData := struct {
		Admins      []Admin `json:"admins"`
		LastUpdated int64   `json:"last_updated"`
	}{
		Admins:      am.admins,
		LastUpdated: time.Now().Unix(),
	}

	data, err := json.MarshalIndent(adminData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal admin data: %w", err)
	}

	// Encrypt data
	encryptedData, err := Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt admin data: %w", err)
	}

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(am.adminFile), 0755); err != nil {
		return fmt.Errorf("failed to create admin directory: %w", err)
	}

	if err := os.WriteFile(am.adminFile, encryptedData, 0644); err != nil {
		return fmt.Errorf("failed to write admin file: %w", err)
	}

	return nil
}

// CreateAdmin creates a new admin
func (am *AdminManager) CreateAdmin(username string, role AdminRole) (string, error) {
	// Check if admin already exists
	for _, admin := range am.admins {
		if admin.Username == username {
			return "", fmt.Errorf("admin with username %s already exists", username)
		}
	}

	// Generate token
	tokenManager := NewTokenManager(filepath.Join(filepath.Dir(am.adminFile), "tokens.json"))
	token, err := tokenManager.GenerateToken(fmt.Sprintf("Admin: %s", username))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// Create admin
	admin := Admin{
		ID:          uuid.New().String(),
		Username:    username,
		Role:        role,
		Status:      AdminStatusActive,
		Token:       token,
		Permissions: []AdminPermission{},
		CreatedAt:   time.Now(),
	}

	am.admins = append(am.admins, admin)

	if err := am.SaveAdmins(); err != nil {
		return "", fmt.Errorf("failed to save admin: %w", err)
	}

	fmt.Printf("Created admin: %s with role: %s\n", username, role)
	return token, nil
}

// GetAdminByUsername gets an admin by username
func (am *AdminManager) GetAdminByUsername(username string) (*Admin, error) {
	for i, admin := range am.admins {
		if admin.Username == username {
			return &am.admins[i], nil
		}
	}

	return nil, fmt.Errorf("admin not found")
}

// GetAdminByToken gets an admin by token
func (am *AdminManager) GetAdminByToken(token string) (*Admin, error) {
	for i, admin := range am.admins {
		if admin.Token == token {
			return &am.admins[i], nil
		}
	}

	return nil, fmt.Errorf("admin not found")
}

// UpdateAdminStatus updates the status of an admin
func (am *AdminManager) UpdateAdminStatus(username string, status AdminStatus) error {
	for i, admin := range am.admins {
		if admin.Username == username {
			am.admins[i].Status = status
			if err := am.SaveAdmins(); err != nil {
				return fmt.Errorf("failed to save admin: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("admin not found")
}

// AddContainerPermission adds a container permission to an admin
func (am *AdminManager) AddContainerPermission(username, containerID string, canManage, canRead, canWrite bool) error {
	for i, admin := range am.admins {
		if admin.Username == username {
			// Check if permission already exists
			for j, perm := range admin.Permissions {
				if perm.ContainerID == containerID {
					// Update existing permission
					am.admins[i].Permissions[j].CanManage = canManage
					am.admins[i].Permissions[j].CanRead = canRead
					am.admins[i].Permissions[j].CanWrite = canWrite
					if err := am.SaveAdmins(); err != nil {
						return fmt.Errorf("failed to save admin: %w", err)
					}
					return nil
				}
			}

			// Add new permission
			permission := AdminPermission{
				ContainerID: containerID,
				CanManage:   canManage,
				CanRead:     canRead,
				CanWrite:    canWrite,
			}
			am.admins[i].Permissions = append(am.admins[i].Permissions, permission)

			if err := am.SaveAdmins(); err != nil {
				return fmt.Errorf("failed to save admin: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("admin not found")
}

// RemoveContainerPermission removes a container permission from an admin
func (am *AdminManager) RemoveContainerPermission(username, containerID string) error {
	for i, admin := range am.admins {
		if admin.Username == username {
			newPermissions := []AdminPermission{}
			for _, perm := range admin.Permissions {
				if perm.ContainerID != containerID {
					newPermissions = append(newPermissions, perm)
				}
			}
			am.admins[i].Permissions = newPermissions

			if err := am.SaveAdmins(); err != nil {
				return fmt.Errorf("failed to save admin: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("admin not found")
}

// ListAdmins lists all admins
func (am *AdminManager) ListAdmins() []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(am.admins))

	for _, admin := range am.admins {
		adminInfo := map[string]interface{}{
			"id":          admin.ID,
			"username":    admin.Username,
			"role":        admin.Role,
			"status":      admin.Status,
			"created_at":  admin.CreatedAt.Unix(),
			"permissions": admin.Permissions,
		}
		result = append(result, adminInfo)
	}

	return result
}

// ValidateAdmin validates an admin token and checks permission for a container
func (am *AdminManager) ValidateAdmin(token, containerID string, action string) (bool, string) {
	// Find admin by token
	admin, err := am.GetAdminByToken(token)
	if err != nil {
		return false, "Admin not found"
	}

	// Check if admin is active
	if admin.Status != AdminStatusActive {
		return false, "Admin is not active"
	}

	// Super admin has all permissions
	if admin.Role == AdminRoleSuper {
		return true, "Super admin has full permissions"
	}

	// Check container permission
	for _, perm := range admin.Permissions {
		if perm.ContainerID == containerID {
			switch action {
			case "manage":
				if perm.CanManage {
					return true, "Permission granted"
				}
			case "read":
				if perm.CanRead {
					return true, "Permission granted"
				}
			case "write":
				if perm.CanWrite {
					return true, "Permission granted"
				}
			}
			return false, "Permission denied for this action"
		}
	}

	return false, "No permission for this container"
}
