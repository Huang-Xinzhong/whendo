package requests

import (
	"whendo/internal/models"
)

// WorkspaceCreateReq 是创建工作区的请求。
type WorkspaceCreateReq struct {
	Name string `json:"name"`
}

// WorkspaceUpdateReq 是更新工作区的请求。
type WorkspaceUpdateReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TaskCreateReq 是创建任务的请求。
type TaskCreateReq struct {
	WorkspaceID   int64               `json:"workspace_id"`
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	Type          models.TaskType     `json:"type"`
	DueAt         *string             `json:"due_at,omitempty"`
	RemindAt      *string             `json:"remind_at,omitempty"`
	StartTime     string              `json:"start_time,omitempty"`
	EndTime       string              `json:"end_time,omitempty"`
	IntervalValue int                 `json:"interval_value,omitempty"`
	IntervalUnit  models.IntervalUnit `json:"interval_unit,omitempty"`
	RepeatMode    models.RepeatMode   `json:"repeat_mode,omitempty"`
	Weekdays      string              `json:"weekdays,omitempty"`
	MonthDay      int                 `json:"month_day,omitempty"`
	PausedDate    *string             `json:"paused_date,omitempty"`
	PausedUntil   *string             `json:"paused_until,omitempty"`
	SkipFrom      *string             `json:"skip_from,omitempty"`
	SkipUntil     *string             `json:"skip_until,omitempty"`
}

// TaskUpdateReq 是更新任务的请求。
type TaskUpdateReq struct {
	ID            int64               `json:"id"`
	WorkspaceID   int64               `json:"workspace_id,omitempty"`
	Title         string              `json:"title,omitempty"`
	Description   string              `json:"description,omitempty"`
	Type          models.TaskType     `json:"type,omitempty"`
	DueAt         *string             `json:"due_at,omitempty"`
	RemindAt      *string             `json:"remind_at,omitempty"`
	StartTime     string              `json:"start_time,omitempty"`
	EndTime       string              `json:"end_time,omitempty"`
	IntervalValue int                 `json:"interval_value,omitempty"`
	IntervalUnit  models.IntervalUnit `json:"interval_unit,omitempty"`
	RepeatMode    models.RepeatMode   `json:"repeat_mode,omitempty"`
	Weekdays      string              `json:"weekdays,omitempty"`
	MonthDay      int                 `json:"month_day,omitempty"`
	PausedDate    *string             `json:"paused_date,omitempty"`
	PausedUntil   *string             `json:"paused_until,omitempty"`
	SkipFrom      *string             `json:"skip_from,omitempty"`
	SkipUntil     *string             `json:"skip_until,omitempty"`
}

// SettingsUpdateReq 是更新设置的请求。
type SettingsUpdateReq struct {
	DefaultWorkspaceID string `json:"default_workspace_id,omitempty"`
	DefaultSort        string `json:"default_sort,omitempty"`
}
