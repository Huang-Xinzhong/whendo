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

// TaskService provides task business logic.
type TaskService struct {
	store *database.TaskStore
}

// NewTaskService creates a new TaskService.
func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{store: database.NewTaskStore(db)}
}

// List returns tasks for a workspace.
func (s *TaskService) List(workspaceID int64, filter string) ([]models.Task, error) {
	list, err := s.store.List(workspaceID, filter)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	for i := range list {
		enrichTask(&list[i])
	}
	return list, nil
}

// Get returns a single task.
func (s *TaskService) Get(id int64) (*models.Task, error) {
	t, err := s.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	enrichTask(t)
	return t, nil
}

// Create validates and creates a task.
func (s *TaskService) Create(req requests.TaskCreateReq) (*models.Task, error) {
	if err := validateTaskCreate(req); err != nil {
		return nil, err
	}
	t := models.Task{
		WorkspaceID:   req.WorkspaceID,
		Title:         strings.TrimSpace(req.Title),
		Description:   strings.TrimSpace(req.Description),
		Type:          req.Type,
		DueAt:         req.DueAt,
		RemindAt:      req.RemindAt,
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
	}
	created, err := s.store.Create(&t)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	enrichTask(created)
	return created, nil
}

// Update validates and updates a task.
func (s *TaskService) Update(req requests.TaskUpdateReq) (*models.Task, error) {
	existing, err := s.store.Get(req.ID)
	if err != nil {
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
		existing.DueAt = req.DueAt
	}
	if req.RemindAt != nil {
		existing.RemindAt = req.RemindAt
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

	if err := validateTask(*existing); err != nil {
		return nil, err
	}

	if existing.Type == models.TaskTypeReminder {
		next, err := calcNextTrigger(*existing)
		if err != nil {
			return nil, fmt.Errorf("calc next trigger: %w", err)
		}
		existing.NextTriggerAt = next
	} else {
		existing.NextTriggerAt = nil
	}

	updated, err := s.store.Update(existing)
	if err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}
	enrichTask(updated)
	return updated, nil
}

// Delete removes a task.
func (s *TaskService) Delete(id int64) error {
	if err := s.store.Delete(id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

// ToggleCompleted flips the completion status.
func (s *TaskService) ToggleCompleted(id int64) (*models.Task, error) {
	t, err := s.store.ToggleCompleted(id)
	if err != nil {
		return nil, fmt.Errorf("toggle completed: %w", err)
	}
	enrichTask(t)
	return t, nil
}

// TogglePause flips the pause status for today.
func (s *TaskService) TogglePause(id int64) (*models.Task, error) {
	t, err := s.store.TogglePause(id)
	if err != nil {
		return nil, fmt.Errorf("toggle pause: %w", err)
	}
	enrichTask(t)
	return t, nil
}

// ClearCompleted removes all completed tasks.
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
		DueAt:         req.DueAt,
		RemindAt:      req.RemindAt,
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
	if t.Type != models.TaskTypeReminder {
		return
	}
	t.RemindText = buildRemindText(t)
	now := time.Now()
	todayStr := now.Format("2006-01-02")
	if t.PausedDate != nil && t.PausedDate.Format("2006-01-02") == todayStr {
		t.PausedToday = true
	} else {
		t.PausedToday = false
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

// calcNextTrigger computes the next trigger time for a reminder task.
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

	// Try today first.
	candidate := today.Add(time.Duration(sh)*time.Hour + time.Duration(sm)*time.Minute)
	end := today.Add(time.Duration(eh)*time.Hour + time.Duration(em)*time.Minute)

	for !candidate.After(now) {
		candidate = candidate.Add(interval)
	}

	for {
		if !candidate.After(end) {
			if isValidDay(t, candidate) {
				return &candidate, nil
			}
		}
		// Move to next valid day start.
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
	}
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

// RecalcNextTrigger recalculates and persists the next trigger time after a reminder fires.
// ListPendingReminders returns active reminders ready to fire.
func (s *TaskService) ListPendingReminders(now time.Time) ([]models.Task, error) {
	list, err := s.store.ListPendingReminders(now)
	if err != nil {
		return nil, fmt.Errorf("list pending reminders: %w", err)
	}
	for i := range list {
		enrichTask(&list[i])
	}
	return list, nil
}

func (s *TaskService) RecalcNextTrigger(id int64) error {
	t, err := s.store.Get(id)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}
	if t.Type != models.TaskTypeReminder {
		return nil
	}
	next, err := calcNextTrigger(*t)
	if err != nil {
		return fmt.Errorf("calc next trigger: %w", err)
	}
	t.NextTriggerAt = next
	_, err = s.store.Update(t)
	return err
}
