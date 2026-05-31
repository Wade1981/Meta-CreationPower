package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type BaseDAO struct {
	db *sql.DB
}

func NewBaseDAO(db *sql.DB) *BaseDAO {
	return &BaseDAO{db: db}
}

func (dao *BaseDAO) GetDB() *sql.DB {
	return dao.db
}

func (dao *BaseDAO) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dao.db.Exec(query, args...)
}

func (dao *BaseDAO) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return dao.db.Query(query, args...)
}

func (dao *BaseDAO) QueryRow(query string, args ...interface{}) *sql.Row {
	return dao.db.QueryRow(query, args...)
}

func (dao *BaseDAO) BeginTx() (*sql.Tx, error) {
	return dao.db.Begin()
}

func (dao *BaseDAO) WithTx(fn func(*sql.Tx) error) error {
	tx, err := dao.db.Begin()
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

func (dao *BaseDAO) Count(table string, where string, args ...interface{}) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	if where != "" {
		query += " WHERE " + where
	}
	
	var count int
	err := dao.db.QueryRow(query, args...).Scan(&count)
	return count, err
}

func (dao *BaseDAO) Exists(table string, where string, args ...interface{}) (bool, error) {
	count, err := dao.Count(table, where, args...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *BaseDAO) Delete(table string, where string, args ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	return dao.db.Exec(query, args...)
}

func (dao *BaseDAO) Update(table string, updates string, where string, args ...interface{}) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, updates, where)
	return dao.db.Exec(query, args...)
}

func toJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func fromJSON(s string, v interface{}) error {
	if s == "" {
		return nil
	}
	return json.Unmarshal([]byte(s), v)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func nullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: i != 0}
}

func nullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: f != 0}
}

func nullTime(t time.Time) sql.NullString {
	if t.IsZero() {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: t.Format(time.RFC3339), Valid: true}
}

func fromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func fromNullInt64(ns sql.NullInt64) int64 {
	if ns.Valid {
		return ns.Int64
	}
	return 0
}

func fromNullFloat64(ns sql.NullFloat64) float64 {
	if ns.Valid {
		return ns.Float64
	}
	return 0
}

func fromNullTime(ns sql.NullString) time.Time {
	if ns.Valid {
		return parseTime(ns.String)
	}
	return time.Time{}
}