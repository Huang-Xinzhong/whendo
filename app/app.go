package app

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"whendo/internal/models"
	"whendo/internal/requests"
	"whendo/internal/scheduler"
	"whendo/internal/services"
)

// App represents the Wails application state and bound methods.
type App struct {
	ctx             context.Context
	db              *sql.DB
	scheduler       *scheduler.Scheduler
	workspaceSvc    *services.WorkspaceService
	taskSvc         *services.TaskService
	settingsSvc     *services.SettingsService
}

// NewApp creates a new App instance.
func NewApp(db *sql.DB) *App {
	return &App{
		db:           db,
		scheduler:    scheduler.New(db),
		workspaceSvc: services.NewWorkspaceService(db),
		taskSvc:      services.NewTaskService(db),
		settingsSvc:  services.NewSettingsService(db),
	}
}

// Startup is called when the app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	runtime.MenuSetApplicationMenu(ctx, NewMenu(ctx))
	a.scheduler.Start(ctx)
}

// Shutdown is called when the app is about to quit.
func (a *App) Shutdown(_ context.Context) {
	a.scheduler.Stop()
}

// DomReady is called when the frontend dom has loaded.
func (a *App) DomReady(_ context.Context) {
}

// GetAppVersion returns the current application version.
func (a *App) GetAppVersion() string {
	return "0.1.0"
}

// --- Workspaces ---

// WorkspaceList returns all workspaces.
func (a *App) WorkspaceList() ([]models.Workspace, error) {
	list, err := a.workspaceSvc.List()
	if err != nil {
		return nil, mapError(err)
	}
	return list, nil
}

// WorkspaceCreate creates a new workspace.
func (a *App) WorkspaceCreate(req requests.WorkspaceCreateReq) (*models.Workspace, error) {
	ws, err := a.workspaceSvc.Create(req)
	if err != nil {
		return nil, mapError(err)
	}
	return ws, nil
}

// WorkspaceUpdate renames a workspace.
func (a *App) WorkspaceUpdate(req requests.WorkspaceUpdateReq) (*models.Workspace, error) {
	ws, err := a.workspaceSvc.Update(req)
	if err != nil {
		return nil, mapError(err)
	}
	return ws, nil
}

// WorkspaceDelete removes a workspace.
func (a *App) WorkspaceDelete(id int64) error {
	if err := a.workspaceSvc.Delete(id); err != nil {
		return mapError(err)
	}
	return nil
}

// --- Tasks ---

// TaskList returns tasks for a workspace.
func (a *App) TaskList(workspaceID int64, filter string) ([]models.Task, error) {
	list, err := a.taskSvc.List(workspaceID, filter)
	if err != nil {
		return nil, mapError(err)
	}
	return list, nil
}

// TaskGet returns a single task.
func (a *App) TaskGet(id int64) (*models.Task, error) {
	t, err := a.taskSvc.Get(id)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// TaskCreate creates a new task.
func (a *App) TaskCreate(req requests.TaskCreateReq) (*models.Task, error) {
	t, err := a.taskSvc.Create(req)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// TaskUpdate updates a task.
func (a *App) TaskUpdate(req requests.TaskUpdateReq) (*models.Task, error) {
	t, err := a.taskSvc.Update(req)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// TaskDelete removes a task.
func (a *App) TaskDelete(id int64) error {
	if err := a.taskSvc.Delete(id); err != nil {
		return mapError(err)
	}
	return nil
}

// TaskToggleCompleted flips completion.
func (a *App) TaskToggleCompleted(id int64) (*models.Task, error) {
	t, err := a.taskSvc.ToggleCompleted(id)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// TaskTogglePause flips pause for today.
func (a *App) TaskTogglePause(id int64) (*models.Task, error) {
	t, err := a.taskSvc.TogglePause(id)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// --- Settings ---

// SettingsGet returns all settings.
func (a *App) SettingsGet() (map[string]string, error) {
	m, err := a.settingsSvc.Get()
	if err != nil {
		return nil, mapError(err)
	}
	return m, nil
}

// SettingsUpdate updates settings.
func (a *App) SettingsUpdate(req requests.SettingsUpdateReq) error {
	if err := a.settingsSvc.Update(req); err != nil {
		return mapError(err)
	}
	return nil
}

// DataExport exports the SQLite database as base64.
func (a *App) DataExport() (string, error) {
	dir, _ := os.UserConfigDir()
	dbPath := filepath.Join(dir, "whendo", "data.db")
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return "", mapError(fmt.Errorf("read db: %w", err))
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// DataImport replaces the SQLite database from base64.
func (a *App) DataImport(fileData string) error {
	data, err := base64.StdEncoding.DecodeString(fileData)
	if err != nil {
		return mapError(fmt.Errorf("decode data: %w", err))
	}
	dir, _ := os.UserConfigDir()
	dbPath := filepath.Join(dir, "whendo", "data.db")
	if err := os.WriteFile(dbPath, data, 0644); err != nil {
		return mapError(fmt.Errorf("write db: %w", err))
	}
	return nil
}

// DataClearCompleted removes all completed tasks.
func (a *App) DataClearCompleted() error {
	if err := a.taskSvc.ClearCompleted(); err != nil {
		return mapError(err)
	}
	return nil
}

// mapError transforms internal errors into user-friendly messages.
func mapError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w", err)
}
