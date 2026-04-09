package services

import (
	"database/sql"
	"fmt"
	"strconv"

	"whendo/internal/database"
	"whendo/internal/requests"
)

// SettingsService provides settings business logic.
type SettingsService struct {
	store *database.SettingsStore
}

// NewSettingsService creates a new SettingsService.
func NewSettingsService(db *sql.DB) *SettingsService {
	s := &SettingsService{store: database.NewSettingsStore(db)}
	_ = s.store.Open()
	return s
}

// Get returns all settings as a map.
func (s *SettingsService) Get() (map[string]string, error) {
	m, err := s.store.All()
	if err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}
	return m, nil
}

// Update updates settings.
func (s *SettingsService) Update(req requests.SettingsUpdateReq) error {
	if req.DefaultWorkspaceID != 0 {
		if err := s.store.SetInt("default_workspace_id", req.DefaultWorkspaceID); err != nil {
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

// DefaultWorkspaceID returns the default workspace id.
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
