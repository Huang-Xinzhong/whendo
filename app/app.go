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

// App 表示 Wails 应用的状态和绑定方法。
type App struct {
	ctx          context.Context
	db           *sql.DB
	scheduler    *scheduler.Scheduler
	workspaceSvc *services.WorkspaceService
	taskSvc      *services.TaskService
	settingsSvc  *services.SettingsService
}

// NewApp 创建一个新的 App 实例。
func NewApp(db *sql.DB) *App {
	return &App{
		db:           db,
		scheduler:    scheduler.New(db),
		workspaceSvc: services.NewWorkspaceService(db),
		taskSvc:      services.NewTaskService(db),
		settingsSvc:  services.NewSettingsService(db),
	}
}

// Startup 在应用启动时调用。
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	runtime.MenuSetApplicationMenu(ctx, NewMenu(ctx))
	a.scheduler.Start(ctx)
}

// Shutdown 在应用即将退出时调用。
func (a *App) Shutdown(_ context.Context) {
	a.scheduler.Stop()
}

// DomReady 在前端 DOM 加载完成时调用。
func (a *App) DomReady(_ context.Context) {
}

// GetAppVersion 返回当前应用版本。
func (a *App) GetAppVersion() string {
	return "0.1.0"
}

// --- 工作区 ---

// WorkspaceList 返回所有工作区。
func (a *App) WorkspaceList() ([]models.Workspace, error) {
	list, err := a.workspaceSvc.List()
	if err != nil {
		return nil, mapError(err)
	}
	return list, nil
}

// WorkspaceCreate 创建新工作区。
func (a *App) WorkspaceCreate(req requests.WorkspaceCreateReq) (*models.Workspace, error) {
	ws, err := a.workspaceSvc.Create(req)
	if err != nil {
		return nil, mapError(err)
	}
	return ws, nil
}

// WorkspaceUpdate 重命名工作区。
func (a *App) WorkspaceUpdate(req requests.WorkspaceUpdateReq) (*models.Workspace, error) {
	ws, err := a.workspaceSvc.Update(req)
	if err != nil {
		return nil, mapError(err)
	}
	return ws, nil
}

// WorkspaceDelete 删除工作区。
func (a *App) WorkspaceDelete(id int64) error {
	if err := a.workspaceSvc.Delete(id); err != nil {
		return mapError(err)
	}
	return nil
}

// --- 任务 ---

// TaskList 返回工作区下的任务。
func (a *App) TaskList(workspaceID int64, filter string) ([]models.Task, error) {
	list, err := a.taskSvc.List(workspaceID, filter)
	if err != nil {
		return nil, mapError(err)
	}
	return list, nil
}

// TaskGet 返回单个任务。
func (a *App) TaskGet(id int64) (*models.Task, error) {
	t, err := a.taskSvc.Get(id)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// TaskCreate 创建新任务。
func (a *App) TaskCreate(req requests.TaskCreateReq) (*models.Task, error) {
	t, err := a.taskSvc.Create(req)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// TaskUpdate 更新任务。
func (a *App) TaskUpdate(req requests.TaskUpdateReq) (*models.Task, error) {
	t, err := a.taskSvc.Update(req)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// TaskDelete 删除任务。
func (a *App) TaskDelete(id int64) error {
	if err := a.taskSvc.Delete(id); err != nil {
		return mapError(err)
	}
	return nil
}

// TaskToggleCompleted 切换完成状态。
func (a *App) TaskToggleCompleted(id int64) (*models.Task, error) {
	t, err := a.taskSvc.ToggleCompleted(id)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// TaskTogglePause 切换今日暂停状态。
func (a *App) TaskTogglePause(id int64) (*models.Task, error) {
	t, err := a.taskSvc.TogglePause(id)
	if err != nil {
		return nil, mapError(err)
	}
	return t, nil
}

// --- 设置 ---

// SettingsGet 返回所有设置。
func (a *App) SettingsGet() (map[string]string, error) {
	m, err := a.settingsSvc.Get()
	if err != nil {
		return nil, mapError(err)
	}
	return m, nil
}

// SettingsUpdate 更新设置。
func (a *App) SettingsUpdate(req requests.SettingsUpdateReq) error {
	if err := a.settingsSvc.Update(req); err != nil {
		return mapError(err)
	}
	return nil
}

// DataExport 将 SQLite 数据库以 base64 形式导出。
func (a *App) DataExport() (string, error) {
	dir, _ := os.UserConfigDir()
	dbPath := filepath.Join(dir, "whendo", "data.db")
	data, err := os.ReadFile(dbPath)
	if err != nil {
		return "", mapError(fmt.Errorf("read db: %w", err))
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// DataImport 从 base64 替换 SQLite 数据库。
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

// DataClearCompleted 删除所有已完成的任务。
func (a *App) DataClearCompleted() error {
	if err := a.taskSvc.ClearCompleted(); err != nil {
		return mapError(err)
	}
	return nil
}

// mapError 将内部错误转换为对用户友好的消息。
func mapError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w", err)
}
