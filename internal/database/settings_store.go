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
	debugLog("SettingsStore.ensureTable 检查设置表")
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	)`)
	if err != nil {
		debugLog("SettingsStore.ensureTable 失败: %v", err)
		return fmt.Errorf("create settings table: %w", err)
	}
	debugLog("SettingsStore.ensureTable 成功")
	return nil
}

func (s *SettingsStore) initDefaults() error {
	debugLog("SettingsStore.initDefaults 初始化默认设置")
	defaults := map[string]string{
		"default_workspace_id": "1",
		"default_sort":         "created_at_desc",
	}
	for k, v := range defaults {
		_, err := s.db.Exec(`INSERT OR IGNORE INTO settings (key, value) VALUES (?, ?)`, k, v)
		if err != nil {
			debugLog("SettingsStore.initDefaults 设置 %s 失败: %v", k, err)
			return fmt.Errorf("init setting %s: %w", k, err)
		}
	}
	debugLog("SettingsStore.initDefaults 完成")
	return nil
}

// Open 确保表存在并初始化默认值。
func (s *SettingsStore) Open() error {
	debugLog("SettingsStore.Open 初始化")
	if err := s.ensureTable(); err != nil {
		return err
	}
	return s.initDefaults()
}

// All 以 map 形式返回所有设置。
func (s *SettingsStore) All() (map[string]string, error) {
	debugLog("SettingsStore.All 查询所有设置")
	rows, err := s.db.Query(`SELECT key, value FROM settings`)
	if err != nil {
		debugLog("SettingsStore.All 查询失败: %v", err)
		return nil, fmt.Errorf("query settings: %w", err)
	}
	defer rows.Close()

	m := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			debugLog("SettingsStore.All 扫描失败: %v", err)
			return nil, fmt.Errorf("scan setting: %w", err)
		}
		m[k] = v
	}
	debugLog("SettingsStore.All 成功，共 %d 条", len(m))
	return m, rows.Err()
}

// Set 更新设置值。
func (s *SettingsStore) Set(key, value string) error {
	debugLog("SettingsStore.Set 更新设置: key=%s value=%s", key, value)
	_, err := s.db.Exec(`INSERT INTO settings (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, key, value)
	if err != nil {
		debugLog("SettingsStore.Set 失败: %v", err)
		return fmt.Errorf("set setting %s: %w", key, err)
	}
	debugLog("SettingsStore.Set 成功: key=%s", key)
	return nil
}

// SetInt 以整数形式更新设置值。
func (s *SettingsStore) SetInt(key string, value int64) error {
	debugLog("SettingsStore.SetInt 更新设置: key=%s value=%d", key, value)
	return s.Set(key, strconv.FormatInt(value, 10))
}
