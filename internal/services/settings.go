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
	fmt.Println("[调试] SettingsService.Get 获取设置")
	m, err := s.store.All()
	if err != nil {
		fmt.Printf("[调试] SettingsService.Get 失败: %v\n", err)
		return nil, fmt.Errorf("get settings: %w", err)
	}
	fmt.Printf("[调试] SettingsService.Get 成功: %+v\n", m)
	return m, nil
}

// Update 更新设置。
func (s *SettingsService) Update(req requests.SettingsUpdateReq) error {
	fmt.Printf("[调试] SettingsService.Update 更新设置: %+v\n", req)
	if req.DefaultWorkspaceID != "" {
		id, err := strconv.ParseInt(req.DefaultWorkspaceID, 10, 64)
		if err != nil {
			fmt.Printf("[调试] SettingsService.Update 解析 default_workspace_id 失败: %v\n", err)
			return fmt.Errorf("invalid default_workspace_id: %w", err)
		}
		if err := s.store.SetInt("default_workspace_id", id); err != nil {
			fmt.Printf("[调试] SettingsService.Update 设置 default_workspace_id 失败: %v\n", err)
			return fmt.Errorf("update default_workspace_id: %w", err)
		}
	}
	if req.DefaultSort != "" {
		if err := s.store.Set("default_sort", req.DefaultSort); err != nil {
			fmt.Printf("[调试] SettingsService.Update 设置 default_sort 失败: %v\n", err)
			return fmt.Errorf("update default_sort: %w", err)
		}
	}
	fmt.Println("[调试] SettingsService.Update 成功")
	return nil
}

// DefaultWorkspaceID 返回默认工作区 ID。
func (s *SettingsService) DefaultWorkspaceID() int64 {
	m, err := s.store.All()
	if err != nil {
		fmt.Printf("[调试] SettingsService.DefaultWorkspaceID 读取失败，使用默认值 1: %v\n", err)
		return 1
	}
	v, ok := m["default_workspace_id"]
	if !ok {
		fmt.Println("[调试] SettingsService.DefaultWorkspaceID 未配置，使用默认值 1")
		return 1
	}
	id, _ := strconv.ParseInt(v, 10, 64)
	if id == 0 {
		fmt.Println("[调试] SettingsService.DefaultWorkspaceID 解析为 0，使用默认值 1")
		return 1
	}
	fmt.Printf("[调试] SettingsService.DefaultWorkspaceID = %d\n", id)
	return id
}
