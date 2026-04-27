package models

import "time"

// TaskType 定义任务的类型。
type TaskType string

const (
	TaskTypeTodo     TaskType = "todo"
	TaskTypeReminder TaskType = "reminder"
)

// IntervalUnit 定义提醒间隔的单位。
type IntervalUnit string

const (
	IntervalUnitMinute IntervalUnit = "minute"
	IntervalUnitHour   IntervalUnit = "hour"
)

// RepeatMode 定义提醒的重复模式。
type RepeatMode string

const (
	RepeatModeDaily   RepeatMode = "daily"
	RepeatModeWorkday RepeatMode = "workday"
	RepeatModeWeekly  RepeatMode = "weekly"
	RepeatModeMonthly RepeatMode = "monthly"
)

// Workspace 表示一个任务容器（工作区）。
type Workspace struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 运行时填充的字段（不存储在数据库中）。
	Color     string `json:"color,omitempty"`
	TaskCount int    `json:"taskCount,omitempty"`
}

// Task 表示一个待办项或定时提醒。
type Task struct {
	ID            int64        `json:"id"`
	WorkspaceID   int64        `json:"workspace_id"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Type          TaskType     `json:"type"`
	DueAt         *time.Time   `json:"due_at,omitempty"`
	RemindAt      *time.Time   `json:"remind_at,omitempty"`
	IsCompleted   bool         `json:"is_completed"`
	StartTime     string       `json:"start_time,omitempty"`
	EndTime       string       `json:"end_time,omitempty"`
	IntervalValue int          `json:"interval_value,omitempty"`
	IntervalUnit  IntervalUnit `json:"interval_unit,omitempty"`
	RepeatMode    RepeatMode   `json:"repeat_mode,omitempty"`
	Weekdays      string       `json:"weekdays,omitempty"`
	MonthDay      int          `json:"month_day,omitempty"`
	NextTriggerAt *time.Time   `json:"next_trigger_at,omitempty"`
	PausedDate    *time.Time   `json:"paused_date,omitempty"`
	PausedUntil   *time.Time   `json:"paused_until,omitempty"`
	SkipFrom      *time.Time   `json:"skip_from,omitempty"`
	SkipUntil     *time.Time   `json:"skip_until,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	// 运行时填充的字段（不存储在数据库中）。
	RemindText    string `json:"remindText,omitempty"`
	PausedToday   bool   `json:"pausedToday,omitempty"`
	IsPaused      bool   `json:"isPaused,omitempty"`
	InSkipRange   bool   `json:"inSkipRange,omitempty"`
	PausedText    string `json:"pausedText,omitempty"`
	SkipRangeText string `json:"skipRangeText,omitempty"`
}
