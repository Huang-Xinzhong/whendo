package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Open initialises the SQLite database, runs migrations and returns the connection.
func Open() (*sql.DB, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("get config dir: %w", err)
	}
	appDir := filepath.Join(dir, "whendo")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, fmt.Errorf("create app dir: %w", err)
	}

	dbPath := filepath.Join(appDir, "data.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	migrations := []string{
		// 001_init.sql content is embedded here for self-containment.
		`CREATE TABLE IF NOT EXISTS workspaces (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			workspace_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			type TEXT DEFAULT 'todo' CHECK(type IN ('todo', 'reminder')),
			due_at DATETIME,
			remind_at DATETIME,
			is_completed BOOLEAN DEFAULT 0,
			start_time TEXT,
			end_time TEXT,
			interval_value INTEGER,
			interval_unit TEXT CHECK(interval_unit IN ('minute', 'hour')),
			repeat_mode TEXT CHECK(repeat_mode IN ('daily', 'workday', 'weekly', 'monthly')),
			weekdays TEXT,
			month_day INTEGER,
			next_trigger_at DATETIME,
			paused_date DATE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_workspace_id ON tasks(workspace_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_is_completed ON tasks(is_completed);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_next_trigger_at ON tasks(next_trigger_at);`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`,
		`INSERT OR IGNORE INTO settings (key, value) VALUES ('default_workspace_id', '1'), ('default_sort', 'created_at_desc');`,
		`INSERT OR IGNORE INTO workspaces (id, name) VALUES (1, '家庭'), (2, '工作'), (3, '个人');`,
	}

	for i, q := range migrations {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("migration step %d: %w", i, err)
		}
	}
	return nil
}
