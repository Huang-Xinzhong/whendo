import React, { useState, useMemo } from 'react'
import styles from './WorkspacePage.module.css'
import { Button } from '../components/common/Button'
import { Select } from '../components/common/Select'
import { Icon } from '../components/common/Icons'
import { TaskCard } from '../components/business/TaskCard'

const SORT_OPTIONS = [
  { value: 'due_at_desc', label: '按截止时间排序' },
  { value: 'created_at_desc', label: '按创建时间排序' },
]

const FILTER_TABS = [
  { key: 'all', label: '全部' },
  { key: 'pending', label: '未完成' },
  { key: 'completed', label: '已完成' },
]

export function WorkspacePage({
  workspace,
  tasks,
  sort,
  filter,
  onSortChange,
  onFilterChange,
  onNewTask,
  onEditTask,
  onDeleteTask,
  onToggleComplete,
  onTogglePause,
}) {
  const [editingName, setEditingName] = useState(false)
  const [nameValue, setNameValue] = useState(workspace?.name || '')

  const stats = useMemo(() => {
    const all = tasks.length
    const done = tasks.filter((t) => t.is_completed).length
    const undone = all - done
    return { all, done, undone }
  }, [tasks])

  const sortedTasks = useMemo(() => {
    const list = [...tasks]
    if (sort === 'due_at_desc') {
      list.sort((a, b) => {
        const da = a.due_at ? new Date(a.due_at) : new Date(0)
        const db = b.due_at ? new Date(b.due_at) : new Date(0)
        return db - da
      })
    } else {
      list.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    }
    return list
  }, [tasks, sort])

  const pendingTodayText = `${stats.undone} 个待办任务`

  return (
    <>
      <header className={styles.header}>
        <div>
          {editingName ? (
            <input
              className={styles.titleInput}
              value={nameValue}
              onChange={(e) => setNameValue(e.target.value)}
              onBlur={() => setEditingName(false)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') setEditingName(false)
              }}
              autoFocus
            />
          ) : (
            <h1 className={styles.title} onDoubleClick={() => setEditingName(true)}>
              {workspace?.name || ''}
            </h1>
          )}
          <p className={styles.subtitle}>今天有 {pendingTodayText}</p>
        </div>
        <div className={styles.headerActions}>
          <Select value={sort} onChange={(e) => onSortChange(e.target.value)} className={styles.sortSelect}>
            {SORT_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>{opt.label}</option>
            ))}
          </Select>
          <Button variant="primary" size="md" onClick={onNewTask} className={styles.newTaskBtn} iconLeft={<Icon name="plus" size={16} />}>
            新建任务
          </Button>
        </div>
      </header>

      <div className={styles.tabs}>
        {FILTER_TABS.map((tab) => {
          const isActive = filter === tab.key
          const count = tab.key === 'all' ? stats.all : tab.key === 'completed' ? stats.done : stats.undone
          return (
            <button
              key={tab.key}
              className={[styles.tab, isActive && styles.tabActive].filter(Boolean).join(' ')}
              onClick={() => onFilterChange(tab.key)}
            >
              {tab.label}
              <span className={styles.tabCount}>{count}</span>
            </button>
          )
        })}
      </div>

      <div className={styles.list}>
        {sortedTasks.length === 0 ? (
          <div className={styles.empty}>
            <p className={styles.emptyText}>暂无任务</p>
            <Button variant="primary" size="md" onClick={onNewTask} iconLeft={<Icon name="plus" size={16} />}>
              新建任务
            </Button>
          </div>
        ) : (
          sortedTasks.map((task) => (
            <TaskCard
              key={task.id}
              task={task}
              onEdit={onEditTask}
              onDelete={onDeleteTask}
              onToggleComplete={onToggleComplete}
              onTogglePause={onTogglePause}
            />
          ))
        )}
      </div>
    </>
  )
}
