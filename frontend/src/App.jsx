import React, { useState, useEffect, useMemo } from 'react'
import { MainLayout } from './layouts/MainLayout'
import { WorkspacePage } from './pages/WorkspacePage'
import { SettingsPage } from './pages/SettingsPage'
import { Modal } from './components/common/Modal'
import { Button } from './components/common/Button'
import { TaskForm } from './components/business/TaskForm'
import { TodoReminderModal } from './components/business/TodoReminderModal'
import { TimerReminderModal } from './components/business/TimerReminderModal'
import { useWorkspaces } from './hooks/useWorkspaces'
import { useTasks } from './hooks/useTasks'
import { useReminders } from './hooks/useReminders'
import { WorkspaceProvider } from './stores/workspaceStore.jsx'
import * as api from './services/api'

const VIEWS = {
  workspace: 'workspace',
  settings: 'settings',
}

function AppContent() {
  const { workspaces, loading: wsLoading, refresh: refreshWorkspaces, create: createWorkspace } = useWorkspaces()
  const [view, setView] = useState(VIEWS.workspace)
  const [activeWorkspaceId, setActiveWorkspaceId] = useState(null)
  const [sort, setSort] = useState('due_at_desc')
  const [filter, setFilter] = useState('all')
  const [taskModalOpen, setTaskModalOpen] = useState(false)
  const [taskModalMode, setTaskModalMode] = useState('create')
  const [editingTask, setEditingTask] = useState(null)
  const [newWorkspaceName, setNewWorkspaceName] = useState('')
  const [wsModalOpen, setWsModalOpen] = useState(false)

  const { tasks, loading: tasksLoading, refresh: refreshTasks, create: createTask, update: updateTask, remove: deleteTask, toggleComplete, togglePause } = useTasks(activeWorkspaceId, filter)
  const { todoReminder, timerReminder, dismissTodoReminder, dismissTimerReminder } = useReminders()

  useEffect(() => {
    if (workspaces.length > 0 && !activeWorkspaceId) {
      api.SettingsGet().then((s) => {
        const defaultId = Number(s.default_workspace_id)
        const found = workspaces.find((w) => w.id === defaultId)
        setActiveWorkspaceId(found ? found.id : workspaces[0].id)
      })
    }
  }, [workspaces, activeWorkspaceId])

  const activeWorkspace = useMemo(() => workspaces.find((w) => w.id === activeWorkspaceId), [workspaces, activeWorkspaceId])

  const handleSelectWorkspace = (id) => {
    setActiveWorkspaceId(id)
    setView(VIEWS.workspace)
    setFilter('all')
  }

  const handleNewTask = () => {
    setEditingTask(null)
    setTaskModalMode('create')
    setTaskModalOpen(true)
  }

  const handleEditTask = (task) => {
    setEditingTask(task)
    setTaskModalMode('edit')
    setTaskModalOpen(true)
  }

  const handleSaveTask = async (form) => {
    if (taskModalMode === 'create') {
      await createTask({ ...form, workspace_id: activeWorkspaceId })
    } else {
      await updateTask({ ...form, id: editingTask.id })
    }
    setTaskModalOpen(false)
    await refreshTasks()
    await refreshWorkspaces()
  }

  const handleDeleteTask = async (id) => {
    await deleteTask(id)
    setTaskModalOpen(false)
    await refreshWorkspaces()
  }

  const handleAddWorkspace = async () => {
    setNewWorkspaceName('')
    setWsModalOpen(true)
  }

  const handleCreateWorkspace = async () => {
    if (!newWorkspaceName.trim()) return
    const ws = await createWorkspace(newWorkspaceName.trim())
    setWsModalOpen(false)
    if (ws) setActiveWorkspaceId(ws.id)
  }

  const currentTasks = tasks

  return (
    <MainLayout
      workspaces={workspaces}
      activeWorkspaceId={activeWorkspaceId}
      onSelectWorkspace={handleSelectWorkspace}
      onAddWorkspace={handleAddWorkspace}
      onOpenSettings={() => setView(VIEWS.settings)}
      isSettingsActive={view === VIEWS.settings}
    >
      {view === VIEWS.workspace && activeWorkspace ? (
        <WorkspacePage
          workspace={activeWorkspace}
          tasks={currentTasks}
          sort={sort}
          filter={filter}
          onSortChange={setSort}
          onFilterChange={setFilter}
          onNewTask={handleNewTask}
          onEditTask={handleEditTask}
          onDeleteTask={deleteTask}
          onToggleComplete={toggleComplete}
          onTogglePause={togglePause}
        />
      ) : (
        <SettingsPage workspaces={workspaces} />
      )}

      <Modal
        isOpen={taskModalOpen}
        onClose={() => setTaskModalOpen(false)}
        title={taskModalMode === 'create' ? '新建任务' : '编辑任务'}
        footer={
          <>
            <Button variant="ghost" onClick={() => setTaskModalOpen(false)}>取消</Button>
            <Button variant="primary" onClick={() => {
              const form = document.getElementById('task-form')
              if (form) form.requestSubmit()
            }}>
              {taskModalMode === 'create' ? '保存任务' : '保存更改'}
            </Button>
          </>
        }
      >
        <TaskForm
          mode={taskModalMode}
          task={editingTask}
          workspaces={workspaces}
          onSubmit={handleSaveTask}
          onDelete={taskModalMode === 'edit' ? () => handleDeleteTask(editingTask.id) : undefined}
        />
      </Modal>

      <Modal
        isOpen={wsModalOpen}
        onClose={() => setWsModalOpen(false)}
        title="新建工作区"
        footer={
          <>
            <Button variant="ghost" onClick={() => setWsModalOpen(false)}>取消</Button>
            <Button variant="primary" onClick={handleCreateWorkspace}>创建</Button>
          </>
        }
      >
        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
          <label style={{ fontSize: 'var(--text-sm)', fontWeight: 500, color: 'var(--color-gray-700)' }}>工作区名称</label>
          <input
            type="text"
            value={newWorkspaceName}
            onChange={(e) => setNewWorkspaceName(e.target.value)}
            placeholder="请输入工作区名称"
            autoFocus
            style={{
              width: '100%',
              padding: '0.625rem 1rem',
              fontSize: 'var(--text-sm)',
              border: '1px solid var(--color-gray-300)',
              borderRadius: 'var(--radius-lg)',
              outline: 'none',
            }}
            onKeyDown={(e) => {
              if (e.key === 'Enter') handleCreateWorkspace()
            }}
          />
        </div>
      </Modal>

      <TodoReminderModal
        isOpen={!!todoReminder}
        onClose={dismissTodoReminder}
        task={todoReminder}
        onComplete={() => {
          dismissTodoReminder()
        }}
        onSnooze={() => {
          dismissTodoReminder()
        }}
      />

      <TimerReminderModal
        isOpen={!!timerReminder}
        onClose={dismissTimerReminder}
        task={timerReminder}
      />
    </MainLayout>
  )
}

function App() {
  return (
    <WorkspaceProvider>
      <div style={{ width: '100%', height: '100%', display: 'flex', backgroundColor: 'var(--color-gray-100)' }}>
        <div style={{ width: '100%', height: '100%' }}>
          <AppContent />
        </div>
      </div>
    </WorkspaceProvider>
  )
}

export default App
