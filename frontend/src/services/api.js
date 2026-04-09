// Real Wails API layer
import {
  WorkspaceList,
  WorkspaceCreate,
  WorkspaceUpdate,
  WorkspaceDelete,
  TaskList,
  TaskGet,
  TaskCreate as GoTaskCreate,
  TaskUpdate as GoTaskUpdate,
  TaskDelete,
  TaskToggleCompleted,
  TaskTogglePause,
  SettingsGet,
  SettingsUpdate,
  DataExport,
  DataImport,
  DataClearCompleted,
} from '../wailsjs/go/app/App.js'

function normalizeTaskPayload(req) {
  const payload = { ...req }
  if (payload.due_at === '') payload.due_at = null
  if (payload.remind_at === '') payload.remind_at = null
  return payload
}

export async function TaskCreate(req) {
  return GoTaskCreate(normalizeTaskPayload(req))
}

export async function TaskUpdate(req) {
  return GoTaskUpdate(normalizeTaskPayload(req))
}

// Workspaces
export { WorkspaceList, WorkspaceCreate, WorkspaceUpdate, WorkspaceDelete }

// Tasks
export { TaskList, TaskGet, TaskDelete, TaskToggleCompleted, TaskTogglePause }

// Settings
export { SettingsGet, SettingsUpdate }

// Data
export { DataExport, DataImport, DataClearCompleted }
