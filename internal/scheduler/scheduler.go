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

// Scheduler polls the database and emits reminder events.
type Scheduler struct {
	db      *sql.DB
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.Mutex
	taskSvc *services.TaskService
}

// New creates a new Scheduler instance.
func New(db *sql.DB) *Scheduler {
	return &Scheduler{
		db:      db,
		taskSvc: services.NewTaskService(db),
	}
}

// Start begins the scheduling loop in a background goroutine.
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

// Stop signals the scheduler to shut down gracefully.
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

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Run immediately on start.
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
	list, err := store.ListPendingReminders(now)
	if err != nil {
		// Silently log error; no runtime logger available without context on each call.
		fmt.Printf("scheduler tick error: %v\n", err)
		return
	}

	for i := range list {
		t := list[i]
		// Emit event to frontend.
		runtime.EventsEmit(ctx, "reminder:triggered", t)
		// Recalculate next trigger.
		if err := store.RecalcNextTrigger(t.ID); err != nil {
			fmt.Printf("scheduler recalc error: %v\n", err)
		}
	}
}
