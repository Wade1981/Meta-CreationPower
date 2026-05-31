package database

import (
	"database/sql"
	"time"
)

type Admin struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string
	Email        string
	GroupID      string
	TokenID      string
	CreatedAt    time.Time
	LastLoginAt  time.Time
	IsActive     bool
	Metadata     map[string]interface{}
}

type AdminGroup struct {
	ID          string
	Name        string
	Description string
	Permissions map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AdminDAO struct {
	*BaseDAO
}

func NewAdminDAO(db *sql.DB) *AdminDAO {
	return &AdminDAO{BaseDAO: NewBaseDAO(db)}
}

func (dao *AdminDAO) Create(admin *Admin) error {
	query := `INSERT INTO admins (
		id, username, password_hash, role, email, group_id, token_id, created_at, last_login_at, is_active, metadata
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	metadataJSON := toJSON(admin.Metadata)
	
	_, err := dao.Exec(query,
		admin.ID,
		admin.Username,
		admin.PasswordHash,
		admin.Role,
		admin.Email,
		admin.GroupID,
		admin.TokenID,
		formatTime(admin.CreatedAt),
		nullTime(admin.LastLoginAt),
		admin.IsActive,
		metadataJSON,
	)
	return err
}

func (dao *AdminDAO) GetByID(id string) (*Admin, error) {
	query := `SELECT id, username, password_hash, role, email, group_id, token_id, created_at, last_login_at, is_active, metadata
		FROM admins WHERE id = ?`
	
	row := dao.QueryRow(query, id)
	
	admin := &Admin{}
	var createdAt, lastLoginAt sql.NullString
	var email, groupID, tokenID, metadataJSON sql.NullString
	var isActive sql.NullInt64
	
	err := row.Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.Role,
		&email,
		&groupID,
		&tokenID,
		&createdAt,
		&lastLoginAt,
		&isActive,
		&metadataJSON,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	admin.Email = fromNullString(email)
	admin.GroupID = fromNullString(groupID)
	admin.TokenID = fromNullString(tokenID)
	admin.CreatedAt = parseTime(createdAt.String)
	admin.LastLoginAt = fromNullTime(lastLoginAt)
	admin.IsActive = fromNullInt64(isActive) == 1
	
	if metadataJSON.Valid {
		fromJSON(metadataJSON.String, &admin.Metadata)
	}
	
	return admin, nil
}

func (dao *AdminDAO) GetByUsername(username string) (*Admin, error) {
	query := `SELECT id, username, password_hash, role, email, group_id, token_id, created_at, last_login_at, is_active, metadata
		FROM admins WHERE username = ?`
	
	row := dao.QueryRow(query, username)
	
	admin := &Admin{}
	var createdAt, lastLoginAt sql.NullString
	var email, groupID, tokenID, metadataJSON sql.NullString
	var isActive sql.NullInt64
	
	err := row.Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.Role,
		&email,
		&groupID,
		&tokenID,
		&createdAt,
		&lastLoginAt,
		&isActive,
		&metadataJSON,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	admin.Email = fromNullString(email)
	admin.GroupID = fromNullString(groupID)
	admin.TokenID = fromNullString(tokenID)
	admin.CreatedAt = parseTime(createdAt.String)
	admin.LastLoginAt = fromNullTime(lastLoginAt)
	admin.IsActive = fromNullInt64(isActive) == 1
	
	if metadataJSON.Valid {
		fromJSON(metadataJSON.String, &admin.Metadata)
	}
	
	return admin, nil
}

func (dao *AdminDAO) Update(admin *Admin) error {
	query := `UPDATE admins SET
		username = ?, password_hash = ?, role = ?, email = ?, group_id = ?, token_id = ?, last_login_at = ?, is_active = ?, metadata = ?
		WHERE id = ?`
	
	metadataJSON := toJSON(admin.Metadata)
	
	_, err := dao.Exec(query,
		admin.Username,
		admin.PasswordHash,
		admin.Role,
		admin.Email,
		admin.GroupID,
		admin.TokenID,
		nullTime(admin.LastLoginAt),
		admin.IsActive,
		metadataJSON,
		admin.ID,
	)
	return err
}

func (dao *AdminDAO) UpdateLastLogin(id string) error {
	query := `UPDATE admins SET last_login_at = ? WHERE id = ?`
	_, err := dao.Exec(query, formatTime(time.Now()), id)
	return err
}

func (dao *AdminDAO) UpdatePassword(id string, passwordHash string) error {
	query := `UPDATE admins SET password_hash = ? WHERE id = ?`
	_, err := dao.Exec(query, passwordHash, id)
	return err
}

func (dao *AdminDAO) UpdateRole(id string, role string) error {
	query := `UPDATE admins SET role = ? WHERE id = ?`
	_, err := dao.Exec(query, role, id)
	return err
}

func (dao *AdminDAO) UpdateToken(id string, tokenID string) error {
	query := `UPDATE admins SET token_id = ? WHERE id = ?`
	_, err := dao.Exec(query, tokenID, id)
	return err
}

func (dao *AdminDAO) SetActive(id string, isActive bool) error {
	query := `UPDATE admins SET is_active = ? WHERE id = ?`
	_, err := dao.Exec(query, isActive, id)
	return err
}

func (dao *AdminDAO) Delete(id string) error {
	query := `DELETE FROM admins WHERE id = ?`
	_, err := dao.Exec(query, id)
	return err
}

func (dao *AdminDAO) List(limit, offset int) ([]*Admin, error) {
	query := `SELECT id, username, password_hash, role, email, group_id, token_id, created_at, last_login_at, is_active, metadata
		FROM admins ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := dao.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var admins []*Admin
	for rows.Next() {
		admin := &Admin{}
		var createdAt, lastLoginAt sql.NullString
		var email, groupID, tokenID, metadataJSON sql.NullString
		var isActive sql.NullInt64
		
		err := rows.Scan(
			&admin.ID,
			&admin.Username,
			&admin.PasswordHash,
			&admin.Role,
			&email,
			&groupID,
			&tokenID,
			&createdAt,
			&lastLoginAt,
			&isActive,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}
		
		admin.Email = fromNullString(email)
		admin.GroupID = fromNullString(groupID)
		admin.TokenID = fromNullString(tokenID)
		admin.CreatedAt = parseTime(createdAt.String)
		admin.LastLoginAt = fromNullTime(lastLoginAt)
		admin.IsActive = fromNullInt64(isActive) == 1
		
		if metadataJSON.Valid {
			fromJSON(metadataJSON.String, &admin.Metadata)
		}
		
		admins = append(admins, admin)
	}
	
	return admins, nil
}

func (dao *AdminDAO) ListByRole(role string) ([]*Admin, error) {
	query := `SELECT id, username, password_hash, role, email, group_id, token_id, created_at, last_login_at, is_active, metadata
		FROM admins WHERE role = ? ORDER BY created_at DESC`
	
	rows, err := dao.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var admins []*Admin
	for rows.Next() {
		admin := &Admin{}
		var createdAt, lastLoginAt sql.NullString
		var email, groupID, tokenID, metadataJSON sql.NullString
		var isActive sql.NullInt64
		
		err := rows.Scan(
			&admin.ID,
			&admin.Username,
			&admin.PasswordHash,
			&admin.Role,
			&email,
			&groupID,
			&tokenID,
			&createdAt,
			&lastLoginAt,
			&isActive,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}
		
		admin.Email = fromNullString(email)
		admin.GroupID = fromNullString(groupID)
		admin.TokenID = fromNullString(tokenID)
		admin.CreatedAt = parseTime(createdAt.String)
		admin.LastLoginAt = fromNullTime(lastLoginAt)
		admin.IsActive = fromNullInt64(isActive) == 1
		
		if metadataJSON.Valid {
			fromJSON(metadataJSON.String, &admin.Metadata)
		}
		
		admins = append(admins, admin)
	}
	
	return admins, nil
}

func (dao *AdminDAO) ListByGroup(groupID string) ([]*Admin, error) {
	query := `SELECT id, username, password_hash, role, email, group_id, token_id, created_at, last_login_at, is_active, metadata
		FROM admins WHERE group_id = ? ORDER BY created_at DESC`
	
	rows, err := dao.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var admins []*Admin
	for rows.Next() {
		admin := &Admin{}
		var createdAt, lastLoginAt sql.NullString
		var email, groupIDVal, tokenID, metadataJSON sql.NullString
		var isActive sql.NullInt64
		
		err := rows.Scan(
			&admin.ID,
			&admin.Username,
			&admin.PasswordHash,
			&admin.Role,
			&email,
			&groupIDVal,
			&tokenID,
			&createdAt,
			&lastLoginAt,
			&isActive,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}
		
		admin.Email = fromNullString(email)
		admin.GroupID = fromNullString(groupIDVal)
		admin.TokenID = fromNullString(tokenID)
		admin.CreatedAt = parseTime(createdAt.String)
		admin.LastLoginAt = fromNullTime(lastLoginAt)
		admin.IsActive = fromNullInt64(isActive) == 1
		
		if metadataJSON.Valid {
			fromJSON(metadataJSON.String, &admin.Metadata)
		}
		
		admins = append(admins, admin)
	}
	
	return admins, nil
}

func (dao *AdminDAO) ListSuperAdmins() ([]*Admin, error) {
	return dao.ListByRole("super_admin")
}

func (dao *AdminDAO) ListActive() ([]*Admin, error) {
	query := `SELECT id, username, password_hash, role, email, group_id, token_id, created_at, last_login_at, is_active, metadata
		FROM admins WHERE is_active = 1 ORDER BY created_at DESC`
	
	rows, err := dao.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var admins []*Admin
	for rows.Next() {
		admin := &Admin{}
		var createdAt, lastLoginAt sql.NullString
		var email, groupID, tokenID, metadataJSON sql.NullString
		var isActive sql.NullInt64
		
		err := rows.Scan(
			&admin.ID,
			&admin.Username,
			&admin.PasswordHash,
			&admin.Role,
			&email,
			&groupID,
			&tokenID,
			&createdAt,
			&lastLoginAt,
			&isActive,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}
		
		admin.Email = fromNullString(email)
		admin.GroupID = fromNullString(groupID)
		admin.TokenID = fromNullString(tokenID)
		admin.CreatedAt = parseTime(createdAt.String)
		admin.LastLoginAt = fromNullTime(lastLoginAt)
		admin.IsActive = fromNullInt64(isActive) == 1
		
		if metadataJSON.Valid {
			fromJSON(metadataJSON.String, &admin.Metadata)
		}
		
		admins = append(admins, admin)
	}
	
	return admins, nil
}

func (dao *AdminDAO) ValidatePassword(username, passwordHash string) (bool, *Admin, error) {
	admin, err := dao.GetByUsername(username)
	if err != nil {
		return false, nil, err
	}
	
	if admin == nil {
		return false, nil, nil
	}
	
	if !admin.IsActive {
		return false, admin, nil
	}
	
	if admin.PasswordHash != passwordHash {
		return false, admin, nil
	}
	
	return true, admin, nil
}

func (dao *AdminDAO) Count() (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM admins").Scan(&count)
	return count, err
}

func (dao *AdminDAO) CountByRole(role string) (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM admins WHERE role = ?", role).Scan(&count)
	return count, err
}

func (dao *AdminDAO) Exists(id string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM admins WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *AdminDAO) ExistsByUsername(username string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM admins WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *AdminDAO) CreateGroup(group *AdminGroup) error {
	query := `INSERT INTO admin_groups (id, name, description, permissions, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	
	permissionsJSON := toJSON(group.Permissions)
	
	_, err := dao.Exec(query,
		group.ID,
		group.Name,
		group.Description,
		permissionsJSON,
		formatTime(group.CreatedAt),
		nullTime(group.UpdatedAt),
	)
	return err
}

func (dao *AdminDAO) GetGroupByID(id string) (*AdminGroup, error) {
	query := `SELECT id, name, description, permissions, created_at, updated_at FROM admin_groups WHERE id = ?`
	
	row := dao.QueryRow(query, id)
	
	group := &AdminGroup{}
	var createdAt, updatedAt sql.NullString
	var description, permissionsJSON sql.NullString
	
	err := row.Scan(
		&group.ID,
		&group.Name,
		&description,
		&permissionsJSON,
		&createdAt,
		&updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	group.Description = fromNullString(description)
	group.CreatedAt = parseTime(createdAt.String)
	group.UpdatedAt = fromNullTime(updatedAt)
	
	if permissionsJSON.Valid {
		fromJSON(permissionsJSON.String, &group.Permissions)
	}
	
	return group, nil
}

func (dao *AdminDAO) GetGroupByName(name string) (*AdminGroup, error) {
	query := `SELECT id, name, description, permissions, created_at, updated_at FROM admin_groups WHERE name = ?`
	
	row := dao.QueryRow(query, name)
	
	group := &AdminGroup{}
	var createdAt, updatedAt sql.NullString
	var description, permissionsJSON sql.NullString
	
	err := row.Scan(
		&group.ID,
		&group.Name,
		&description,
		&permissionsJSON,
		&createdAt,
		&updatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	group.Description = fromNullString(description)
	group.CreatedAt = parseTime(createdAt.String)
	group.UpdatedAt = fromNullTime(updatedAt)
	
	if permissionsJSON.Valid {
		fromJSON(permissionsJSON.String, &group.Permissions)
	}
	
	return group, nil
}

func (dao *AdminDAO) UpdateGroup(group *AdminGroup) error {
	query := `UPDATE admin_groups SET name = ?, description = ?, permissions = ?, updated_at = ? WHERE id = ?`
	
	permissionsJSON := toJSON(group.Permissions)
	
	_, err := dao.Exec(query,
		group.Name,
		group.Description,
		permissionsJSON,
		formatTime(time.Now()),
		group.ID,
	)
	return err
}

func (dao *AdminDAO) DeleteGroup(id string) error {
	query := `DELETE FROM admin_groups WHERE id = ?`
	_, err := dao.Exec(query, id)
	return err
}

func (dao *AdminDAO) ListGroups() ([]*AdminGroup, error) {
	query := `SELECT id, name, description, permissions, created_at, updated_at FROM admin_groups ORDER BY created_at DESC`
	
	rows, err := dao.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var groups []*AdminGroup
	for rows.Next() {
		group := &AdminGroup{}
		var createdAt, updatedAt sql.NullString
		var description, permissionsJSON sql.NullString
		
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&description,
			&permissionsJSON,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		group.Description = fromNullString(description)
		group.CreatedAt = parseTime(createdAt.String)
		group.UpdatedAt = fromNullTime(updatedAt)
		
		if permissionsJSON.Valid {
			fromJSON(permissionsJSON.String, &group.Permissions)
		}
		
		groups = append(groups, group)
	}
	
	return groups, nil
}

func (dao *AdminDAO) CountGroups() (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM admin_groups").Scan(&count)
	return count, err
}

func (dao *AdminDAO) ExistsGroup(id string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM admin_groups WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *AdminDAO) ExistsGroupByName(name string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM admin_groups WHERE name = ?", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}