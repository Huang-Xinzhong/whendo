package services

import (
	"database/sql"
)

// WorkspaceStore defines storage operations for workspaces.
type WorkspaceStore interface {
	// placeholder for future methods.
	any
}

// WorkspaceService provides workspace business logic.
type WorkspaceService struct {
	db *sql.DB
}

// NewWorkspaceService creates a new WorkspaceService.
func NewWorkspaceService(db *sql.DB) *WorkspaceService {
	return &WorkspaceService{db: db}
}
