package scheduler

import (
	"context"
	"database/sql"
)

// Scheduler polls the database and emits reminder events.
type Scheduler struct {
	db     *sql.DB
	cancel context.CancelFunc
}

// New creates a new Scheduler instance.
func New(db *sql.DB) *Scheduler {
	return &Scheduler{db: db}
}

// Start begins the scheduling loop in a background goroutine.
func (s *Scheduler) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	// TODO: run scheduling loop here.
	_ = ctx
}

// Stop signals the scheduler to shut down gracefully.
func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}
