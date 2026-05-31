package database

import (
	"crypto/sha256"
	"encoding/hex"
	"database/sql"
	"time"
)

type Token struct {
	ID              string
	UserID          string
	AdminID         string
	TokenHash       string
	TokenType       string
	EndpointType    string
	Permissions     map[string]interface{}
	ExpiresAt       time.Time
	IssuedAt        time.Time
	LastUsedAt      time.Time
	Status          string
	RefreshTokenHash string
	UsageCount      int64
	MaxUsage        int64
	AllowedIPs      []string
	Metadata        map[string]interface{}
}

type TokenDAO struct {
	*BaseDAO
}

func NewTokenDAO(db *sql.DB) *TokenDAO {
	return &TokenDAO{BaseDAO: NewBaseDAO(db)}
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (dao *TokenDAO) Create(token *Token) error {
	query := `INSERT INTO tokens (
		id, user_id, admin_id, token_hash, token_type, endpoint_type, permissions,
		expires_at, issued_at, last_used_at, status, refresh_token_hash, usage_count, max_usage,
		allowed_ips, metadata
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	permissionsJSON := toJSON(token.Permissions)
	allowedIPsJSON := toJSON(token.AllowedIPs)
	metadataJSON := toJSON(token.Metadata)
	
	_, err := dao.Exec(query,
		token.ID,
		token.UserID,
		token.AdminID,
		token.TokenHash,
		token.TokenType,
		token.EndpointType,
		permissionsJSON,
		nullTime(token.ExpiresAt),
		formatTime(token.IssuedAt),
		nullTime(token.LastUsedAt),
		token.Status,
		token.RefreshTokenHash,
		token.UsageCount,
		token.MaxUsage,
		allowedIPsJSON,
		metadataJSON,
	)
	return err
}

func (dao *TokenDAO) GetByID(id string) (*Token, error) {
	query := `SELECT id, user_id, admin_id, token_hash, token_type, endpoint_type, permissions,
		expires_at, issued_at, last_used_at, status, refresh_token_hash, usage_count, max_usage,
		allowed_ips, metadata
		FROM tokens WHERE id = ?`
	
	row := dao.QueryRow(query, id)
	
	token := &Token{}
	var userID, adminID, permissionsJSON, allowedIPsJSON, metadataJSON sql.NullString
	var expiresAt, lastUsedAt sql.NullString
	var refreshTokenHash sql.NullString
	var usageCount, maxUsage sql.NullInt64
	
	err := row.Scan(
		&token.ID,
		&userID,
		&adminID,
		&token.TokenHash,
		&token.TokenType,
		&token.EndpointType,
		&permissionsJSON,
		&expiresAt,
		&token.IssuedAt,
		&lastUsedAt,
		&token.Status,
		&refreshTokenHash,
		&usageCount,
		&maxUsage,
		&allowedIPsJSON,
		&metadataJSON,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	token.UserID = fromNullString(userID)
	token.AdminID = fromNullString(adminID)
	token.ExpiresAt = fromNullTime(expiresAt)
	token.LastUsedAt = fromNullTime(lastUsedAt)
	token.RefreshTokenHash = fromNullString(refreshTokenHash)
	token.UsageCount = fromNullInt64(usageCount)
	token.MaxUsage = fromNullInt64(maxUsage)
	
	if permissionsJSON.Valid {
		fromJSON(permissionsJSON.String, &token.Permissions)
	}
	if allowedIPsJSON.Valid {
		fromJSON(allowedIPsJSON.String, &token.AllowedIPs)
	}
	if metadataJSON.Valid {
		fromJSON(metadataJSON.String, &token.Metadata)
	}
	
	return token, nil
}

func (dao *TokenDAO) GetByHash(tokenHash string) (*Token, error) {
	query := `SELECT id, user_id, admin_id, token_hash, token_type, endpoint_type, permissions,
		expires_at, issued_at, last_used_at, status, refresh_token_hash, usage_count, max_usage,
		allowed_ips, metadata
		FROM tokens WHERE token_hash = ?`
	
	row := dao.QueryRow(query, tokenHash)
	
	token := &Token{}
	var userID, adminID, permissionsJSON, allowedIPsJSON, metadataJSON sql.NullString
	var expiresAt, lastUsedAt sql.NullString
	var refreshTokenHash sql.NullString
	var usageCount, maxUsage sql.NullInt64
	
	err := row.Scan(
		&token.ID,
		&userID,
		&adminID,
		&token.TokenHash,
		&token.TokenType,
		&token.EndpointType,
		&permissionsJSON,
		&expiresAt,
		&token.IssuedAt,
		&lastUsedAt,
		&token.Status,
		&refreshTokenHash,
		&usageCount,
		&maxUsage,
		&allowedIPsJSON,
		&metadataJSON,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	token.UserID = fromNullString(userID)
	token.AdminID = fromNullString(adminID)
	token.ExpiresAt = fromNullTime(expiresAt)
	token.LastUsedAt = fromNullTime(lastUsedAt)
	token.RefreshTokenHash = fromNullString(refreshTokenHash)
	token.UsageCount = fromNullInt64(usageCount)
	token.MaxUsage = fromNullInt64(maxUsage)
	
	if permissionsJSON.Valid {
		fromJSON(permissionsJSON.String, &token.Permissions)
	}
	if allowedIPsJSON.Valid {
		fromJSON(allowedIPsJSON.String, &token.AllowedIPs)
	}
	if metadataJSON.Valid {
		fromJSON(metadataJSON.String, &token.Metadata)
	}
	
	return token, nil
}

func (dao *TokenDAO) Update(token *Token) error {
	query := `UPDATE tokens SET
		user_id = ?, admin_id = ?, token_type = ?, endpoint_type = ?, permissions = ?,
		expires_at = ?, last_used_at = ?, status = ?, refresh_token_hash = ?, usage_count = ?, max_usage = ?,
		allowed_ips = ?, metadata = ?
		WHERE id = ?`
	
	permissionsJSON := toJSON(token.Permissions)
	allowedIPsJSON := toJSON(token.AllowedIPs)
	metadataJSON := toJSON(token.Metadata)
	
	_, err := dao.Exec(query,
		token.UserID,
		token.AdminID,
		token.TokenType,
		token.EndpointType,
		permissionsJSON,
		nullTime(token.ExpiresAt),
		nullTime(token.LastUsedAt),
		token.Status,
		token.RefreshTokenHash,
		token.UsageCount,
		token.MaxUsage,
		allowedIPsJSON,
		metadataJSON,
		token.ID,
	)
	return err
}

func (dao *TokenDAO) UpdateStatus(id string, status string) error {
	query := `UPDATE tokens SET status = ? WHERE id = ?`
	_, err := dao.Exec(query, status, id)
	return err
}

func (dao *TokenDAO) UpdateLastUsed(id string) error {
	query := `UPDATE tokens SET last_used_at = ?, usage_count = usage_count + 1 WHERE id = ?`
	_, err := dao.Exec(query, formatTime(time.Now()), id)
	return err
}

func (dao *TokenDAO) Revoke(id string) error {
	return dao.UpdateStatus(id, "revoked")
}

func (dao *TokenDAO) Delete(id string) error {
	query := `DELETE FROM tokens WHERE id = ?`
	_, err := dao.Exec(query, id)
	return err
}

func (dao *TokenDAO) List(limit, offset int) ([]*Token, error) {
	query := `SELECT id, user_id, admin_id, token_hash, token_type, endpoint_type, permissions,
		expires_at, issued_at, last_used_at, status, refresh_token_hash, usage_count, max_usage,
		allowed_ips, metadata
		FROM tokens ORDER BY issued_at DESC LIMIT ? OFFSET ?`
	
	rows, err := dao.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tokens []*Token
	for rows.Next() {
		token := &Token{}
		var userID, adminID, permissionsJSON, allowedIPsJSON, metadataJSON sql.NullString
		var expiresAt, lastUsedAt sql.NullString
		var refreshTokenHash sql.NullString
		var usageCount, maxUsage sql.NullInt64
		
		err := rows.Scan(
			&token.ID,
			&userID,
			&adminID,
			&token.TokenHash,
			&token.TokenType,
			&token.EndpointType,
			&permissionsJSON,
			&expiresAt,
			&token.IssuedAt,
			&lastUsedAt,
			&token.Status,
			&refreshTokenHash,
			&usageCount,
			&maxUsage,
			&allowedIPsJSON,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}
		
		token.UserID = fromNullString(userID)
		token.AdminID = fromNullString(adminID)
		token.ExpiresAt = fromNullTime(expiresAt)
		token.LastUsedAt = fromNullTime(lastUsedAt)
		token.RefreshTokenHash = fromNullString(refreshTokenHash)
		token.UsageCount = fromNullInt64(usageCount)
		token.MaxUsage = fromNullInt64(maxUsage)
		
		if permissionsJSON.Valid {
			fromJSON(permissionsJSON.String, &token.Permissions)
		}
		if allowedIPsJSON.Valid {
			fromJSON(allowedIPsJSON.String, &token.AllowedIPs)
		}
		if metadataJSON.Valid {
			fromJSON(metadataJSON.String, &token.Metadata)
		}
		
		tokens = append(tokens, token)
	}
	
	return tokens, nil
}

func (dao *TokenDAO) ListByStatus(status string) ([]*Token, error) {
	query := `SELECT id, user_id, admin_id, token_hash, token_type, endpoint_type, permissions,
		expires_at, issued_at, last_used_at, status, refresh_token_hash, usage_count, max_usage,
		allowed_ips, metadata
		FROM tokens WHERE status = ? ORDER BY issued_at DESC`
	
	rows, err := dao.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tokens []*Token
	for rows.Next() {
		token := &Token{}
		var userID, adminID, permissionsJSON, allowedIPsJSON, metadataJSON sql.NullString
		var expiresAt, lastUsedAt sql.NullString
		var refreshTokenHash sql.NullString
		var usageCount, maxUsage sql.NullInt64
		
		err := rows.Scan(
			&token.ID,
			&userID,
			&adminID,
			&token.TokenHash,
			&token.TokenType,
			&token.EndpointType,
			&permissionsJSON,
			&expiresAt,
			&token.IssuedAt,
			&lastUsedAt,
			&token.Status,
			&refreshTokenHash,
			&usageCount,
			&maxUsage,
			&allowedIPsJSON,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}
		
		token.UserID = fromNullString(userID)
		token.AdminID = fromNullString(adminID)
		token.ExpiresAt = fromNullTime(expiresAt)
		token.LastUsedAt = fromNullTime(lastUsedAt)
		token.RefreshTokenHash = fromNullString(refreshTokenHash)
		token.UsageCount = fromNullInt64(usageCount)
		token.MaxUsage = fromNullInt64(maxUsage)
		
		if permissionsJSON.Valid {
			fromJSON(permissionsJSON.String, &token.Permissions)
		}
		if allowedIPsJSON.Valid {
			fromJSON(allowedIPsJSON.String, &token.AllowedIPs)
		}
		if metadataJSON.Valid {
			fromJSON(metadataJSON.String, &token.Metadata)
		}
		
		tokens = append(tokens, token)
	}
	
	return tokens, nil
}

func (dao *TokenDAO) ListActive() ([]*Token, error) {
	return dao.ListByStatus("active")
}

func (dao *TokenDAO) ListByType(tokenType string) ([]*Token, error) {
	query := `SELECT id, user_id, admin_id, token_hash, token_type, endpoint_type, permissions,
		expires_at, issued_at, last_used_at, status, refresh_token_hash, usage_count, max_usage,
		allowed_ips, metadata
		FROM tokens WHERE token_type = ? ORDER BY issued_at DESC`
	
	rows, err := dao.Query(query, tokenType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tokens []*Token
	for rows.Next() {
		token := &Token{}
		var userID, adminID, permissionsJSON, allowedIPsJSON, metadataJSON sql.NullString
		var expiresAt, lastUsedAt sql.NullString
		var refreshTokenHash sql.NullString
		var usageCount, maxUsage sql.NullInt64
		
		err := rows.Scan(
			&token.ID,
			&userID,
			&adminID,
			&token.TokenHash,
			&token.TokenType,
			&token.EndpointType,
			&permissionsJSON,
			&expiresAt,
			&token.IssuedAt,
			&lastUsedAt,
			&token.Status,
			&refreshTokenHash,
			&usageCount,
			&maxUsage,
			&allowedIPsJSON,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}
		
		token.UserID = fromNullString(userID)
		token.AdminID = fromNullString(adminID)
		token.ExpiresAt = fromNullTime(expiresAt)
		token.LastUsedAt = fromNullTime(lastUsedAt)
		token.RefreshTokenHash = fromNullString(refreshTokenHash)
		token.UsageCount = fromNullInt64(usageCount)
		token.MaxUsage = fromNullInt64(maxUsage)
		
		if permissionsJSON.Valid {
			fromJSON(permissionsJSON.String, &token.Permissions)
		}
		if allowedIPsJSON.Valid {
			fromJSON(allowedIPsJSON.String, &token.AllowedIPs)
		}
		if metadataJSON.Valid {
			fromJSON(metadataJSON.String, &token.Metadata)
		}
		
		tokens = append(tokens, token)
	}
	
	return tokens, nil
}

func (dao *TokenDAO) Validate(tokenHash string) (bool, *Token, error) {
	token, err := dao.GetByHash(tokenHash)
	if err != nil {
		return false, nil, err
	}
	
	if token == nil {
		return false, nil, nil
	}
	
	if token.Status != "active" {
		return false, token, nil
	}
	
	if !token.ExpiresAt.IsZero() && token.ExpiresAt.Before(time.Now()) {
		dao.UpdateStatus(token.ID, "expired")
		return false, token, nil
	}
	
	if token.MaxUsage > 0 && token.UsageCount >= token.MaxUsage {
		dao.UpdateStatus(token.ID, "expired")
		return false, token, nil
	}
	
	return true, token, nil
}

func (dao *TokenDAO) Count() (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM tokens").Scan(&count)
	return count, err
}

func (dao *TokenDAO) CountByStatus(status string) (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM tokens WHERE status = ?", status).Scan(&count)
	return count, err
}

func (dao *TokenDAO) Exists(id string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM tokens WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *TokenDAO) ExistsByHash(tokenHash string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM tokens WHERE token_hash = ?", tokenHash).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *TokenDAO) CreateStats(tokenID string, date string, requestCount, successCount, failureCount int, responseTime float64) error {
	query := `INSERT INTO token_stats (token_id, date, request_count, success_count, failure_count, total_response_time, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := dao.Exec(query, tokenID, date, requestCount, successCount, failureCount, responseTime, formatTime(time.Now()))
	return err
}

func (dao *TokenDAO) UpdateStats(tokenID string, date string, requestCount, successCount, failureCount int, responseTime float64) error {
	query := `UPDATE token_stats SET
		request_count = request_count + ?, success_count = success_count + ?, failure_count = failure_count + ?,
		total_response_time = total_response_time + ?
		WHERE token_id = ? AND date = ?`
	_, err := dao.Exec(query, requestCount, successCount, failureCount, responseTime, tokenID, date)
	return err
}

func (dao *TokenDAO) GetStats(tokenID string, date string) (int, int, int, float64, error) {
	query := `SELECT request_count, success_count, failure_count, total_response_time FROM token_stats WHERE token_id = ? AND date = ?`
	
	var requestCount, successCount, failureCount sql.NullInt64
	var totalResponseTime sql.NullFloat64
	
	err := dao.QueryRow(query, tokenID, date).Scan(&requestCount, &successCount, &failureCount, &totalResponseTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, 0, 0, nil
		}
		return 0, 0, 0, 0, err
	}
	
	return int(fromNullInt64(requestCount)), int(fromNullInt64(successCount)), int(fromNullInt64(failureCount)), fromNullFloat64(totalResponseTime), nil
}