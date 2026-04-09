package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"whendo/internal/scheduler"
)

// App represents the Wails application state and bound methods.
type App struct {
	ctx       context.Context
	db        *sql.DB
	scheduler *scheduler.Scheduler
}

// NewApp creates a new App instance.
func NewApp(db *sql.DB) *App {
	return &App{
		db:        db,
		scheduler: scheduler.New(db),
	}
}

// Startup is called when the app starts. The context is saved
// so we can interact with the Wails runtime later.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	runtime.MenuSetApplicationMenu(ctx, nil)
	a.scheduler.Start(ctx)
}

// Shutdown is called when the app is about to quit.
func (a *App) Shutdown(_ context.Context) {
	a.scheduler.Stop()
}

// DomReady is called when the frontend dom has loaded.
func (a *App) DomReady(_ context.Context) {
	// Frontend is ready; future runtime calls can safely interact with UI.
}

// GetAppVersion returns the current application version.
func (a *App) GetAppVersion() string {
	return "0.1.0"
}

// mapError transforms internal errors into user-friendly messages.
// This is the central place to add error codes later.
func mapError(err error) error {
	if err == nil {
		return nil
	}
	// TODO: classify known errors (not found, validation, etc.)
	return fmt.Errorf("%w", err)
}
