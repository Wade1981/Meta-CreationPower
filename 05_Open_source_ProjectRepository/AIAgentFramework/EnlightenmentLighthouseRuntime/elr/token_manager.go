// Package elr implements token management for Enlightenment Lighthouse Runtime
package elr

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// TokenStatus represents the status of a token
type TokenStatus string

const (
	// TokenStatusActive represents an active token
	TokenStatusActive TokenStatus = "active"
	// TokenStatusRevoked represents a revoked token
	TokenStatusRevoked TokenStatus = "revoked"
	// TokenStatusExpired represents an expired token
	TokenStatusExpired TokenStatus = "expired"
)

// Token represents a token instance
type Token struct {
	ID          string       `json:"id"`
	Secret      string       `json:"secret"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	ExpiresAt   time.Time    `json:"expires_at"`
	Status      TokenStatus  `json:"status"`
}

// TokenManager represents the token manager
type TokenManager struct {
	tokenFile string
	tokens    []Token
}

// NewTokenManager creates a new token manager
func NewTokenManager(tokenFile string) *TokenManager {
	return &TokenManager{
		tokenFile: tokenFile,
		tokens:    []Token{},
	}
}

// LoadTokens loads tokens from file
func (tm *TokenManager) LoadTokens() error {
	if _, err := os.Stat(tm.tokenFile); os.IsNotExist(err) {
		// Token file doesn't exist, initialize with empty tokens
		tm.tokens = []Token{}
		return nil
	}

	data, err := os.ReadFile(tm.tokenFile)
	if err != nil {
		return fmt.Errorf("failed to read token file: %w", err)
	}

	var tokenData struct {
		Tokens       []Token `json:"tokens"`
		LastUpdated  int64   `json:"last_updated"`
	}

	if err := json.Unmarshal(data, &tokenData); err != nil {
		return fmt.Errorf("failed to unmarshal token data: %w", err)
	}

	tm.tokens = tokenData.Tokens
	return nil
}

// SaveTokens saves tokens to file
func (tm *TokenManager) SaveTokens() error {
	tokenData := struct {
		Tokens      []Token `json:"tokens"`
		LastUpdated int64   `json:"last_updated"`
	}{
		Tokens:      tm.tokens,
		LastUpdated: time.Now().Unix(),
	}

	data, err := json.MarshalIndent(tokenData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(tm.tokenFile), 0755); err != nil {
		return fmt.Errorf("failed to create token directory: %w", err)
	}

	if err := os.WriteFile(tm.tokenFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// GenerateToken generates a new token
func (tm *TokenManager) GenerateToken(description string) (string, error) {
	tokenID := uuid.New().String()
	tokenSecret := fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.New().String())))

	token := Token{
		ID:          tokenID,
		Secret:      tokenSecret,
		Description: description,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour), // 7 days expiration
		Status:      TokenStatusActive,
	}

	tm.tokens = append(tm.tokens, token)

	if err := tm.SaveTokens(); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}

	return fmt.Sprintf("%s.%s", tokenID, tokenSecret), nil
}

// ValidateToken validates a token
func (tm *TokenManager) ValidateToken(tokenString string) (bool, string) {
	if tokenString == "" {
		return false, "Token is required"
	}

	// Parse token string
	parts := splitToken(tokenString)
	if len(parts) != 2 {
		return false, "Invalid token format"
	}

	tokenID, tokenSecret := parts[0], parts[1]

	// Find token
	for _, token := range tm.tokens {
		if token.ID == tokenID && token.Secret == tokenSecret {
			// Check if token is expired
			if time.Now().After(token.ExpiresAt) {
				return false, "Token has expired"
			}

			// Check if token is active
			if token.Status != TokenStatusActive {
				return false, "Token is not active"
			}

			return true, "Token is valid"
		}
	}

	return false, "Token not found"
}

// RefreshToken refreshes a token
func (tm *TokenManager) RefreshToken(oldTokenString, description string) (string, error) {
	// Validate old token
	valid, message := tm.ValidateToken(oldTokenString)
	if !valid {
		return "", fmt.Errorf("invalid old token: %s", message)
	}

	// Parse old token
	parts := splitToken(oldTokenString)
	tokenID := parts[0]

	// Revoke old token
	for i, token := range tm.tokens {
		if token.ID == tokenID {
			tm.tokens[i].Status = TokenStatusRevoked
			break
		}
	}

	// Generate new token
	newToken, err := tm.GenerateToken(description)
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	return newToken, nil
}

// ListTokens lists all tokens
func (tm *TokenManager) ListTokens() []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(tm.tokens))

	for _, token := range tm.tokens {
		tokenInfo := map[string]interface{}{
			"id":          token.ID,
			"description": token.Description,
			"created_at":  token.CreatedAt.Unix(),
			"expires_at":  token.ExpiresAt.Unix(),
			"status":      token.Status,
			"expired":     time.Now().After(token.ExpiresAt),
		}
		result = append(result, tokenInfo)
	}

	return result
}

// RevokeToken revokes a token
func (tm *TokenManager) RevokeToken(tokenID string) error {
	for i, token := range tm.tokens {
		if token.ID == tokenID {
			tm.tokens[i].Status = TokenStatusRevoked
			if err := tm.SaveTokens(); err != nil {
				return fmt.Errorf("failed to save token: %w", err)
			}
			return nil
		}
	}

	return fmt.Errorf("token not found")
}

// splitToken splits a token string into ID and secret
func splitToken(token string) []string {
	parts := make([]string, 0, 2)
	for i, c := range token {
		if c == '.' {
			parts = append(parts, token[:i], token[i+1:])
			break
		}
	}
	return parts
}
