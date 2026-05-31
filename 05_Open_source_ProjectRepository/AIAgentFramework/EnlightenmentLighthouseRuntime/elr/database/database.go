package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type DatabaseManager struct {
	dataDB     *sql.DB
	metricsDB  *sql.DB
	logsDB     *sql.DB
	dataDir    string
	mutex      sync.RWMutex
	initialized bool
}

var (
	globalDBManager *DatabaseManager
	dbMutex         sync.Mutex
)

func GetDatabaseManager() *DatabaseManager {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	
	if globalDBManager == nil {
		globalDBManager = &DatabaseManager{}
	}
	return globalDBManager
}

func InitDatabaseManager(dataDir string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	
	if globalDBManager == nil {
		globalDBManager = &DatabaseManager{}
	}
	
	return globalDBManager.Initialize(dataDir)
}

func (dm *DatabaseManager) Initialize(dataDir string) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	
	if dm.initialized {
		return nil
	}
	
	dm.dataDir = dataDir
	
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}
	
	dataDBPath := filepath.Join(dataDir, "elr_data.db")
	metricsDBPath := filepath.Join(dataDir, "elr_metrics.db")
	logsDBPath := filepath.Join(dataDir, "elr_logs.db")
	
	dataDB, err := sql.Open("sqlite", dataDBPath)
	if err != nil {
		return fmt.Errorf("failed to open data database: %v", err)
	}
	dataDB.SetMaxOpenConns(25)
	dataDB.SetMaxIdleConns(10)
	dataDB.SetConnMaxLifetime(5 * time.Minute)
	dm.dataDB = dataDB
	
	metricsDB, err := sql.Open("sqlite", metricsDBPath)
	if err != nil {
		return fmt.Errorf("failed to open metrics database: %v", err)
	}
	metricsDB.SetMaxOpenConns(10)
	metricsDB.SetMaxIdleConns(5)
	metricsDB.SetConnMaxLifetime(5 * time.Minute)
	dm.metricsDB = metricsDB
	
	logsDB, err := sql.Open("sqlite", logsDBPath)
	if err != nil {
		return fmt.Errorf("failed to open logs database: %v", err)
	}
	logsDB.SetMaxOpenConns(10)
	logsDB.SetMaxIdleConns(5)
	logsDB.SetConnMaxLifetime(5 * time.Minute)
	dm.logsDB = logsDB
	
	if err := dm.initSchema(); err != nil {
		return fmt.Errorf("failed to initialize schema: %v", err)
	}
	
	dm.initialized = true
	fmt.Println("Database initialized successfully")
	return nil
}

func (dm *DatabaseManager) GetDataDB() *sql.DB {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	return dm.dataDB
}

func (dm *DatabaseManager) GetMetricsDB() *sql.DB {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	return dm.metricsDB
}

func (dm *DatabaseManager) GetLogsDB() *sql.DB {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	return dm.logsDB
}

func (dm *DatabaseManager) Close() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	
	var errors []error
	
	if dm.dataDB != nil {
		if err := dm.dataDB.Close(); err != nil {
			errors = append(errors, err)
		}
	}
	
	if dm.metricsDB != nil {
		if err := dm.metricsDB.Close(); err != nil {
			errors = append(errors, err)
		}
	}
	
	if dm.logsDB != nil {
		if err := dm.logsDB.Close(); err != nil {
			errors = append(errors, err)
		}
	}
	
	dm.initialized = false
	
	if len(errors) > 0 {
		return fmt.Errorf("errors closing databases: %v", errors)
	}
	return nil
}

func (dm *DatabaseManager) IsInitialized() bool {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	return dm.initialized
}

func (dm *DatabaseManager) GetDataDir() string {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	return dm.dataDir
}

func (dm *DatabaseManager) Backup(backupDir string) error {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	
	if !dm.initialized {
		return fmt.Errorf("database not initialized")
	}
	
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}
	
	timestamp := time.Now().Format("20260102_150405")
	
	dataDBPath := filepath.Join(dm.dataDir, "elr_data.db")
	metricsDBPath := filepath.Join(dm.dataDir, "elr_metrics.db")
	logsDBPath := filepath.Join(dm.dataDir, "elr_logs.db")
	
	dataBackup := filepath.Join(backupDir, fmt.Sprintf("elr_data_%s.db", timestamp))
	metricsBackup := filepath.Join(backupDir, fmt.Sprintf("elr_metrics_%s.db", timestamp))
	logsBackup := filepath.Join(backupDir, fmt.Sprintf("elr_logs_%s.db", timestamp))
	
	if err := copyFile(dataDBPath, dataBackup); err != nil {
		return fmt.Errorf("failed to backup data database: %v", err)
	}
	
	if err := copyFile(metricsDBPath, metricsBackup); err != nil {
		return fmt.Errorf("failed to backup metrics database: %v", err)
	}
	
	if err := copyFile(logsDBPath, logsBackup); err != nil {
		return fmt.Errorf("failed to backup logs database: %v", err)
	}
	
	fmt.Printf("Database backup completed: %s\n", backupDir)
	return nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}