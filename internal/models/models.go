package models

import "time"

// TaskType defines the kind of a task.
type TaskType string

const (
	TaskTypeTodo     TaskType = "todo"
	TaskTypeReminder TaskType = "reminder"
)

// IntervalUnit defines the unit for reminder intervals.
type IntervalUnit string

const (
	IntervalUnitMinute IntervalUnit = "minute"
	IntervalUnitHour   IntervalUnit = "hour"
)

// RepeatMode defines how a reminder repeats.
type RepeatMode string

const (
	RepeatModeDaily   RepeatMode = "daily"
	RepeatModeWorkday RepeatMode = "workday"
	RepeatModeWeekly  RepeatMode = "weekly"
	RepeatModeMonthly RepeatMode = "monthly"
)

// Workspace represents a task container.
type Workspace struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Runtime enriched fields (not stored in DB).
	Color     string `json:"color,omitempty"`
	TaskCount int    `json:"taskCount,omitempty"`
}

// Task represents a todo item or reminder.
type Task struct {
	ID             int64        `json:"id"`
	WorkspaceID    int64        `json:"workspace_id"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Type           TaskType     `json:"type"`
	DueAt          *time.Time   `json:"due_at,omitempty"`
	RemindAt       *time.Time   `json:"remind_at,omitempty"`
	IsCompleted    bool         `json:"is_completed"`
	StartTime      string       `json:"start_time,omitempty"`
	EndTime        string       `json:"end_time,omitempty"`
	IntervalValue  int          `json:"interval_value,omitempty"`
	IntervalUnit   IntervalUnit `json:"interval_unit,omitempty"`
	RepeatMode     RepeatMode   `json:"repeat_mode,omitempty"`
	Weekdays       string       `json:"weekdays,omitempty"`
	MonthDay       int          `json:"month_day,omitempty"`
	NextTriggerAt  *time.Time   `json:"next_trigger_at,omitempty"`
	PausedDate     *time.Time   `json:"paused_date,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	// Runtime enriched fields (not stored in DB).
	RemindText  string `json:"remindText,omitempty"`
	PausedToday bool   `json:"pausedToday,omitempty"`
}
