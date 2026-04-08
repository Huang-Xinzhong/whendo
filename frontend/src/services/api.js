// Mock API layer matching api_contract.md
// Replace with real Wails bindings once backend is ready

const MOCK_WORKSPACES = [
  { id: 1, name: '工作', color: 'blue', taskCount: 4, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 2, name: '家庭', color: 'green', taskCount: 4, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 3, name: '个人', color: 'purple', taskCount: 6, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
]

const now = new Date()
const today = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 15, 0)
const later = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 18, 0)
const remind = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 14, 30)

const MOCK_TASKS = [
  { id: 1, workspace_id: 1, title: '完成季度报告', description: '整理 Q2 的数据并生成 PPT', type: 'todo', due_at: later.toISOString(), remind_at: null, is_completed: true, start_time: '', end_time: '', interval_value: 0, interval_unit: '', repeat_mode: '', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 2, workspace_id: 1, title: '回复客户邮件', description: '关于下周合作方案的沟通', type: 'todo', due_at: today.toISOString(), remind_at: remind.toISOString(), is_completed: false, start_time: '', end_time: '', interval_value: 0, interval_unit: '', repeat_mode: '', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString(), remindText: '14:30 提醒' },
  { id: 3, workspace_id: 1, title: '团队周会准备', description: '收集各成员进度更新，整理会议议程', type: 'reminder', due_at: null, remind_at: null, is_completed: false, start_time: '09:00', end_time: '18:00', interval_value: 30, interval_unit: 'minute', repeat_mode: 'weekly', weekdays: '1,4', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString(), remindText: '每周一 / 周四 09:00 提醒' },
  { id: 4, workspace_id: 1, title: '喝水提醒', description: '每半个小时喝一次水，保持身体水分', type: 'reminder', due_at: null, remind_at: null, is_completed: false, start_time: '09:00', end_time: '18:00', interval_value: 30, interval_unit: 'minute', repeat_mode: 'daily', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString(), remindText: '每 30 分钟提醒' },

  { id: 5, workspace_id: 2, title: '缴纳物业费', description: '第二季度物业费 580 元', type: 'reminder', due_at: null, remind_at: null, is_completed: false, start_time: '09:00', end_time: '18:00', interval_value: 60, interval_unit: 'minute', repeat_mode: 'daily', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString(), remindText: '19:00 提醒' },
  { id: 6, workspace_id: 2, title: '采购周末食材', description: '牛肉、蔬菜、牛奶、水果', type: 'todo', due_at: new Date(now.getFullYear(), now.getMonth(), now.getDate() + 2, 10, 0).toISOString(), remind_at: null, is_completed: false, start_time: '', end_time: '', interval_value: 0, interval_unit: '', repeat_mode: '', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 7, workspace_id: 2, title: '订机票', description: '五一假期回老家的往返机票', type: 'todo', due_at: new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1, 10, 0).toISOString(), remind_at: null, is_completed: true, start_time: '', end_time: '', interval_value: 0, interval_unit: '', repeat_mode: '', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 8, workspace_id: 2, title: '更换净水器滤芯', description: '厨房台下净水器 3 号滤芯', type: 'reminder', due_at: null, remind_at: null, is_completed: false, start_time: '09:00', end_time: '18:00', interval_value: 60, interval_unit: 'minute', repeat_mode: 'weekly', weekdays: '7', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString(), remindText: '15:00 提醒' },

  { id: 9, workspace_id: 3, title: '晨跑 5 公里', description: '配速控制在 6 分钟以内', type: 'reminder', due_at: null, remind_at: null, is_completed: false, start_time: '06:30', end_time: '07:30', interval_value: 15, interval_unit: 'minute', repeat_mode: 'daily', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString(), remindText: '06:30 提醒' },
  { id: 10, workspace_id: 3, title: '阅读《深度学习》', description: '完成第三章并做笔记', type: 'todo', due_at: new Date(now.getFullYear(), now.getMonth(), now.getDate() + 6, 22, 0).toISOString(), remind_at: null, is_completed: false, start_time: '', end_time: '', interval_value: 0, interval_unit: '', repeat_mode: '', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 11, workspace_id: 3, title: '预约牙医', description: '半年一次的洗牙检查', type: 'reminder', due_at: null, remind_at: null, is_completed: false, start_time: '09:00', end_time: '18:00', interval_value: 60, interval_unit: 'minute', repeat_mode: 'monthly', weekdays: '', month_day: 15, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString(), remindText: '前一天 18:00 提醒' },
  { id: 12, workspace_id: 3, title: '整理相册', description: '备份手机照片到 NAS', type: 'todo', due_at: new Date(now.getFullYear(), now.getMonth(), now.getDate() - 2, 20, 0).toISOString(), remind_at: null, is_completed: true, start_time: '', end_time: '', interval_value: 0, interval_unit: '', repeat_mode: '', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 13, workspace_id: 3, title: '买生日礼物', description: '给妈妈的生日礼物挑选', type: 'reminder', due_at: null, remind_at: null, is_completed: false, start_time: '09:00', end_time: '18:00', interval_value: 60, interval_unit: 'minute', repeat_mode: 'monthly', weekdays: '', month_day: 12, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString(), remindText: '10:00 提醒' },
  { id: 14, workspace_id: 3, title: '学习吉他', description: '练习《晴天》前奏指法', type: 'todo', due_at: new Date(now.getFullYear(), now.getMonth(), now.getDate() - 3, 19, 0).toISOString(), remind_at: null, is_completed: true, start_time: '', end_time: '', interval_value: 0, interval_unit: '', repeat_mode: '', weekdays: '', month_day: 0, next_trigger_at: null, paused_date: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
]

