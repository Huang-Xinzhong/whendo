package requests

import (
	"time"

	"whendo/internal/models"
)

// WorkspaceCreateReq is the request for creating a workspace.
type WorkspaceCreateReq struct {
	Name string `json:"name"`
}

// WorkspaceUpdateReq is the request for updating a workspace.
type WorkspaceUpdateReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TaskCreateReq is the request for creating a task.
type TaskCreateReq struct {
	WorkspaceID   int64                `json:"workspace_id"`
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	Type          models.TaskType      `json:"type"`
	DueAt         *time.Time           `json:"due_at,omitempty"`
	RemindAt      *time.Time           `json:"remind_at,omitempty"`
	StartTime     string               `json:"start_time,omitempty"`
	EndTime       string               `json:"end_time,omitempty"`
	IntervalValue int                  `json:"interval_value,omitempty"`
	IntervalUnit  models.IntervalUnit  `json:"interval_unit,omitempty"`
	RepeatMode    models.RepeatMode    `json:"repeat_mode,omitempty"`
	Weekdays      string               `json:"weekdays,omitempty"`
	MonthDay      int                  `json:"month_day,omitempty"`
}

// TaskUpdateReq is the request for updating a task.
type TaskUpdateReq struct {
	ID            int64                `json:"id"`
	WorkspaceID   int64                `json:"workspace_id,omitempty"`
	Title         string               `json:"title,omitempty"`
	Description   string               `json:"description,omitempty"`
	Type          models.TaskType      `json:"type,omitempty"`
	DueAt         *time.Time           `json:"due_at,omitempty"`
	RemindAt      *time.Time           `json:"remind_at,omitempty"`
	StartTime     string               `json:"start_time,omitempty"`
	EndTime       string               `json:"end_time,omitempty"`
	IntervalValue int                  `json:"interval_value,omitempty"`
	IntervalUnit  models.IntervalUnit  `json:"interval_unit,omitempty"`
	RepeatMode    models.RepeatMode    `json:"repeat_mode,omitempty"`
	Weekdays      string               `json:"weekdays,omitempty"`
	MonthDay      int                  `json:"month_day,omitempty"`
}

// SettingsUpdateReq is the request for updating settings.
type SettingsUpdateReq struct {
	DefaultWorkspaceID int64  `json:"default_workspace_id,omitempty"`
	DefaultSort        string `json:"default_sort,omitempty"`
}
