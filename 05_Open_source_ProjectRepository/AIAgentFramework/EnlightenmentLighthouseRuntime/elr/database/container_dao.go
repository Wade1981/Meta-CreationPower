package database

import (
	"database/sql"
	"time"
)

type Container struct {
	ID            string
	Name          string
	Image         string
	Status        string
	IsolationType string
	CPULimit      int64
	MemoryLimit   int64
	DiskLimit     int64
	NetworkMode   string
	CreatedAt     time.Time
	StartedAt     time.Time
	StoppedAt     time.Time
	Config        map[string]interface{}
	Metadata      map[string]interface{}
}

type ContainerDAO struct {
	*BaseDAO
}

func NewContainerDAO(db *sql.DB) *ContainerDAO {
	return &ContainerDAO{BaseDAO: NewBaseDAO(db)}
}

func (dao *ContainerDAO) Create(container *Container) error {
	query := `INSERT INTO containers (
		id, name, image, status, isolation_type, cpu_limit, memory_limit, disk_limit,
		network_mode, created_at, started_at, stopped_at, config, metadata
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	configJSON := toJSON(container.Config)
	metadataJSON := toJSON(container.Metadata)
	
	_, err := dao.Exec(query,
		container.ID,
		container.Name,
		container.Image,
		container.Status,
		container.IsolationType,
		container.CPULimit,
		container.MemoryLimit,
		container.DiskLimit,
		container.NetworkMode,
		formatTime(container.CreatedAt),
		nullTime(container.StartedAt),
		nullTime(container.StoppedAt),
		configJSON,
		metadataJSON,
	)
	return err
}

func (dao *ContainerDAO) GetByID(id string) (*Container, error) {
	query := `SELECT id, name, image, status, isolation_type, cpu_limit, memory_limit, disk_limit,
		network_mode, created_at, started_at, stopped_at, config, metadata
		FROM containers WHERE id = ?`
	
	row := dao.QueryRow(query, id)
	
	container := &Container{}
	var createdAt, startedAt, stoppedAt sql.NullString
	var configJSON, metadataJSON sql.NullString
	
	err := row.Scan(
		&container.ID,
		&container.Name,
		&container.Image,
		&container.Status,
		&container.IsolationType,
		&container.CPULimit,
		&container.MemoryLimit,
		&container.DiskLimit,
		&container.NetworkMode,
		&createdAt,
		&startedAt,
		&stoppedAt,
		&configJSON,
		&metadataJSON,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	container.CreatedAt = parseTime(createdAt.String)
	container.StartedAt = fromNullTime(startedAt)
	container.StoppedAt = fromNullTime(stoppedAt)
	
	if configJSON.Valid {
		fromJSON(configJSON.String, &container.Config)
	}
	if metadataJSON.Valid {
		fromJSON(metadataJSON.String, &container.Metadata)
	}
	
	return container, nil
}

func (dao *ContainerDAO) Update(container *Container) error {
	query := `UPDATE containers SET
		name = ?, image = ?, status = ?, isolation_type = ?, cpu_limit = ?, memory_limit = ?,
		disk_limit = ?, network_mode = ?, started_at = ?, stopped_at = ?, config = ?, metadata = ?
		WHERE id = ?`
	
	configJSON := toJSON(container.Config)
	metadataJSON := toJSON(container.Metadata)
	
	_, err := dao.Exec(query,
		container.Name,
		container.Image,
		container.Status,
		container.IsolationType,
		container.CPULimit,
		container.MemoryLimit,
		container.DiskLimit,
		container.NetworkMode,
		nullTime(container.StartedAt),
		nullTime(container.StoppedAt),
		configJSON,
		metadataJSON,
		container.ID,
	)
	return err
}

func (dao *ContainerDAO) UpdateStatus(id string, status string) error {
	query := `UPDATE containers SET status = ? WHERE id = ?`
	_, err := dao.Exec(query, status, id)
	return err
}

func (dao *ContainerDAO) UpdateStartedAt(id string, startedAt time.Time) error {
	query := `UPDATE containers SET started_at = ?, status = 'running' WHERE id = ?`
	_, err := dao.Exec(query, formatTime(startedAt), id)
	return err
}

func (dao *ContainerDAO) UpdateStoppedAt(id string, stoppedAt time.Time) error {
	query := `UPDATE containers SET stopped_at = ?, status = 'stopped' WHERE id = ?`
	_, err := dao.Exec(query, formatTime(stoppedAt), id)
	return err
}

func (dao *ContainerDAO) Delete(id string) error {
	query := `DELETE FROM containers WHERE id = ?`
	_, err := dao.Exec(query, id)
	return err
}

func (dao *ContainerDAO) List(limit, offset int) ([]*Container, error) {
	query := `SELECT id, name, image, status, isolation_type, cpu_limit, memory_limit, disk_limit,
		network_mode, created_at, started_at, stopped_at, config, metadata
		FROM containers ORDER BY created_at DESC LIMIT ? OFFSET ?`
	
	rows, err := dao.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var containers []*Container
	for rows.Next() {
		container := &Container{}
		var createdAt, startedAt, stoppedAt sql.NullString
		var configJSON, metadataJSON sql.NullString
		
		err := rows.Scan(
			&container.ID,
			&container.Name,
			&container.Image,
			&container.Status,
			&container.IsolationType,
			&container.CPULimit,
			&container.MemoryLimit,
			&container.DiskLimit,
			&container.NetworkMode,
			&createdAt,
			&startedAt,
			&stoppedAt,
			&configJSON,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}
		
		container.CreatedAt = parseTime(createdAt.String)
		container.StartedAt = fromNullTime(startedAt)
		container.StoppedAt = fromNullTime(stoppedAt)
		
		if configJSON.Valid {
			fromJSON(configJSON.String, &container.Config)
		}
		if metadataJSON.Valid {
			fromJSON(metadataJSON.String, &container.Metadata)
		}
		
		containers = append(containers, container)
	}
	
	return containers, nil
}

func (dao *ContainerDAO) ListByStatus(status string) ([]*Container, error) {
	query := `SELECT id, name, image, status, isolation_type, cpu_limit, memory_limit, disk_limit,
		network_mode, created_at, started_at, stopped_at, config, metadata
		FROM containers WHERE status = ? ORDER BY created_at DESC`
	
	rows, err := dao.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var containers []*Container
	for rows.Next() {
		container := &Container{}
		var createdAt, startedAt, stoppedAt sql.NullString
		var configJSON, metadataJSON sql.NullString
		
		err := rows.Scan(
			&container.ID,
			&container.Name,
			&container.Image,
			&container.Status,
			&container.IsolationType,
			&container.CPULimit,
			&container.MemoryLimit,
			&container.DiskLimit,
			&container.NetworkMode,
			&createdAt,
			&startedAt,
			&stoppedAt,
			&configJSON,
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}
		
		container.CreatedAt = parseTime(createdAt.String)
		container.StartedAt = fromNullTime(startedAt)
		container.StoppedAt = fromNullTime(stoppedAt)
		
		if configJSON.Valid {
			fromJSON(configJSON.String, &container.Config)
		}
		if metadataJSON.Valid {
			fromJSON(metadataJSON.String, &container.Metadata)
		}
		
		containers = append(containers, container)
	}
	
	return containers, nil
}

func (dao *ContainerDAO) ListRunning() ([]*Container, error) {
	return dao.ListByStatus("running")
}

func (dao *ContainerDAO) Count() (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM containers").Scan(&count)
	return count, err
}

func (dao *ContainerDAO) CountByStatus(status string) (int, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM containers WHERE status = ?", status).Scan(&count)
	return count, err
}

func (dao *ContainerDAO) Exists(id string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM containers WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *ContainerDAO) ExistsByName(name string) (bool, error) {
	var count int
	err := dao.QueryRow("SELECT COUNT(*) FROM containers WHERE name = ?", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}