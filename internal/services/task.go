package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"whendo/internal/database"
	"whendo/internal/models"
	"whendo/internal/requests"
)

// TaskService 提供任务相关的业务逻辑。
type TaskService struct {
	store *database.TaskStore
}

// NewTaskService 创建一个新的 TaskService 实例。
func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{store: database.NewTaskStore(db)}
}

// List 返回指定工作区下的任务列表。
func (s *TaskService) List(workspaceID int64, filter string) ([]models.Task, error) {
	fmt.Printf("[调试] TaskService.List 获取任务列表: workspaceID=%d filter=%s\n", workspaceID, filter)
	list, err := s.store.List(workspaceID, filter)
	if err != nil {
		fmt.Printf("[调试] TaskService.List 失败: %v\n", err)
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	for i := range list {
		enrichTask(&list[i])
	}
	fmt.Printf("[调试] TaskService.List 成功，共 %d 条任务\n", len(list))
	return list, nil
}

// Get 返回单个任务。
func (s *TaskService) Get(id int64) (*models.Task, error) {
	t, err := s.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	enrichTask(t)
	return t, nil
}

// Create 校验并创建任务。
func (s *TaskService) Create(req requests.TaskCreateReq) (*models.Task, error) {
	fmt.Printf("[调试] TaskService.Create 创建任务: req=%+v\n", req)
	if err := validateTaskCreate(req); err != nil {
		fmt.Printf("[调试] TaskService.Create 校验失败: %v\n", err)
		return nil, err
	}
	dueAt, err := parseOptionalTime(req.DueAt)
	if err != nil {
		return nil, err
	}
	remindAt, err := parseOptionalTime(req.RemindAt)
	if err != nil {
		return nil, err
	}
	pausedDate, err := parseOptionalTime(req.PausedDate)
	if err != nil {
		return nil, err
	}
	pausedUntil, err := parseOptionalTime(req.PausedUntil)
	if err != nil {
		return nil, err
	}
	skipFrom, err := parseOptionalTime(req.SkipFrom)
	if err != nil {
		return nil, err
	}
	skipUntil, err := parseOptionalTime(req.SkipUntil)
	if err != nil {
		return nil, err
	}

	t := models.Task{
		WorkspaceID:   req.WorkspaceID,
		Title:         strings.TrimSpace(req.Title),
		Description:   strings.TrimSpace(req.Description),
		Type:          req.Type,
		DueAt:         dueAt,
		RemindAt:      remindAt,
		PausedDate:    pausedDate,
		PausedUntil:   pausedUntil,
		SkipFrom:      skipFrom,
		SkipUntil:     skipUntil,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		IntervalValue: req.IntervalValue,
		IntervalUnit:  req.IntervalUnit,
		RepeatMode:    req.RepeatMode,
		Weekdays:      req.Weekdays,
		MonthDay:      req.MonthDay,
		IsCompleted:   false,
	}
	if t.Type == models.TaskTypeReminder {
		next, err := calcNextTrigger(t)
		if err != nil {
			return nil, fmt.Errorf("calc next trigger: %w", err)
		}
		t.NextTriggerAt = next
	} else if t.Type == models.TaskTypeTodo {
		t.NextTriggerAt = remindAt
	}
	fmt.Printf("[调试] TaskService.Create 入库前 next_trigger_at=%v\n", t.NextTriggerAt)
	created, err := s.store.Create(&t)
	if err != nil {
		fmt.Printf("[调试] TaskService.Create 入库失败: %v\n", err)
		return nil, fmt.Errorf("create task: %w", err)
	}
	enrichTask(created)
	fmt.Printf("[调试] TaskService.Create 成功: id=%d next_trigger_at=%v\n", created.ID, created.NextTriggerAt)
	return created, nil
}

// Update 校验并更新任务。
func (s *TaskService) Update(req requests.TaskUpdateReq) (*models.Task, error) {
	fmt.Printf("[调试] TaskService.Update 更新任务: req=%+v\n", req)
	existing, err := s.store.Get(req.ID)
	if err != nil {
		fmt.Printf("[调试] TaskService.Update 查询失败: %v\n", err)
		return nil, fmt.Errorf("get task: %w", err)
	}

	if req.WorkspaceID != 0 {
		existing.WorkspaceID = req.WorkspaceID
	}
	if req.Title != "" {
		existing.Title = strings.TrimSpace(req.Title)
	}
	existing.Description = strings.TrimSpace(req.Description)
	if req.Type != "" {
		existing.Type = req.Type
	}
	if req.DueAt != nil {
		dueAt, err := parseOptionalTime(req.DueAt)
		if err != nil {
			return nil, err
		}
		existing.DueAt = dueAt
	}
	if req.RemindAt != nil {
		remindAt, err := parseOptionalTime(req.RemindAt)
		if err != nil {
			return nil, err
		}
		existing.RemindAt = remindAt
	}
	if req.StartTime != "" {
		existing.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		existing.EndTime = req.EndTime
	}
	if req.IntervalValue != 0 {
		existing.IntervalValue = req.IntervalValue
	}
	if req.IntervalUnit != "" {
		existing.IntervalUnit = req.IntervalUnit
	}
	if req.RepeatMode != "" {
		existing.RepeatMode = req.RepeatMode
	}
	if req.Weekdays != "" {
		existing.Weekdays = req.Weekdays
	}
	if req.MonthDay != 0 {
		existing.MonthDay = req.MonthDay
	}
	if req.PausedDate != nil {
		pausedDate, err := parseOptionalTime(req.PausedDate)
		if err != nil {
			return nil, err
		}
		existing.PausedDate = pausedDate
	}
	if req.PausedUntil != nil {
		pausedUntil, err := parseOptionalTime(req.PausedUntil)
		if err != nil {
			return nil, err
		}
		existing.PausedUntil = pausedUntil
	}
	if req.SkipFrom != nil {
		skipFrom, err := parseOptionalTime(req.SkipFrom)
		if err != nil {
			return nil, err
		}
		existing.SkipFrom = skipFrom
	}
	if req.SkipUntil != nil {
		skipUntil, err := parseOptionalTime(req.SkipUntil)
		if err != nil {
			return nil, err
		}
		existing.SkipUntil = skipUntil
	}

	if err := validateTask(*existing); err != nil {
		return nil, err
	}

	if existing.Type == models.TaskTypeReminder {
		next, err := calcNextTrigger(*existing)
		if err != nil {
			return nil, fmt.Errorf("calc next trigger: %w", err)
		}
		existing.NextTriggerAt = next
	} else if existing.Type == models.TaskTypeTodo {
		existing.NextTriggerAt = existing.RemindAt
	} else {
		existing.NextTriggerAt = nil
	}
	fmt.Printf("[调试] TaskService.Update 入库前 id=%d next_trigger_at=%v\n", existing.ID, existing.NextTriggerAt)

	updated, err := s.store.Update(existing)
	if err != nil {
		fmt.Printf("[调试] TaskService.Update 入库失败: %v\n", err)
		return nil, fmt.Errorf("update task: %w", err)
	}
	enrichTask(updated)
	fmt.Printf("[调试] TaskService.Update 成功: id=%d next_trigger_at=%v\n", updated.ID, updated.NextTriggerAt)
	return updated, nil
}

// Delete 删除任务。
func (s *TaskService) Delete(id int64) error {
	if err := s.store.Delete(id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

// ToggleCompleted 切换任务的完成状态。
func (s *TaskService) ToggleCompleted(id int64) (*models.Task, error) {
	t, err := s.store.ToggleCompleted(id)
	if err != nil {
		return nil, fmt.Errorf("toggle completed: %w", err)
	}
	enrichTask(t)
	return t, nil
}

// TogglePause 切换今日暂停状态。
// 内部改用 paused_until 实现：设置到今晚 23:59。
func (s *TaskService) TogglePause(id int64) (*models.Task, error) {
	now := time.Now()
	tonight := now.Truncate(24 * time.Hour).Add(24*time.Hour - time.Second)
	return s.SetPauseUntil(id, &tonight)
}

// SetPauseUntil 设置/清除暂停截止时间。
func (s *TaskService) SetPauseUntil(id int64, until *time.Time) (*models.Task, error) {
	t, err := s.store.SetPauseUntil(id, until)
	if err != nil {
		return nil, fmt.Errorf("set pause until: %w", err)
	}
	// 取消暂停时重新计算触发时间。
	if until == nil && t.Type == models.TaskTypeReminder {
		next, err := calcNextTrigger(*t)
		if err != nil {
			return nil, fmt.Errorf("calc next trigger: %w", err)
		}
		t.NextTriggerAt = next
		if _, err := s.store.Update(t); err != nil {
			return nil, fmt.Errorf("update after unpause: %w", err)
		}
		t, _ = s.store.Get(t.ID)
	}
	enrichTask(t)
	return t, nil
}

// SetSkipRange 设置跳过区间。
func (s *TaskService) SetSkipRange(id int64, from, until *time.Time) (*models.Task, error) {
	t, err := s.store.SetSkipRange(id, from, until)
	if err != nil {
		return nil, fmt.Errorf("set skip range: %w", err)
	}
	enrichTask(t)
	return t, nil
}

// ClearSkipRange 清除跳过区间。
func (s *TaskService) ClearSkipRange(id int64) (*models.Task, error) {
	t, err := s.store.ClearSkipRange(id)
	if err != nil {
		return nil, fmt.Errorf("clear skip range: %w", err)
	}
	enrichTask(t)
	return t, nil
}

// ClearCompleted 删除所有已完成的任务。
func (s *TaskService) ClearCompleted() error {
	if err := s.store.ClearCompleted(); err != nil {
		return fmt.Errorf("clear completed: %w", err)
	}
	return nil
}

func validateTaskCreate(req requests.TaskCreateReq) error {
	t := models.Task{
		WorkspaceID:   req.WorkspaceID,
		Title:         req.Title,
		Type:          req.Type,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		IntervalValue: req.IntervalValue,
		IntervalUnit:  req.IntervalUnit,
		RepeatMode:    req.RepeatMode,
		Weekdays:      req.Weekdays,
		MonthDay:      req.MonthDay,
	}
	return validateTask(t)
}

func parseOptionalTime(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	// 前端传入的 ISO 字符串（RFC3339）
	if t, err := time.Parse(time.RFC3339, *s); err == nil {
		return &t, nil
	}
	// 带秒的 datetime-local
	if t, err := time.ParseInLocation("2006-01-02T15:04:05", *s, time.Local); err == nil {
		return &t, nil
	}
	// 不带秒的 datetime-local（浏览器默认）
	if t, err := time.ParseInLocation("2006-01-02T15:04", *s, time.Local); err == nil {
		return &t, nil
	}
	return nil, fmt.Errorf("invalid time format: %s", *s)
}

func validateTask(t models.Task) error {
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("task title is required")
	}
	if t.WorkspaceID == 0 {
		return fmt.Errorf("workspace_id is required")
	}
	if t.Type != models.TaskTypeTodo && t.Type != models.TaskTypeReminder {
		return fmt.Errorf("invalid task type")
	}
	if t.Type == models.TaskTypeReminder {
		if t.StartTime == "" || t.EndTime == "" {
			return fmt.Errorf("start_time and end_time are required for reminder")
		}
		if t.IntervalValue <= 0 {
			return fmt.Errorf("interval_value is required for reminder")
		}
		if t.IntervalUnit != models.IntervalUnitMinute && t.IntervalUnit != models.IntervalUnitHour {
			return fmt.Errorf("interval_unit is required for reminder")
		}
		if t.RepeatMode != models.RepeatModeDaily && t.RepeatMode != models.RepeatModeWorkday &&
			t.RepeatMode != models.RepeatModeWeekly && t.RepeatMode != models.RepeatModeMonthly {
			return fmt.Errorf("repeat_mode is required for reminder")
		}
	}
	return nil
}

func enrichTask(t *models.Task) {
	now := time.Now()

	// 填充 RemindText（仅限 reminder）
	if t.Type == models.TaskTypeReminder {
		t.RemindText = buildRemindText(t)
	}

	// 当日暂停（兼容旧 paused_date）
	todayStr := now.Format("2006-01-02")
	if t.PausedDate != nil && t.PausedDate.Format("2006-01-02") == todayStr {
		t.PausedToday = true
	} else {
		t.PausedToday = false
	}

	// 持续暂停
	if t.PausedUntil != nil && t.PausedUntil.After(now) {
		t.IsPaused = true
		if t.PausedUntil.Year() >= 9999 {
			t.PausedText = "持续暂停"
		} else {
			t.PausedText = fmt.Sprintf("暂停到 %s", t.PausedUntil.Format("2006-01-02 15:04"))
		}
	} else {
		t.IsPaused = false
		t.PausedText = ""
	}

	// 跳过区间
	if t.SkipFrom != nil && t.SkipUntil != nil && now.After(*t.SkipFrom) && now.Before(*t.SkipUntil) {
		t.InSkipRange = true
		t.SkipRangeText = fmt.Sprintf("跳过到 %s", t.SkipUntil.Format("2006-01-02 15:04"))
	} else {
		t.InSkipRange = false
		t.SkipRangeText = ""
	}
}

func buildRemindText(t *models.Task) string {
	if t.Type != models.TaskTypeReminder {
		return ""
	}
	var intervalLabel string
	if t.IntervalUnit == models.IntervalUnitHour {
		if t.IntervalValue == 1 {
			intervalLabel = "每小时"
		} else {
			intervalLabel = fmt.Sprintf("每 %d 小时", t.IntervalValue)
		}
	} else {
		intervalLabel = fmt.Sprintf("每 %d 分钟", t.IntervalValue)
	}

	switch t.RepeatMode {
	case models.RepeatModeDaily:
		return fmt.Sprintf("%s提醒", intervalLabel)
	case models.RepeatModeWorkday:
		return fmt.Sprintf("每个工作日 %s 提醒", t.StartTime)
	case models.RepeatModeWeekly:
		var days []string
		for _, d := range strings.Split(t.Weekdays, ",") {
			idx := 0
			fmt.Sscanf(d, "%d", &idx)
			if idx >= 1 && idx <= 7 {
				days = append(days, "周"+"一二三四五六日"[idx-1:idx])
			}
		}
		return fmt.Sprintf("每周%s %s 提醒", strings.Join(days, " / "), t.StartTime)
	case models.RepeatModeMonthly:
		return fmt.Sprintf("每月 %d 号 %s 提醒", t.MonthDay, t.StartTime)
	}
	return ""
}

func parseTimeOfDay(t string) (hour, minute int, err error) {
	var tm time.Time
	tm, err = time.Parse("15:04", t)
	if err != nil {
		return 0, 0, err
	}
	return tm.Hour(), tm.Minute(), nil
}

// calcNextTrigger 计算定时提醒的下次触发时间。
// 修复：weekly 空 weekdays 报错；最大扫描 366 天兜底。
func calcNextTrigger(t models.Task) (*time.Time, error) {
	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	sh, sm, err := parseTimeOfDay(t.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time: %w", err)
	}
	eh, em, err := parseTimeOfDay(t.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time: %w", err)
	}

	var interval time.Duration
	if t.IntervalUnit == models.IntervalUnitHour {
		interval = time.Duration(t.IntervalValue) * time.Hour
	} else {
		interval = time.Duration(t.IntervalValue) * time.Minute
	}

	// weekly 模式下检查 weekdays 有效性
	if t.RepeatMode == models.RepeatModeWeekly {
		valid := false
		for _, s := range strings.Split(t.Weekdays, ",") {
			var d int
			fmt.Sscanf(s, "%d", &d)
			if d >= 1 && d <= 7 {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("weekdays must contain at least one valid day (1-7)")
		}
	}

	// 优先尝试今天的候选时间。
	candidate := today.Add(time.Duration(sh)*time.Hour + time.Duration(sm)*time.Minute)
	end := today.Add(time.Duration(eh)*time.Hour + time.Duration(em)*time.Minute)

	for !candidate.After(now) {
		candidate = candidate.Add(interval)
	}

	maxDays := 366
	daysChecked := 0

	for daysChecked < maxDays {
		if !candidate.After(end) {
			if isValidDay(t, candidate) {
				return &candidate, nil
			}
		}
		// 移到下一天的起始时间。
		candidate = candidate.Add(24 * time.Hour)
		candidate = candidate.Truncate(24 * time.Hour).Add(time.Duration(sh)*time.Hour + time.Duration(sm)*time.Minute)
		end = candidate.Truncate(24 * time.Hour).Add(time.Duration(eh)*time.Hour + time.Duration(em)*time.Minute)
		if isValidDay(t, candidate) {
			for !candidate.After(now) && !candidate.After(end) {
				if isValidDay(t, candidate) {
					return &candidate, nil
				}
				candidate = candidate.Add(interval)
			}
			if !candidate.After(end) {
				return &candidate, nil
			}
		}
		daysChecked++
	}

	return nil, fmt.Errorf("unable to find next trigger within %d days", maxDays)
}

func isValidDay(t models.Task, day time.Time) bool {
	switch t.RepeatMode {
	case models.RepeatModeDaily:
		return true
	case models.RepeatModeWorkday:
		wd := day.Weekday()
		return wd >= time.Monday && wd <= time.Friday
	case models.RepeatModeWeekly:
		wd := int(day.Weekday())
		if wd == 0 {
			wd = 7
		}
		for _, s := range strings.Split(t.Weekdays, ",") {
			var d int
			fmt.Sscanf(s, "%d", &d)
			if d == wd {
				return true
			}
		}
		return false
	case models.RepeatModeMonthly:
		return day.Day() == t.MonthDay
	}
	return true
}

// ListPendingReminders 返回准备触发的活跃提醒（仅 reminder）。
func (s *TaskService) ListPendingReminders(now time.Time) ([]models.Task, error) {
	fmt.Printf("[调试] TaskService.ListPendingReminders 查询待触发提醒: now=%v\n", now)
	list, err := s.store.ListPendingReminders(now)
	if err != nil {
		fmt.Printf("[调试] TaskService.ListPendingReminders 失败: %v\n", err)
		return nil, fmt.Errorf("list pending reminders: %w", err)
	}
	for i := range list {
		enrichTask(&list[i])
	}
	fmt.Printf("[调试] TaskService.ListPendingReminders 成功，共 %d 条\n", len(list))
	return list, nil
}

// ListPendingTodos 返回准备触发的待办任务（仅 todo）。
func (s *TaskService) ListPendingTodos(now time.Time) ([]models.Task, error) {
	fmt.Printf("[调试] TaskService.ListPendingTodos 查询待触发待办: now=%v\n", now)
	list, err := s.store.ListPendingTodos(now)
	if err != nil {
		fmt.Printf("[调试] TaskService.ListPendingTodos 失败: %v\n", err)
		return nil, fmt.Errorf("list pending todos: %w", err)
	}
	for i := range list {
		enrichTask(&list[i])
	}
	fmt.Printf("[调试] TaskService.ListPendingTodos 成功，共 %d 条\n", len(list))
	return list, nil
}

// ClearExpiredPauses 清理已过期的暂停。
func (s *TaskService) ClearExpiredPauses(now time.Time) error {
	return s.store.ClearExpiredPauses(now)
}

// ClearExpiredSkips 清理已过期的跳过区间。
func (s *TaskService) ClearExpiredSkips(now time.Time) error {
	return s.store.ClearExpiredSkips(now)
}

// RecalcNextTrigger 在提醒触发后重新计算并持久化下次触发时间。
func (s *TaskService) RecalcNextTrigger(id int64) error {
	fmt.Printf("[调试] TaskService.RecalcNextTrigger 重新计算触发时间: id=%d\n", id)
	t, err := s.store.Get(id)
	if err != nil {
		fmt.Printf("[调试] TaskService.RecalcNextTrigger 查询失败: %v\n", err)
		return fmt.Errorf("get task: %w", err)
	}
	if t.Type == models.TaskTypeReminder {
		// 清理过期的 skip 区间
		if t.SkipUntil != nil && !t.SkipUntil.After(time.Now()) {
			t.SkipFrom = nil
			t.SkipUntil = nil
		}
		next, err := calcNextTrigger(*t)
		if err != nil {
			return fmt.Errorf("calc next trigger: %w", err)
		}
		t.NextTriggerAt = next
	} else if t.Type == models.TaskTypeTodo {
		// 待办提醒触发一次后清除
		t.NextTriggerAt = nil
	} else {
		fmt.Printf("[调试] TaskService.RecalcNextTrigger 跳过未知类型: type=%s\n", t.Type)
		return nil
	}
	fmt.Printf("[调试] TaskService.RecalcNextTrigger 新的触发时间: next_trigger_at=%v\n", t.NextTriggerAt)
	_, err = s.store.Update(t)
	if err != nil {
		fmt.Printf("[调试] TaskService.RecalcNextTrigger 更新失败: %v\n", err)
	}
	return err
}
