package database

import (
	"database/sql"
	"fmt"
	"strconv"
)

// SettingsStore provides key-value settings operations.
type SettingsStore struct {
	db *sql.DB
}

// NewSettingsStore creates a new SettingsStore.
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

// Open ensures the table exists and seeds defaults.
func (s *SettingsStore) Open() error {
	if err := s.ensureTable(); err != nil {
		return err
	}
	return s.initDefaults()
}

// All returns all settings as a map.
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

// Set updates a setting value.
func (s *SettingsStore) Set(key, value string) error {
	_, err := s.db.Exec(`INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, key, value)
	if err != nil {
		return fmt.Errorf("set setting %s: %w", key, err)
	}
	return nil
}

// SetInt updates an integer setting value.
func (s *SettingsStore) SetInt(key string, value int64) error {
	return s.Set(key, strconv.FormatInt(value, 10))
}
