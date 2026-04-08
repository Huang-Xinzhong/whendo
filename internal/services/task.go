package services

import (
	"database/sql"
)

// TaskStore defines storage operations for tasks.
type TaskStore interface {
	// placeholder for future methods.
	any
}

// TaskService provides task business logic.
type TaskService struct {
	db *sql.DB
}

// NewTaskService creates a new TaskService.
func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{db: db}
}
