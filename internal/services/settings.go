package services

import (
	"database/sql"
	"fmt"
	"strconv"

	"whendo/internal/database"
	"whendo/internal/requests"
)

// SettingsService 提供设置相关的业务逻辑。
type SettingsService struct {
	store *database.SettingsStore
}

// NewSettingsService 创建一个新的 SettingsService 实例。
func NewSettingsService(db *sql.DB) *SettingsService {
	s := &SettingsService{store: database.NewSettingsStore(db)}
	_ = s.store.Open()
	return s
}

// Get 以 map 形式返回所有设置。
func (s *SettingsService) Get() (map[string]string, error) {
	m, err := s.store.All()
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}
	return m, nil
}

// Update 更新设置。
func (s *SettingsService) Update(req requests.SettingsUpdateReq) error {
	if req.DefaultWorkspaceID != "" {
		id, err := strconv.ParseInt(req.DefaultWorkspaceID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid default_workspace_id: %w", err)
		}
		if err := s.store.SetInt("default_workspace_id", id); err != nil {
			return fmt.Errorf("update default_workspace_id: %w", err)
		}
	}
	if req.DefaultSort != "" {
		if err := s.store.Set("default_sort", req.DefaultSort); err != nil {
			return fmt.Errorf("update default_sort: %w", err)
		}
	}
	return nil
}

// DefaultWorkspaceID 返回默认工作区 ID。
func (s *SettingsService) DefaultWorkspaceID() int64 {
	m, err := s.store.All()
	if err != nil {
		return 1
	}
	v, ok := m["default_workspace_id"]
	if !ok {
		return 1
	}
	id, _ := strconv.ParseInt(v, 10, 64)
	if id == 0 {
		return 1
	}
	return id
}
