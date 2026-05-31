package database

import (
	"database/sql"
	"time"
)

type Sandbox struct {
	ID          string
	Name        string
	ContainerID string
	Status      string
	ProcessID   int64
	PipeName    string
	CreatedAt   time.Time
	StartedAt   time.Time
	Uptime      int64
	Config      map[string]interface{}
}

type SandboxDAO struct {
	*BaseDAO
}

func NewSandboxDAO(db *sql.DB) *SandboxDAO {
	return &SandboxDAO{BaseDAO: NewBaseDAO(db)}
}

func (dao *SandboxDAO) Create(sandbox *Sandbox) error {
	query := `INSERT INTO sandboxes (
		id, name, container_id, status, process_id, pipe_name, created_at, started_at, uptime, config
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	configJSON := toJSON(sandbox.Config)
	
	_, err := dao.Exec(query,
		sandbox.ID,
		sandbox.Name,
		sandbox.ContainerID,
		sandbox.Status,
		sandbox.ProcessID,
		sandbox.PipeName,
		formatTime(sandbox.CreatedAt),
		nullTime(sandbox.StartedAt),
		sandbox.Uptime,
		configJSON,
	)
	return err
}

func (dao *SandboxDAO) GetByID(id string) (*Sandbox, error) {
	query := `SELECT id, name, container_id, status, process_id, pipe_name, created_at, started_at, uptime, config
		FROM sandboxes WHERE id = ?`
	
	row := dao.QueryRow(query, id)
	
	sandbox := &Sandbox{}
	var createdAt, startedAt sql.NullString
	var containerID, pipeName, configJSON sql.NullString
	var processID, uptime sql.NullInt64
	
	err := row.Scan(
		&sandbox.ID,
		&sandbox.Name,
		&containerID,
		&sandbox.Status,
		&processID,
		&pipeName,
		&createdAt,
		&startedAt,
		&uptime,
		&configJSON,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	sandbox.ContainerID = fromNullString(containerID)
	sandbox.ProcessID = fromNullInt64(processID)
	sandbox.PipeName = fromNullString(pipeName)
	sandbox.CreatedAt = parseTime(createdAt.String)
	sandbox.StartedAt = fromNullTime(startedAt)
	sandbox.Uptime = fromNullInt64(uptime)
	
	if configJSON.Valid {
		fromJSON(configJSON.String, &sandbox.Config)
	}
	
	return sandbox, nil
}

func (dao *SandboxDAO) Update(sandbox *Sandbox) error {
	query := `UPDATE sandboxes SET
		name = ?, container_id = ?, status = ?, process_id = ?, pipe_name = ?, started_at = ?, uptime = ?, config = ?
		WHERE id = ?`
	
	configJSON := toJSON(sandbox.Config)
	
	_, err := dao.Exec(query,
		sandbox.Name,
		sandbox.ContainerID,
		sandbox.Status,
		sandbox.ProcessID,
		sandbox.PipeName,
		nullTime(sandbox.StartedAt),
		sandbox.Uptime,
		configJSON,
		sandbox.ID,
	)
	return err
}

func (dao *SandboxDAO) UpdateStatus(id string, status string) error {
	query := `UPDATE sandboxes SET status = ? WHERE id = ?`
	_, err := dao.Exec(query, status, id)
	return err
}

func (dao *SandboxDAO) UpdateStartedAt(id string, startedAt time.Time) error {
	query := `UPDATE sandboxes SET started_at = ?, status = 'running' WHERE id = ?`
	_, err := dao.Exec(query, formatTime(startedAt), id)
	return err
}

func (dao *SandboxDAO) UpdateUptime(id string, uptime int64) error {
	query := `UPDATE sandboxes SET uptime = ? WHERE id = ?`
	_, err := dao.Exec(query, uptime, id)
	return err
}

func (dao *SandboxDAO) UpdateProcessInfo(id string, processID int64, pipeName string) error {
	query := `UPDATE sandboxes SET process_id = ?, pipe_name = ? WHERE id = ?`
	_, err := dao.Exec(query, processID, pipeName, id)
	return err
}

func (dao *SandboxDAO) Delete(id string) error {
	query := `DELETE FROM sandboxes WHERE id = ?`
	_, err := dao.Exec(query, id)
	return err
}

func (dao *SandboxDAO) List(limit, offset int) ([]*Sandbox, error) {
	query := `SELECT id, name, container_id, status, process_id, pipe_name, created_at, started_at, uptime, config
		FROM sandboxes ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := dao.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var sandboxes []*Sandbox
	for rows.Next() {
		sandbox := &Sandbox{}
		var createdAt, startedAt sql.NullString
		var containerID, pipeName, configJSON sql.NullString
		var processID, uptime sql.NullInt64
		
		err := rows.Scan(
			&sandbox.ID,
			&sandbox.Name,
			&containerID,
			&sandbox.Status,
			&processID,
			&pipeName,
			&createdAt,
			&startedAt,
			&uptime,
			&configJSON,
		)
		if err != nil {
			return nil, err
		}
		
		sandbox.ContainerID = fromNullString(containerID)
		sandbox.ProcessID = fromNullInt64(processID)
		sandbox.PipeName = fromNullString(pipeName)
		sandbox.CreatedAt = parseTime(createdAt.String)
		sandbox.StartedAt = fromNullTime(startedAt)
		sandbox.Uptime = fromNullInt64(uptime)
		
		if configJSON.Valid {
			fromJSON(configJSON.String, &sandbox.Config)
		}
		
		sandboxes = append(sandboxes, sandbox)
	}
	
	return sandboxes, nil
}

