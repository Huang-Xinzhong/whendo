package database

import (
	"database/sql"
	"fmt"
	"time"

	"whendo/internal/models"
)

// debugLog 格式化调试输出的辅助函数。
func debugLog(format string, args ...interface{}) {
	fmt.Printf("[调试] "+format+"\n", args...)
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
	var dueAt, remindAt, nextTriggerAt, pausedDate, pausedUntil, skipFrom, skipUntil sql.NullTime
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
		&pausedUntil,
		&skipFrom,
		&skipUntil,
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
	if pausedUntil.Valid {
		t.PausedUntil = &pausedUntil.Time
	}
	if skipFrom.Valid {
		t.SkipFrom = &skipFrom.Time
	}
	if skipUntil.Valid {
		t.SkipUntil = &skipUntil.Time
	}
	return t, err
}

func scanTaskRow(row *sql.Row) (models.Task, error) {
	var t models.Task
	var dueAt, remindAt, nextTriggerAt, pausedDate, pausedUntil, skipFrom, skipUntil sql.NullTime
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
		&pausedUntil,
		&skipFrom,
		&skipUntil,
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
	if pausedUntil.Valid {
		t.PausedUntil = &pausedUntil.Time
	}
	if skipFrom.Valid {
		t.SkipFrom = &skipFrom.Time
	}
	if skipUntil.Valid {
		t.SkipUntil = &skipUntil.Time
	}
	return t, err
}

// selectAllColumns 返回全量 SELECT 列列表（保持顺序与 scanTask 一致）。
func selectAllColumns() string {
	return `id, workspace_id, title, description, type, due_at, remind_at, is_completed, start_time, end_time, interval_value, interval_unit, repeat_mode, weekdays, month_day, next_trigger_at, paused_date, paused_until, skip_from, skip_until, created_at, updated_at`
}

// List 返回指定工作区下的任务列表，支持可选过滤。
func (s *TaskStore) List(workspaceID int64, filter string) ([]models.Task, error) {
	debugLog("TaskStore.List 查询任务: workspaceID=%d filter=%s", workspaceID, filter)
	query := `SELECT ` + selectAllColumns() + ` FROM tasks WHERE workspace_id = ?`
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
		debugLog("TaskStore.List 查询失败: %v", err)
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var list []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			debugLog("TaskStore.List 扫描失败: %v", err)
			return nil, fmt.Errorf("scan task: %w", err)
		}
		list = append(list, t)
	}
	debugLog("TaskStore.List 成功，共 %d 条任务", len(list))
	return list, rows.Err()
}

// Get 根据 ID 返回单个任务。
func (s *TaskStore) Get(id int64) (*models.Task, error) {
	debugLog("TaskStore.Get 查询任务: id=%d", id)
	row := s.db.QueryRow(`SELECT `+selectAllColumns()+` FROM tasks WHERE id = ?`, id)
	t, err := scanTaskRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			debugLog("TaskStore.Get 任务不存在: id=%d", id)
			return nil, fmt.Errorf("task %d not found", id)
		}
		debugLog("TaskStore.Get 扫描失败: %v", err)
		return nil, fmt.Errorf("scan task: %w", err)
	}
	debugLog("TaskStore.Get 成功: id=%d type=%s next_trigger_at=%v", t.ID, t.Type, t.NextTriggerAt)
	return &t, nil
}

