package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"whendo/internal/services"
)

// Scheduler 轮询数据库并发送提醒事件。
type Scheduler struct {
	db      *sql.DB
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.Mutex
	taskSvc *services.TaskService
}

// New 创建一个新的 Scheduler 实例。
func New(db *sql.DB) *Scheduler {
	return &Scheduler{
		db:      db,
		taskSvc: services.NewTaskService(db),
	}
}

// Start 在后台 goroutine 中启动调度循环。
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("[调试] Scheduler.Start 启动调度器")
	if s.cancel != nil {
		fmt.Println("[调试] Scheduler.Start 停止旧调度器")
		s.cancel()
	}

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	// 启动两条独立调度链。
	s.wg.Add(2)
	go s.loopReminder(ctx)
	go s.loopTodo(ctx)
}

// Stop 优雅地关闭调度器。
func (s *Scheduler) Stop() {
	fmt.Println("[调试] Scheduler.Stop 停止调度器")
	s.mu.Lock()
	cancel := s.cancel
	s.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	s.wg.Wait()
	fmt.Println("[调试] Scheduler.Stop 调度器已停止")
}

// alignToNextMinute 返回到下一整分（00秒）的等待时长。
func alignToNextMinute() time.Duration {
	now := time.Now()
	return time.Until(now.Add(time.Minute).Truncate(time.Minute))
}

func (s *Scheduler) loopReminder(ctx context.Context) {
	defer s.wg.Done()

	select {
	case <-ctx.Done():
		return
	case <-time.After(alignToNextMinute()):
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	s.tickReminder(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tickReminder(ctx)
		}
	}
}

func (s *Scheduler) loopTodo(ctx context.Context) {
	defer s.wg.Done()

	select {
	case <-ctx.Done():
		return
	case <-time.After(alignToNextMinute()):
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	s.tickTodo(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tickTodo(ctx)
		}
	}
}

func (s *Scheduler) tickReminder(ctx context.Context) {
	store := services.NewTaskService(s.db)
	now := time.Now()
	fmt.Printf("[调试] Scheduler.tickReminder 开始扫描，时间=%v\n", now)

	// 阶段 1: 清理过期暂停与跳过区间。
	if err := store.ClearExpiredPauses(now); err != nil {
		fmt.Printf("[调试] Scheduler.tickReminder 清理过期暂停失败: %v\n", err)
	}
	if err := store.ClearExpiredSkips(now); err != nil {
		fmt.Printf("[调试] Scheduler.tickReminder 清理过期跳过区间失败: %v\n", err)
	}

	// 阶段 2: 查询并触发 reminder。
	list, err := store.ListPendingReminders(now)
	if err != nil {
		fmt.Printf("[调试] Scheduler.tickReminder 扫描失败: %v\n", err)
		return
	}
	fmt.Printf("[调试] Scheduler.tickReminder 发现 %d 个待触发提醒\n", len(list))

	for i := range list {
		t := list[i]
		if t.NextTriggerAt != nil && now.Sub(*t.NextTriggerAt) > 5*time.Minute {
			// 跳过延迟补提（休眠恢复 / 跳过区间恢复后）。
			fmt.Printf("[调试] Scheduler.tickReminder 延迟跳过: task_id=%d next_trigger_at=%v\n", t.ID, t.NextTriggerAt)
			if err := store.RecalcNextTrigger(t.ID); err != nil {
				fmt.Printf("[调试] Scheduler.tickReminder 延迟推进失败: task_id=%d error=%v\n", t.ID, err)
			}
			continue
		}
		fmt.Printf("[调试] Scheduler.tickReminder 发送事件: task_id=%d type=%s next_trigger_at=%v\n", t.ID, t.Type, t.NextTriggerAt)
		runtime.EventsEmit(ctx, "reminder:triggered", t)
		if err := store.RecalcNextTrigger(t.ID); err != nil {
			fmt.Printf("[调试] Scheduler.tickReminder 重新计算触发时间失败: task_id=%d error=%v\n", t.ID, err)
		}
	}
}

func (s *Scheduler) tickTodo(ctx context.Context) {
	store := services.NewTaskService(s.db)
	now := time.Now()
	fmt.Printf("[调试] Scheduler.tickTodo 开始扫描，时间=%v\n", now)
	list, err := store.ListPendingTodos(now)
	if err != nil {
		fmt.Printf("[调试] Scheduler.tickTodo 扫描失败: %v\n", err)
		return
	}
	fmt.Printf("[调试] Scheduler.tickTodo 发现 %d 个待触发待办\n", len(list))

	for i := range list {
		t := list[i]
		fmt.Printf("[调试] Scheduler.tickTodo 发送事件: task_id=%d type=%s next_trigger_at=%v\n", t.ID, t.Type, t.NextTriggerAt)
		runtime.EventsEmit(ctx, "reminder:triggered", t)
		if err := store.RecalcNextTrigger(t.ID); err != nil {
			fmt.Printf("[调试] Scheduler.tickTodo 重新计算触发时间失败: task_id=%d error=%v\n", t.ID, err)
		}
	}
}
