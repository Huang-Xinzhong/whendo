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
	debugLog("WorkspaceStore.List 查询工作区列表")
	rows, err := s.db.Query(`SELECT id, name, created_at, updated_at FROM workspaces ORDER BY created_at ASC`)
	if err != nil {
		debugLog("WorkspaceStore.List 查询失败: %v", err)
		return nil, fmt.Errorf("query workspaces: %w", err)
	}
	defer rows.Close()

	var list []models.Workspace
	for rows.Next() {
		var ws models.Workspace
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.CreatedAt, &ws.UpdatedAt); err != nil {
			debugLog("WorkspaceStore.List 扫描失败: %v", err)
			return nil, fmt.Errorf("scan workspace: %w", err)
		}
		list = append(list, ws)
	}
	debugLog("WorkspaceStore.List 成功，共 %d 个工作区", len(list))
	return list, rows.Err()
}

// Create 插入一个新的工作区。
func (s *WorkspaceStore) Create(name string) (*models.Workspace, error) {
	debugLog("WorkspaceStore.Create 插入工作区: name=%s", name)
	res, err := s.db.Exec(`INSERT INTO workspaces (name) VALUES (?)`, name)
	if err != nil {
		debugLog("WorkspaceStore.Create 插入失败: %v", err)
		return nil, fmt.Errorf("insert workspace: %w", err)
	}
	id, _ := res.LastInsertId()
	debugLog("WorkspaceStore.Create 成功: id=%d", id)
	return s.Get(id)
}

// Get 根据 ID 返回工作区。
func (s *WorkspaceStore) Get(id int64) (*models.Workspace, error) {
	debugLog("WorkspaceStore.Get 查询工作区: id=%d", id)
	var ws models.Workspace
	row := s.db.QueryRow(`SELECT id, name, created_at, updated_at FROM workspaces WHERE id = ?`, id)
	if err := row.Scan(&ws.ID, &ws.Name, &ws.CreatedAt, &ws.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			debugLog("WorkspaceStore.Get 工作区不存在: id=%d", id)
			return nil, fmt.Errorf("workspace %d not found", id)
		}
		debugLog("WorkspaceStore.Get 扫描失败: %v", err)
		return nil, fmt.Errorf("scan workspace: %w", err)
	}
	debugLog("WorkspaceStore.Get 成功: id=%d name=%s", ws.ID, ws.Name)
	return &ws, nil
}

// Update 重命名工作区。
func (s *WorkspaceStore) Update(id int64, name string) (*models.Workspace, error) {
	debugLog("WorkspaceStore.Update 更新工作区: id=%d name=%s", id, name)
	_, err := s.db.Exec(`UPDATE workspaces SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, name, id)
	if err != nil {
		debugLog("WorkspaceStore.Update 失败: %v", err)
		return nil, fmt.Errorf("update workspace: %w", err)
	}
	debugLog("WorkspaceStore.Update 成功: id=%d", id)
	return s.Get(id)
}

// Delete 删除工作区，通过外键级联删除其任务。
func (s *WorkspaceStore) Delete(id int64) error {
	debugLog("WorkspaceStore.Delete 删除工作区: id=%d", id)
	_, err := s.db.Exec(`DELETE FROM workspaces WHERE id = ?`, id)
	if err != nil {
		debugLog("WorkspaceStore.Delete 失败: %v", err)
		return fmt.Errorf("delete workspace: %w", err)
	}
	debugLog("WorkspaceStore.Delete 成功: id=%d", id)
	return nil
}
