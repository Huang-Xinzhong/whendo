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
  else if (payload.due_at) payload.due_at = new Date(payload.due_at).toISOString()
  if (payload.remind_at === '') payload.remind_at = null
  else if (payload.remind_at) payload.remind_at = new Date(payload.remind_at).toISOString()
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
