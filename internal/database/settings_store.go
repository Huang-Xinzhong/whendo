package database

import (
	"database/sql"
	"fmt"
	"strconv"
)

// SettingsStore 提供键值对设置的存储操作。
type SettingsStore struct {
	db *sql.DB
}

// NewSettingsStore 创建一个新的 SettingsStore 实例。
func NewSettingsStore(db *sql.DB) *SettingsStore {
	return &SettingsStore{db: db}
}

func (s *SettingsStore) ensureTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("create settings table: %w", err)
	}
	return nil
}

func (s *SettingsStore) initDefaults() error {
	defaults := map[string]string{
		"default_workspace_id": "1",
		"default_sort":         "created_at_desc",
	}
	for k, v := range defaults {
		_, err := s.db.Exec(`INSERT OR IGNORE INTO settings (key, value) VALUES (?, ?)`, k, v)
		if err != nil {
			return fmt.Errorf("init setting %s: %w", k, err)
		}
	}
	return nil
}

// Open 确保表存在并初始化默认值。
func (s *SettingsStore) Open() error {
	if err := s.ensureTable(); err != nil {
		return err
	}
	return s.initDefaults()
}

// All 以 map 形式返回所有设置。
func (s *SettingsStore) All() (map[string]string, error) {
	rows, err := s.db.Query(`SELECT key, value FROM settings`)
	if err != nil {
		return nil, fmt.Errorf("query settings: %w", err)
	}
	defer rows.Close()

	m := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, fmt.Errorf("scan setting: %w", err)
		}
		m[k] = v
	}
	return m, rows.Err()
}

// Set 更新设置值。
func (s *SettingsStore) Set(key, value string) error {
	_, err := s.db.Exec(`INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, key, value)
	if err != nil {
		return fmt.Errorf("set setting %s: %w", key, err)
	}
	return nil
}

// SetInt 以整数形式更新设置值。
func (s *SettingsStore) SetInt(key string, value int64) error {
	return s.Set(key, strconv.FormatInt(value, 10))
}
