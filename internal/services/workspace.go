package services

import (
	"database/sql"
	"fmt"
	"strings"

	"whendo/internal/database"
	"whendo/internal/models"
	"whendo/internal/requests"
)

// WorkspaceService provides workspace business logic.
type WorkspaceService struct {
	store *database.WorkspaceStore
}

// NewWorkspaceService creates a new WorkspaceService.
func NewWorkspaceService(db *sql.DB) *WorkspaceService {
	return &WorkspaceService{store: database.NewWorkspaceStore(db)}
}

// List returns all workspaces.
func (s *WorkspaceService) List() ([]models.Workspace, error) {
	list, err := s.store.List()
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	return list, nil
}

// Create validates and creates a workspace.
func (s *WorkspaceService) Create(req requests.WorkspaceCreateReq) (*models.Workspace, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 50 {
		return nil, fmt.Errorf("invalid workspace name")
	}
	ws, err := s.store.Create(name)
	if err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}
	return ws, nil
}

// Update validates and renames a workspace.
func (s *WorkspaceService) Update(req requests.WorkspaceUpdateReq) (*models.Workspace, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 50 {
		return nil, fmt.Errorf("invalid workspace name")
	}
	ws, err := s.store.Update(req.ID, name)
	if err != nil {
		return nil, fmt.Errorf("update workspace: %w", err)
	}
	return ws, nil
}

// Delete removes a workspace.
func (s *WorkspaceService) Delete(id int64) error {
	if err := s.store.Delete(id); err != nil {
		return fmt.Errorf("delete workspace: %w", err)
	}
	return nil
}
