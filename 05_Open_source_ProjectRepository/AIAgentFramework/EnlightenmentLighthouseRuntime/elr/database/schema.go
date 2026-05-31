package database

import (
	"database/sql"
	"fmt"
)

func (dm *DatabaseManager) initSchema() error {
	if err := dm.initDataSchema(); err != nil {
		return err
	}
	
	if err := dm.initMetricsSchema(); err != nil {
		return err
	}
	
	if err := dm.initLogsSchema(); err != nil {
		return err
	}
	
	return nil
}

func (dm *DatabaseManager) initDataSchema() error {
	db := dm.dataDB
	
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS containers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			image TEXT NOT NULL,
			status TEXT NOT NULL,
			isolation_type TEXT,
			cpu_limit INTEGER,
			memory_limit INTEGER,
			disk_limit INTEGER,
			network_mode TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			started_at TIMESTAMP,
			stopped_at TIMESTAMP,
			config TEXT,
			metadata TEXT
		)`,
		
		`CREATE TABLE IF NOT EXISTS sandboxes (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			container_id TEXT,
			status TEXT NOT NULL,
			process_id INTEGER,
			pipe_name TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			started_at TIMESTAMP,
			uptime INTEGER,
			config TEXT,
			FOREIGN KEY (container_id) REFERENCES containers(id)
		)`,
		
		`CREATE TABLE IF NOT EXISTS models (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			version TEXT NOT NULL,
			path TEXT NOT NULL,
			url TEXT,
			size INTEGER,
			status TEXT NOT NULL,
			dependencies TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP,
			metadata TEXT
		)`,
		
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			version TEXT,
			description TEXT,
			path TEXT NOT NULL,
			status TEXT NOT NULL,
			sandbox_id TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deployed_at TIMESTAMP,
			FOREIGN KEY (sandbox_id) REFERENCES sandboxes(id)
		)`,
		
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			email TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login_at TIMESTAMP,
			is_active INTEGER DEFAULT 1
		)`,
		
		`CREATE TABLE IF NOT EXISTS admins (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			email TEXT,
			group_id TEXT,
			token_id TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login_at TIMESTAMP,
			is_active INTEGER DEFAULT 1,
			metadata TEXT,
			FOREIGN KEY (group_id) REFERENCES admin_groups(id),
			FOREIGN KEY (token_id) REFERENCES tokens(id)
		)`,
		
		`CREATE TABLE IF NOT EXISTS admin_groups (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			permissions TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS tokens (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			admin_id TEXT,
			token_hash TEXT NOT NULL UNIQUE,
			token_type TEXT NOT NULL,
			endpoint_type TEXT NOT NULL,
			permissions TEXT NOT NULL,
			expires_at TIMESTAMP,
			issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_used_at TIMESTAMP,
			status TEXT NOT NULL,
			refresh_token_hash TEXT,
			usage_count INTEGER DEFAULT 0,
			max_usage INTEGER,
			allowed_ips TEXT,
			metadata TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (admin_id) REFERENCES admins(id)
		)`,
		
		`CREATE TABLE IF NOT EXISTS token_stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token_id TEXT NOT NULL,
			date DATE NOT NULL,
			request_count INTEGER DEFAULT 0,
			success_count INTEGER DEFAULT 0,
			failure_count INTEGER DEFAULT 0,
			total_response_time REAL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (token_id) REFERENCES tokens(id)
		)`,
		
		`CREATE TABLE IF NOT EXISTS sandbox_container_mapping (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sandbox_id TEXT NOT NULL UNIQUE,
			container_id TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	
	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}
	
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_containers_status ON containers(status)`,
		`CREATE INDEX IF NOT EXISTS idx_containers_created_at ON containers(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_sandboxes_container_id ON sandboxes(container_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sandboxes_status ON sandboxes(status)`,
		`CREATE INDEX IF NOT EXISTS idx_sandboxes_process_id ON sandboxes(process_id)`,
		`CREATE INDEX IF NOT EXISTS idx_models_type ON models(type)`,
		`CREATE INDEX IF NOT EXISTS idx_models_status ON models(status)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_sandbox_id ON projects(sandbox_id)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status)`,
		`CREATE INDEX IF NOT EXISTS idx_admins_group_id ON admins(group_id)`,
		`CREATE INDEX IF NOT EXISTS idx_admins_role ON admins(role)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_admin_id ON tokens(admin_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_status ON tokens(status)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_token_type ON tokens(token_type)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_endpoint_type ON tokens(endpoint_type)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_expires_at ON tokens(expires_at)`,
		`CREATE INDEX IF NOT EXISTS idx_token_stats_token_id ON token_stats(token_id)`,
		`CREATE INDEX IF NOT EXISTS idx_token_stats_date ON token_stats(date)`,
	}
	
	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return fmt.Errorf("failed to create index: %v", err)
		}
	}
	
	return nil
}

func (dm *DatabaseManager) initMetricsSchema() error {
	db := dm.metricsDB
	
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS resource_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			target_id TEXT NOT NULL,
			target_type TEXT NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			cpu_usage REAL,
			memory_usage INTEGER,
			memory_percent REAL,
			disk_usage INTEGER,
			disk_percent REAL,
			network_in INTEGER,
			network_out INTEGER,
			gpu_usage REAL,
			gpu_memory INTEGER,
			process_count INTEGER,
			thread_count INTEGER
		)`,
		
		`CREATE TABLE IF NOT EXISTS alert_rules (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			target_type TEXT NOT NULL,
			metric_type TEXT NOT NULL,
			operator TEXT NOT NULL,
			threshold REAL NOT NULL,
			duration INTEGER DEFAULT 60,
			severity TEXT NOT NULL,
			enabled INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			rule_id TEXT,
			target_id TEXT NOT NULL,
			target_type TEXT NOT NULL,
			metric_type TEXT NOT NULL,
			current_value REAL NOT NULL,
			threshold REAL NOT NULL,
			severity TEXT NOT NULL,
			status TEXT NOT NULL,
			message TEXT,
			fired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			resolved_at TIMESTAMP,
			FOREIGN KEY (rule_id) REFERENCES alert_rules(id)
		)`,
	}
	
	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create metrics table: %v", err)
		}
	}
	
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_metrics_target ON resource_metrics(target_id, target_type)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON resource_metrics(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_alert_rules_target ON alert_rules(target_type, metric_type)`,
		`CREATE INDEX IF NOT EXISTS idx_alert_rules_enabled ON alert_rules(enabled)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_rule_id ON alerts(rule_id)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_severity ON alerts(severity)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_timestamp ON alerts(fired_at)`,
	}
	
	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return fmt.Errorf("failed to create metrics index: %v", err)
		}
	}
	
	return nil
}

func (dm *DatabaseManager) initLogsSchema() error {
	db := dm.logsDB
	
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS access_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			container_id TEXT,
			client_ip TEXT NOT NULL,
			request_method TEXT,
			request_path TEXT,
			status_code INTEGER,
			response_time REAL,
			user_agent TEXT,
			token_id TEXT,
			endpoint_type TEXT,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT,
			admin_id TEXT,
			operation TEXT NOT NULL,
			target_type TEXT NOT NULL,
			target_id TEXT,
			target_name TEXT,
			details TEXT,
			result TEXT NOT NULL,
			error_message TEXT,
			ip_address TEXT,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		`CREATE TABLE IF NOT EXISTS runtime_data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sandbox_id TEXT,
			model_id TEXT,
			data_type TEXT NOT NULL,
			content TEXT,
			metadata TEXT,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	
	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create logs table: %v", err)
		}
	}
	
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_access_container ON access_logs(container_id)`,
		`CREATE INDEX IF NOT EXISTS idx_access_timestamp ON access_logs(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_access_client_ip ON access_logs(client_ip)`,
		`CREATE INDEX IF NOT EXISTS idx_access_token_id ON access_logs(token_id)`,
		`CREATE INDEX IF NOT EXISTS idx_access_endpoint_type ON access_logs(endpoint_type)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_user ON audit_logs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_admin ON audit_logs(admin_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_target ON audit_logs(target_type, target_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit_logs(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_runtime_sandbox ON runtime_data(sandbox_id)`,
		`CREATE INDEX IF NOT EXISTS idx_runtime_model ON runtime_data(model_id)`,
		`CREATE INDEX IF NOT EXISTS idx_runtime_timestamp ON runtime_data(timestamp)`,
	}
	
	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return fmt.Errorf("failed to create logs index: %v", err)
		}
	}
	
	return nil
}

func (dm *DatabaseManager) ExecInDataDB(query string, args ...interface{}) (sql.Result, error) {
	return dm.dataDB.Exec(query, args...)
}

func (dm *DatabaseManager) QueryInDataDB(query string, args ...interface{}) (*sql.Rows, error) {
	return dm.dataDB.Query(query, args...)
}

func (dm *DatabaseManager) QueryRowInDataDB(query string, args ...interface{}) *sql.Row {
	return dm.dataDB.QueryRow(query, args...)
}

func (dm *DatabaseManager) ExecInMetricsDB(query string, args ...interface{}) (sql.Result, error) {
	return dm.metricsDB.Exec(query, args...)
}

func (dm *DatabaseManager) QueryInMetricsDB(query string, args ...interface{}) (*sql.Rows, error) {
	return dm.metricsDB.Query(query, args...)
}

func (dm *DatabaseManager) ExecInLogsDB(query string, args ...interface{}) (sql.Result, error) {
	return dm.logsDB.Exec(query, args...)
}

func (dm *DatabaseManager) QueryInLogsDB(query string, args ...interface{}) (*sql.Rows, error) {
	return dm.logsDB.Query(query, args...)
}

func (dm *DatabaseManager) BeginTransaction(db *sql.DB) (*sql.Tx, error) {
	return db.Begin()
}

func (dm *DatabaseManager) WithTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit()
}