// Create 插入一条新任务。
func (s *TaskStore) Create(t *models.Task) (*models.Task, error) {
	debugLog("TaskStore.Create 插入任务: title=%s type=%s remind_at=%v next_trigger_at=%v", t.Title, t.Type, t.RemindAt, t.NextTriggerAt)
	res, err := s.db.Exec(
		`INSERT INTO tasks (workspace_id, title, description, type, due_at, remind_at, is_completed, start_time, end_time, interval_value, interval_unit, repeat_mode, weekdays, month_day, next_trigger_at, paused_date, paused_until, skip_from, skip_until) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.WorkspaceID, t.Title, t.Description, t.Type,
		t.DueAt, t.RemindAt, t.IsCompleted,
		t.StartTime, t.EndTime, t.IntervalValue, t.IntervalUnit,
		t.RepeatMode, t.Weekdays, t.MonthDay,
		t.NextTriggerAt, t.PausedDate, t.PausedUntil, t.SkipFrom, t.SkipUntil,
	)
	if err != nil {
		debugLog("TaskStore.Create 插入失败: %v", err)
		return nil, fmt.Errorf("insert task: %w", err)
	}
	id, _ := res.LastInsertId()
	debugLog("TaskStore.Create 成功: id=%d", id)
	return s.Get(id)
}

// Update 修改已有任务。
func (s *TaskStore) Update(t *models.Task) (*models.Task, error) {
	_, err := s.db.Exec(
		`UPDATE tasks SET workspace_id = ?, title = ?, description = ?, type = ?, due_at = ?, remind_at = ?, is_completed = ?, start_time = ?, end_time = ?, interval_value = ?, interval_unit = ?, repeat_mode = ?, weekdays = ?, month_day = ?, next_trigger_at = ?, paused_date = ?, paused_until = ?, skip_from = ?, skip_until = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		t.WorkspaceID, t.Title, t.Description, t.Type,
		t.DueAt, t.RemindAt, t.IsCompleted,
		t.StartTime, t.EndTime, t.IntervalValue, t.IntervalUnit,
		t.RepeatMode, t.Weekdays, t.MonthDay,
		t.NextTriggerAt, t.PausedDate, t.PausedUntil, t.SkipFrom, t.SkipUntil,
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

// SetPauseUntil 设置/清除暂停截止时间。
func (s *TaskStore) SetPauseUntil(id int64, until *time.Time) (*models.Task, error) {
	_, err := s.db.Exec(`UPDATE tasks SET paused_until = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, until, id)
	if err != nil {
		return nil, fmt.Errorf("set pause until: %w", err)
	}
	return s.Get(id)
}

// SetSkipRange 设置跳过区间。
func (s *TaskStore) SetSkipRange(id int64, from, until *time.Time) (*models.Task, error) {
	_, err := s.db.Exec(`UPDATE tasks SET skip_from = ?, skip_until = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, from, until, id)
	if err != nil {
		return nil, fmt.Errorf("set skip range: %w", err)
	}
	return s.Get(id)
}

// ClearSkipRange 清除跳过区间。
func (s *TaskStore) ClearSkipRange(id int64) (*models.Task, error) {
	_, err := s.db.Exec(`UPDATE tasks SET skip_from = NULL, skip_until = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("clear skip range: %w", err)
	}
	return s.Get(id)
}

// ListPendingReminders 返回待触发的 reminder 任务（不含 todo）。
func (s *TaskStore) ListPendingReminders(now time.Time) ([]models.Task, error) {
	nowStr := now.Format("2006-01-02 15:04:05")
	debugLog("TaskStore.ListPendingReminders 查询待触发提醒: now=%v nowStr=%s", now, nowStr)
	rows, err := s.db.Query(
		`SELECT `+selectAllColumns()+`
		FROM tasks
		WHERE type = 'reminder'
		  AND is_completed = 0
		  AND next_trigger_at IS NOT NULL
		  AND datetime(next_trigger_at) <= datetime(?)
		  AND (paused_until IS NULL OR datetime(?) < datetime(paused_until))
		  AND NOT (
				skip_from IS NOT NULL AND skip_until IS NOT NULL
				AND datetime(?) >= datetime(skip_from)
				AND datetime(?) < datetime(skip_until)
			)`,
		nowStr, nowStr, nowStr, nowStr,
	)
	if err != nil {
		debugLog("TaskStore.ListPendingReminders 查询失败: %v", err)
		return nil, fmt.Errorf("query reminders: %w", err)
	}
	defer rows.Close()

	var list []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			debugLog("TaskStore.ListPendingReminders 扫描失败: %v", err)
			return nil, fmt.Errorf("scan reminder: %w", err)
		}
		list = append(list, t)
	}
	debugLog("TaskStore.ListPendingReminders 成功，共 %d 条", len(list))
	return list, rows.Err()
}

// ListPendingTodos 返回待触发的 todo 任务。
func (s *TaskStore) ListPendingTodos(now time.Time) ([]models.Task, error) {
	nowStr := now.Format("2006-01-02 15:04:05")
	debugLog("TaskStore.ListPendingTodos 查询待触发待办: now=%v nowStr=%s", now, nowStr)
	rows, err := s.db.Query(
		`SELECT `+selectAllColumns()+`
		FROM tasks
		WHERE type = 'todo'
		  AND is_completed = 0
		  AND next_trigger_at IS NOT NULL
		  AND datetime(next_trigger_at) <= datetime(?)`,
		nowStr,
	)
	if err != nil {
		debugLog("TaskStore.ListPendingTodos 查询失败: %v", err)
		return nil, fmt.Errorf("query todos: %w", err)
	}
	defer rows.Close()

	var list []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			debugLog("TaskStore.ListPendingTodos 扫描失败: %v", err)
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		list = append(list, t)
	}
	debugLog("TaskStore.ListPendingTodos 成功，共 %d 条", len(list))
	return list, rows.Err()
}

// ClearExpiredPauses 将已过期的 paused_until 置为 NULL。
func (s *TaskStore) ClearExpiredPauses(now time.Time) error {
	nowStr := now.Format("2006-01-02 15:04:05")
	res, err := s.db.Exec(
		`UPDATE tasks SET paused_until = NULL, updated_at = CURRENT_TIMESTAMP WHERE paused_until IS NOT NULL AND datetime(paused_until) <= datetime(?)`,
		nowStr,
	)
	if err != nil {
		return fmt.Errorf("clear expired pauses: %w", err)
	}
	if n, _ := res.RowsAffected(); n > 0 {
		debugLog("TaskStore.ClearExpiredPauses 清理 %d 条过期暂停", n)
	}
	return nil
}

// ClearExpiredSkips 将已过期的 skip 区间置为 NULL。
func (s *TaskStore) ClearExpiredSkips(now time.Time) error {
	nowStr := now.Format("2006-01-02 15:04:05")
	res, err := s.db.Exec(
		`UPDATE tasks SET skip_from = NULL, skip_until = NULL, updated_at = CURRENT_TIMESTAMP WHERE skip_until IS NOT NULL AND datetime(skip_until) <= datetime(?)`,
		nowStr,
	)
	if err != nil {
		return fmt.Errorf("clear expired skips: %w", err)
	}
	if n, _ := res.RowsAffected(); n > 0 {
		debugLog("TaskStore.ClearExpiredSkips 清理 %d 条过期跳过区间", n)
	}
	return nil
}

// ClearCompleted 删除所有已完成的任务。
func (s *TaskStore) ClearCompleted() error {
	_, err := s.db.Exec(`DELETE FROM tasks WHERE is_completed = 1`)
	if err != nil {
		return fmt.Errorf("clear completed: %w", err)
	}
	return nil
}