let tasks = [...MOCK_TASKS]
let workspaces = [...MOCK_WORKSPACES]
let settings = {
  default_workspace_id: '1',
  default_sort: 'created_at_desc',
}

function updateWorkspaceCounts() {
  workspaces = workspaces.map((ws) => ({
    ...ws,
    taskCount: tasks.filter((t) => t.workspace_id === ws.id).length,
  }))
}

// Workspaces
export async function WorkspaceList() {
  updateWorkspaceCounts()
  return [...workspaces]
}

export async function WorkspaceCreate(req) {
  const ws = {
    id: Date.now(),
    name: req.name,
    color: 'blue',
    taskCount: 0,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  }
  workspaces.push(ws)
  return ws
}

export async function WorkspaceUpdate(req) {
  const idx = workspaces.findIndex((w) => w.id === req.id)
  if (idx === -1) throw new Error('workspace not found')
  workspaces[idx] = { ...workspaces[idx], name: req.name, updated_at: new Date().toISOString() }
  return workspaces[idx]
}

export async function WorkspaceDelete(id) {
  workspaces = workspaces.filter((w) => w.id !== id)
  tasks = tasks.filter((t) => t.workspace_id !== id)
  return null
}

// Tasks
export async function TaskList(workspaceID, filter = 'all') {
  let list = tasks.filter((t) => t.workspace_id === workspaceID)
  if (filter === 'pending') list = list.filter((t) => !t.is_completed)
  if (filter === 'completed') list = list.filter((t) => t.is_completed)
  // sort
  list = list.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
  return [...list]
}

export async function TaskGet(id) {
  const t = tasks.find((x) => x.id === id)
  if (!t) throw new Error('task not found')
  return { ...t }
}

export async function TaskCreate(req) {
  const t = {
    id: Date.now(),
    workspace_id: req.workspace_id || req.workspaceId,
    title: req.title,
    description: req.description || '',
    type: req.type || 'todo',
    due_at: req.due_at || req.dueAt || null,
    remind_at: req.remind_at || req.remindAt || null,
    is_completed: false,
    start_time: req.start_time || req.startTime || '',
    end_time: req.end_time || req.endTime || '',
    interval_value: req.interval_value || req.intervalValue || 0,
    interval_unit: req.interval_unit || req.intervalUnit || '',
    repeat_mode: req.repeat_mode || req.repeatMode || '',
    weekdays: req.weekdays || '',
    month_day: req.month_day || req.monthDay || 0,
    next_trigger_at: null,
    paused_date: null,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  }
  if (t.type === 'reminder') {
    const intervalLabel = t.interval_value === 60 ? '每小时' : t.interval_value === 120 ? '每 2 小时' : `每 ${t.interval_value} 分钟`
    if (t.repeat_mode === 'daily') t.remindText = `${intervalLabel}提醒`
    else if (t.repeat_mode === 'workday') t.remindText = `每个工作日 ${t.start_time} 提醒`
    else if (t.repeat_mode === 'weekly') {
      const days = t.weekdays.split(',').map((d) => '周' + '一二三四五六日'[d - 1]).join(' / ')
      t.remindText = `每周${days} ${t.start_time} 提醒`
    } else if (t.repeat_mode === 'monthly') t.remindText = `每月 ${t.month_day} 号 ${t.start_time} 提醒`
  }
  tasks.push(t)
  updateWorkspaceCounts()
  return { ...t }
}

export async function TaskUpdate(req) {
  const idx = tasks.findIndex((t) => t.id === req.id)
  if (idx === -1) throw new Error('task not found')
  tasks[idx] = { ...tasks[idx], ...req, updated_at: new Date().toISOString() }
  return { ...tasks[idx] }
}

export async function TaskDelete(id) {
  tasks = tasks.filter((t) => t.id !== id)
  updateWorkspaceCounts()
  return null
}

export async function TaskToggleCompleted(id) {
  const idx = tasks.findIndex((t) => t.id === id)
  if (idx === -1) throw new Error('task not found')
  tasks[idx].is_completed = !tasks[idx].is_completed
  return { ...tasks[idx] }
}

export async function TaskTogglePause(id) {
  const idx = tasks.findIndex((t) => t.id === id)
  if (idx === -1) throw new Error('task not found')
  const todayStr = new Date().toISOString().slice(0, 10)
  const paused = tasks[idx].paused_date?.slice(0, 10) === todayStr
  tasks[idx].paused_date = paused ? null : todayStr
  tasks[idx].pausedToday = !paused
  return { ...tasks[idx] }
}

// Settings
export async function SettingsGet() {
  return { ...settings }
}

export async function SettingsUpdate(req) {
  settings = { ...settings, ...req }
  return null
}

export async function DataExport() {
  return btoa(JSON.stringify({ workspaces, tasks, settings }))
}

export async function DataImport(fileData) {
  const data = JSON.parse(atob(fileData))
  workspaces = data.workspaces || workspaces
  tasks = data.tasks || tasks
  settings = data.settings || settings
  return null
}

export async function DataClearCompleted() {
  tasks = tasks.filter((t) => !t.is_completed)
  updateWorkspaceCounts()
  return null
}
