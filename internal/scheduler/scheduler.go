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

	if s.cancel != nil {
		s.cancel()
	}

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	s.wg.Add(1)
	go s.loop(ctx)
}

// Stop 优雅地关闭调度器。
func (s *Scheduler) Stop() {
	s.mu.Lock()
	cancel := s.cancel
	s.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	s.wg.Wait()
}

func (s *Scheduler) loop(ctx context.Context) {
	defer s.wg.Done()

	// 等待到下一个整分（00秒）对齐执行。
	now := time.Now()
	sleep := time.Until(now.Add(time.Minute).Truncate(time.Minute))
	select {
	case <-ctx.Done():
		return
	case <-time.After(sleep):
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	s.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	store := services.NewTaskService(s.db)
	now := time.Now()
	fmt.Printf("[DEBUG] Scheduler.tick started at %v\n", now)
	list, err := store.ListPendingReminders(now)
	if err != nil {
		fmt.Printf("[DEBUG] Scheduler.tick error: %v\n", err)
		return
	}
	fmt.Printf("[DEBUG] Scheduler.tick found %d pending reminders\n", len(list))

	for i := range list {
		t := list[i]
		fmt.Printf("[DEBUG] Scheduler.tick emitting event for task id=%d type=%s next_trigger_at=%v\n", t.ID, t.Type, t.NextTriggerAt)
		// 向前端发送事件。
		runtime.EventsEmit(ctx, "reminder:triggered", t)
		// 重新计算下次触发时间。
		if err := store.RecalcNextTrigger(t.ID); err != nil {
			fmt.Printf("[DEBUG] Scheduler.tick recalc error for task %d: %v\n", t.ID, err)
		}
	}
}
