package database

import (
	"database/sql"
	"fmt"
	"time"

	"whendo/internal/models"
)

// TaskStore provides task storage operations.
type TaskStore struct {
	db *sql.DB
}

// NewTaskStore creates a new TaskStore.
func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{db: db}
}

func scanTask(rows *sql.Rows) (models.Task, error) {
	var t models.Task
	var dueAt, remindAt, nextTriggerAt, pausedDate sql.NullTime
	err := rows.Scan(
		&t.ID,
		&t.WorkspaceID,
		&t.Title,
		&t.Description,
		&t.Type,
		&dueAt,
		&remindAt,
		&t.IsCompleted,
		&t.StartTime,
		&t.EndTime,
		&t.IntervalValue,
		&t.IntervalUnit,
		&t.RepeatMode,
		&t.Weekdays,
		&t.MonthDay,
		&nextTriggerAt,
		&pausedDate,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if dueAt.Valid {
		t.DueAt = &dueAt.Time
	}
	if remindAt.Valid {
		t.RemindAt = &remindAt.Time
	}
	if nextTriggerAt.Valid {
		t.NextTriggerAt = &nextTriggerAt.Time
	}
	if pausedDate.Valid {
		t.PausedDate = &pausedDate.Time
	}
	return t, err
}

func scanTaskRow(row *sql.Row) (models.Task, error) {
	var t models.Task
	var dueAt, remindAt, nextTriggerAt, pausedDate sql.NullTime
	err := row.Scan(
		&t.ID,
		&t.WorkspaceID,
		&t.Title,
		&t.Description,
		&t.Type,
		&dueAt,
		&remindAt,
		&t.IsCompleted,
		&t.StartTime,
		&t.EndTime,
		&t.IntervalValue,
		&t.IntervalUnit,
		&t.RepeatMode,
		&t.Weekdays,
		&t.MonthDay,
		&nextTriggerAt,
		&pausedDate,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if dueAt.Valid {
		t.DueAt = &dueAt.Time
	}
	if remindAt.Valid {
		t.RemindAt = &remindAt.Time
	}
	if nextTriggerAt.Valid {
		t.NextTriggerAt = &nextTriggerAt.Time
	}
	if pausedDate.Valid {
		t.PausedDate = &pausedDate.Time
	}
	return t, err
}

// List returns tasks for a workspace with optional filtering.
func (s *TaskStore) List(workspaceID int64, filter string) ([]models.Task, error) {
	query := `SELECT id, workspace_id, title, description, type, due_at, remind_at, is_completed, start_time, end_time, interval_value, interval_unit, repeat_mode, weekdays, month_day, next_trigger_at, paused_date, created_at, updated_at FROM tasks WHERE workspace_id = ?`
	args := []interface{}{workspaceID}

	switch filter {
	case "pending":
		query += ` AND is_completed = 0`
	case "completed":
		query += ` AND is_completed = 1`
	}

	query += ` ORDER BY created_at DESC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var list []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

// Get returns a single task by id.
func (s *TaskStore) Get(id int64) (*models.Task, error) {
	row := s.db.QueryRow(`SELECT id, workspace_id, title, description, type, due_at, remind_at, is_completed, start_time, end_time, interval_value, interval_unit, repeat_mode, weekdays, month_day, next_trigger_at, paused_date, created_at, updated_at FROM tasks WHERE id = ?`, id)
	t, err := scanTaskRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task %d not found", id)
		}
		return nil, fmt.Errorf("scan task: %w", err)
	}
	return &t, nil
}

// Create inserts a new task.
func (s *TaskStore) Create(t *models.Task) (*models.Task, error) {
	res, err := s.db.Exec(
		`INSERT INTO tasks (workspace_id, title, description, type, due_at, remind_at, is_completed, start_time, end_time, interval_value, interval_unit, repeat_mode, weekdays, month_day, next_trigger_at, paused_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.WorkspaceID, t.Title, t.Description, t.Type,
		t.DueAt, t.RemindAt, t.IsCompleted,
		t.StartTime, t.EndTime, t.IntervalValue, t.IntervalUnit,
		t.RepeatMode, t.Weekdays, t.MonthDay,
		t.NextTriggerAt, t.PausedDate,
	)
	if err != nil {
		return nil, fmt.Errorf("insert task: %w", err)
	}
	id, _ := res.LastInsertId()
	return s.Get(id)
}

// Update modifies an existing task.
func (s *TaskStore) Update(t *models.Task) (*models.Task, error) {
	_, err := s.db.Exec(
		`UPDATE tasks SET workspace_id = ?, title = ?, description = ?, type = ?, due_at = ?, remind_at = ?, is_completed = ?, start_time = ?, end_time = ?, interval_value = ?, interval_unit = ?, repeat_mode = ?, weekdays = ?, month_day = ?, next_trigger_at = ?, paused_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		t.WorkspaceID, t.Title, t.Description, t.Type,
		t.DueAt, t.RemindAt, t.IsCompleted,
		t.StartTime, t.EndTime, t.IntervalValue, t.IntervalUnit,
		t.RepeatMode, t.Weekdays, t.MonthDay,
		t.NextTriggerAt, t.PausedDate,
		t.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	return s.Get(t.ID)
}

// Delete removes a task.
func (s *TaskStore) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

// ToggleCompleted flips the completion status.
func (s *TaskStore) ToggleCompleted(id int64) (*models.Task, error) {
	_, err := s.db.Exec(`UPDATE tasks SET is_completed = NOT is_completed, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("toggle completed: %w", err)
	}
	return s.Get(id)
}

// TogglePause sets or clears paused_date to today.
func (s *TaskStore) TogglePause(id int64) (*models.Task, error) {
	t, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	today := time.Now().Truncate(24 * time.Hour)
	var pausedDate *time.Time
	if t.PausedDate == nil || !t.PausedDate.Equal(today) {
		pausedDate = &today
	}

	_, err = s.db.Exec(`UPDATE tasks SET paused_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, pausedDate, id)
	if err != nil {
		return nil, fmt.Errorf("toggle pause: %w", err)
	}
	return s.Get(id)
}

// ListPendingReminders returns active reminders that are not paused today and have next_trigger_at <= now.
func (s *TaskStore) ListPendingReminders(now time.Time) ([]models.Task, error) {
	todayStr := now.Format("2006-01-02")
	rows, err := s.db.Query(
		`SELECT id, workspace_id, title, description, type, due_at, remind_at, is_completed, start_time, end_time, interval_value, interval_unit, repeat_mode, weekdays, month_day, next_trigger_at, paused_date, created_at, updated_at FROM tasks WHERE type = 'reminder' AND (paused_date IS NULL OR paused_date != ?) AND next_trigger_at IS NOT NULL AND next_trigger_at <= ?`,
		todayStr, now,
	)
	if err != nil {
		return nil, fmt.Errorf("query reminders: %w", err)
	}
	defer rows.Close()

	var list []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan reminder: %w", err)
		}
		list = append(list, t)
	}
	return list, rows.Err()
}

// ClearCompleted removes all completed tasks.
func (s *TaskStore) ClearCompleted() error {
	_, err := s.db.Exec(`DELETE FROM tasks WHERE is_completed = 1`)
	if err != nil {
		return fmt.Errorf("clear completed: %w", err)
	}
	return nil
}
