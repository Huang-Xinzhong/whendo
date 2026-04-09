package database

import (
	"database/sql"
	"fmt"

	"whendo/internal/models"
)

// WorkspaceStore 提供工作区的存储操作。
type WorkspaceStore struct {
	db *sql.DB
}

// NewWorkspaceStore 创建一个新的 WorkspaceStore 实例。
func NewWorkspaceStore(db *sql.DB) *WorkspaceStore {
	return &WorkspaceStore{db: db}
}

// List 返回所有工作区，按 created_at 升序排列。
func (s *WorkspaceStore) List() ([]models.Workspace, error) {
	rows, err := s.db.Query(`SELECT id, name, created_at, updated_at FROM workspaces ORDER BY created_at ASC`)
	if err != nil {
		return nil, fmt.Errorf("query workspaces: %w", err)
	}
	defer rows.Close()

	var list []models.Workspace
	for rows.Next() {
		var ws models.Workspace
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.CreatedAt, &ws.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan workspace: %w", err)
		}
		list = append(list, ws)
	}
	return list, rows.Err()
}

// Create 插入一个新的工作区。
func (s *WorkspaceStore) Create(name string) (*models.Workspace, error) {
	res, err := s.db.Exec(`INSERT INTO workspaces (name) VALUES (?)`, name)
	if err != nil {
		return nil, fmt.Errorf("insert workspace: %w", err)
	}
	id, _ := res.LastInsertId()
	return s.Get(id)
}

// Get 根据 ID 返回工作区。
func (s *WorkspaceStore) Get(id int64) (*models.Workspace, error) {
	var ws models.Workspace
	row := s.db.QueryRow(`SELECT id, name, created_at, updated_at FROM workspaces WHERE id = ?`, id)
	if err := row.Scan(&ws.ID, &ws.Name, &ws.CreatedAt, &ws.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("workspace %d not found", id)
		}
		return nil, fmt.Errorf("scan workspace: %w", err)
	}
	return &ws, nil
}

// Update 重命名工作区。
func (s *WorkspaceStore) Update(id int64, name string) (*models.Workspace, error) {
	_, err := s.db.Exec(`UPDATE workspaces SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, name, id)
	if err != nil {
		return nil, fmt.Errorf("update workspace: %w", err)
	}
	return s.Get(id)
}

// Delete 删除工作区，通过外键级联删除其任务。
func (s *WorkspaceStore) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM workspaces WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete workspace: %w", err)
	}
	return nil
}
