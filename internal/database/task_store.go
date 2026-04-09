package database

import (
	"database/sql"
	"fmt"
	"time"

	"whendo/internal/models"
)

// [DEBUG] 格式化调试输出的辅助函数。
func debugLog(format string, args ...interface{}) {
	fmt.Printf("[DEBUG] "+format+"\n", args...)
}

// TaskStore 提供任务相关的存储操作。
type TaskStore struct {
	db *sql.DB
}

// NewTaskStore 创建一个新的 TaskStore 实例。
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

// List 返回指定工作区下的任务列表，支持可选过滤。
func (s *TaskStore) List(workspaceID int64, filter string) ([]models.Task, error) {
	debugLog("TaskStore.List workspaceID=%d filter=%s", workspaceID, filter)
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
		debugLog("TaskStore.List query error: %v", err)
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var list []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			debugLog("TaskStore.List scan error: %v", err)
			return nil, fmt.Errorf("scan task: %w", err)
		}
		list = append(list, t)
	}
	debugLog("TaskStore.List returned %d tasks", len(list))
	return list, rows.Err()
}

// Get 根据 ID 返回单个任务。
func (s *TaskStore) Get(id int64) (*models.Task, error) {
	debugLog("TaskStore.Get id=%d", id)
	row := s.db.QueryRow(`SELECT id, workspace_id, title, description, type, due_at, remind_at, is_completed, start_time, end_time, interval_value, interval_unit, repeat_mode, weekdays, month_day, next_trigger_at, paused_date, created_at, updated_at FROM tasks WHERE id = ?`, id)
	t, err := scanTaskRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			debugLog("TaskStore.Get task %d not found", id)
			return nil, fmt.Errorf("task %d not found", id)
		}
		debugLog("TaskStore.Get scan error: %v", err)
		return nil, fmt.Errorf("scan task: %w", err)
	}
	debugLog("TaskStore.Get found task id=%d type=%s next_trigger_at=%v", t.ID, t.Type, t.NextTriggerAt)
	return &t, nil
}

// Create 插入一条新任务。
func (s *TaskStore) Create(t *models.Task) (*models.Task, error) {
	debugLog("TaskStore.Create title=%s type=%s remind_at=%v next_trigger_at=%v", t.Title, t.Type, t.RemindAt, t.NextTriggerAt)
	res, err := s.db.Exec(
		`INSERT INTO tasks (workspace_id, title, description, type, due_at, remind_at, is_completed, start_time, end_time, interval_value, interval_unit, repeat_mode, weekdays, month_day, next_trigger_at, paused_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.WorkspaceID, t.Title, t.Description, t.Type,
		t.DueAt, t.RemindAt, t.IsCompleted,
		t.StartTime, t.EndTime, t.IntervalValue, t.IntervalUnit,
		t.RepeatMode, t.Weekdays, t.MonthDay,
		t.NextTriggerAt, t.PausedDate,
	)
	if err != nil {
		debugLog("TaskStore.Create insert error: %v", err)
		return nil, fmt.Errorf("insert task: %w", err)
	}
	id, _ := res.LastInsertId()
	debugLog("TaskStore.Created id=%d", id)
	return s.Get(id)
}

// Update 修改已有任务。
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

// Delete 删除任务。
func (s *TaskStore) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

// ToggleCompleted 切换任务的完成状态。
func (s *TaskStore) ToggleCompleted(id int64) (*models.Task, error) {
	_, err := s.db.Exec(`UPDATE tasks SET is_completed = NOT is_completed, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("toggle completed: %w", err)
	}
	return s.Get(id)
}

// TogglePause 将 paused_date 设为今天或清空。
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

// ListPendingReminders 返回今天未暂停且 next_trigger_at <= now 的活跃提醒。
func (s *TaskStore) ListPendingReminders(now time.Time) ([]models.Task, error) {
	todayStr := now.Format("2006-01-02")
	nowStr := now.Format("2006-01-02 15:04:05")
	debugLog("TaskStore.ListPendingReminders now=%v nowStr=%s todayStr=%s", now, nowStr, todayStr)
	rows, err := s.db.Query(
		`SELECT
			id,
			workspace_id,
			title,
			description,
			type,
			due_at,
			remind_at,
			is_completed,
			start_time,
			end_time,
			interval_value,
			interval_unit,
			repeat_mode,
			weekdays,
			month_day,
			next_trigger_at,
			paused_date,
			created_at,
			updated_at
		FROM tasks
		WHERE type IN ('reminder', 'todo')
		  AND (paused_date IS NULL OR datetime(paused_date) != datetime(?))
		  AND next_trigger_at IS NOT NULL
		  AND datetime(next_trigger_at) <= datetime(?)`,
		todayStr, nowStr,
	)
	if err != nil {
		debugLog("TaskStore.ListPendingReminders query error: %v", err)
		return nil, fmt.Errorf("query reminders: %w", err)
	}
	defer rows.Close()

	var list []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			debugLog("TaskStore.ListPendingReminders scan error: %v", err)
			return nil, fmt.Errorf("scan reminder: %w", err)
		}
		list = append(list, t)
	}
	debugLog("TaskStore.ListPendingReminders returned %d tasks", len(list))
	return list, rows.Err()
}

// ClearCompleted 删除所有已完成的任务。
func (s *TaskStore) ClearCompleted() error {
	_, err := s.db.Exec(`DELETE FROM tasks WHERE is_completed = 1`)
	if err != nil {
		return fmt.Errorf("clear completed: %w", err)
	}
	return nil
}
