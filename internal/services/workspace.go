package services

import (
	"database/sql"
	"fmt"
	"strings"

	"whendo/internal/database"
	"whendo/internal/models"
	"whendo/internal/requests"
)

var workspaceColors = []string{"blue", "green", "purple"}

// WorkspaceService provides workspace business logic.
type WorkspaceService struct {
	db    *sql.DB
	store *database.WorkspaceStore
}

// NewWorkspaceService creates a new WorkspaceService.
func NewWorkspaceService(db *sql.DB) *WorkspaceService {
	return &WorkspaceService{db: db, store: database.NewWorkspaceStore(db)}
}

// List returns all workspaces enriched with color and task count.
func (s *WorkspaceService) List() ([]models.Workspace, error) {
	list, err := s.store.List()
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}

	counts, err := s.taskCounts()
	if err != nil {
		return nil, fmt.Errorf("task counts: %w", err)
	}

	for i := range list {
		list[i].Color = workspaceColors[i%len(workspaceColors)]
		list[i].TaskCount = counts[list[i].ID]
	}
	return list, nil
}

func (s *WorkspaceService) taskCounts() (map[int64]int, error) {
	rows, err := s.db.Query(`SELECT workspace_id, COUNT(*) FROM tasks GROUP BY workspace_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[int64]int)
	for rows.Next() {
		var wsID int64
		var c int
		if err := rows.Scan(&wsID, &c); err != nil {
			return nil, err
		}
		m[wsID] = c
	}
	return m, rows.Err()
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