func (dao *SandboxDAO) ListByContainer(containerID string) ([]*Sandbox, error) {
	query := `SELECT id, name, container_id, status, process_id, pipe_name, created_at, started_at, uptime, config
		FROM sandboxes WHERE container_id = ? ORDER BY created_at DESC`
	
	rows, err := dao.Query(query, containerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var sandboxes []*Sandbox
	for rows.Next() {
		sandbox := &Sandbox{}
		var createdAt, startedAt sql.NullString
		var containerIDVal, pipeName, configJSON sql.NullString
		var processID, uptime sql.NullInt64
		
		err := rows.Scan(
			&sandbox.ID,
			&sandbox.Name,
			&containerIDVal,
			&sandbox.Status,
			&processID,
			&pipeName,
			&createdAt,
			&startedAt,
			&uptime,
			&configJSON,
		)
		if err != nil {
			return nil, err
		}
		
		sandbox.ContainerID = fromNullString(containerIDVal)
		sandbox.ProcessID = fromNullInt64(processID)
		sandbox.PipeName = fromNullString(pipeName)
		sandbox.CreatedAt = parseTime(createdAt.String)
		sandbox.StartedAt = fromNullTime(startedAt)
		sandbox.Uptime = fromNullInt64(uptime)
		
		if configJSON.Valid {
			fromJSON(configJSON.String, &sandbox.Config)
		}
		
		sandboxes = append(sandboxes, sandbox)
	}
	
	return sandboxes, nil
}

func (dao *SandboxDAO) ListByStatus(status string) ([]*Sandbox, error) {
	query := `SELECT id, name, container_id, status, process_id, pipe_name, created_at, started_at, uptime, config
		FROM sandboxes WHERE status = ? ORDER BY created_at DESC`
	
	rows, err := dao.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var sandboxes []*Sandbox
	for rows.Next() {
		sandbox := &Sandbox{}
		var createdAt, startedAt sql.NullString
		var containerID, pipeName, configJSON sql.NullString
		var processID, uptime sql.NullInt64
		
		err := rows.Scan(
			&sandbox.ID,
			&sandbox.Name,
			&containerID,
			&sandbox.Status,
			&processID,
			&pipeName,
			&createdAt,
			&startedAt,
			&uptime,
			&configJSON,
		)
		if err != nil {
			return nil, err
		}
		
		sandbox.ContainerID = fromNullString(containerID)
		sandbox.ProcessID = fromNullInt64(processID)
		sandbox.PipeName = fromNullString(pipeName)
		sandbox.CreatedAt = parseTime(createdAt.String)
		sandbox.StartedAt = fromNullTime(startedAt)
		sandbox.Uptime = fromNullInt64(uptime)
		
		if configJSON.Valid {
			fromJSON(configJSON.String, &sandbox.Config)
		}
		
		sandboxes = append(sandboxes, sandbox)
	}
	
	return sandboxes, nil
}

func (dao *SandboxDAO) ListRunning() ([]*Sandbox, error) {
	return dao.ListByStatus("running")
}

func (dao *SandboxDAO) Count() (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM sandboxes").Scan(&count)
	return count, err
}

func (dao *SandboxDAO) CountByContainer(containerID string) (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM sandboxes WHERE container_id = ?", containerID).Scan(&count)
	return count, err
}

func (dao *SandboxDAO) CountByStatus(status string) (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM sandboxes WHERE status = ?", status).Scan(&count)
	return count, err
}

func (dao *SandboxDAO) Exists(id string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM sandboxes WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *SandboxDAO) ExistsByName(name string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM sandboxes WHERE name = ?", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *SandboxDAO) CreateMapping(sandboxID, containerID string) error {
	query := `INSERT INTO sandbox_container_mapping (sandbox_id, container_id, created_at) VALUES (?, ?, ?)`
	_, err := dao.Exec(query, sandboxID, containerID, formatTime(time.Now()))
	return err
}

func (dao *SandboxDAO) GetMapping(sandboxID string) (string, error) {
	query := `SELECT container_id FROM sandbox_container_mapping WHERE sandbox_id = ?`
	
	var containerID sql.NullString
	err := dao.QueryRow(query, sandboxID).Scan(&containerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	
	return fromNullString(containerID), nil
}

func (dao *SandboxDAO) DeleteMapping(sandboxID string) error {
	query := `DELETE FROM sandbox_container_mapping WHERE sandbox_id = ?`
	_, err := dao.Exec(query, sandboxID)
	return err
}