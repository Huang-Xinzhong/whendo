import React, { useState, useEffect } from 'react'
import styles from './TaskForm.module.css'
import { Input } from '../common/Input'
import { Textarea } from '../common/Textarea'
import { Select } from '../common/Select'
import { Checkbox } from '../common/Checkbox'
import { Button } from '../common/Button'

const INTERVAL_OPTIONS = [
  { value: 15, label: '每 15 分钟' },
  { value: 30, label: '每 30 分钟' },
  { value: 60, label: '每小时' },
  { value: 120, label: '每 2 小时' },
]

const REPEAT_OPTIONS = [
  { value: 'daily', label: '每天' },
  { value: 'workday', label: '每个工作日' },
  { value: 'weekly', label: '每周（自定义）' },
  { value: 'monthly', label: '每月（自定义）' },
]

const WEEKDAYS = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']

export function TaskForm({
  task,
  workspaces,
  onSubmit,
  onDelete,
  mode = 'create',
}) {
  const [form, setForm] = useState({
    title: '',
    description: '',
    workspace_id: workspaces[0]?.id || 1,
    type: 'todo',
    due_at: '',
    remind_at: '',
    start_time: '09:00',
    end_time: '18:00',
    interval_value: 30,
    interval_unit: 'minute',
    repeat_mode: 'daily',
    weekdays: [],
    month_day: 15,
  })

  useEffect(() => {
    if (task) {
      setForm({
        title: task.title || '',
        description: task.description || '',
        workspace_id: task.workspace_id || workspaces[0]?.id || 1,
        type: task.type || 'todo',
        due_at: task.due_at ? toDatetimeLocal(task.due_at) : '',
        remind_at: task.remind_at ? toDatetimeLocal(task.remind_at) : '',
        start_time: task.start_time || '09:00',
        end_time: task.end_time || '18:00',
        interval_value: task.interval_value || 30,
        interval_unit: task.interval_unit || 'minute',
        repeat_mode: task.repeat_mode || 'daily',
        weekdays: task.weekdays ? task.weekdays.split(',').map(Number) : [],
        month_day: task.month_day || 15,
      })
    }
  }, [task, workspaces])

  const handleChange = (key, value) => {
    setForm((prev) => ({ ...prev, [key]: value }))
  }

  const toggleWeekday = (idx) => {
    setForm((prev) => {
      const next = prev.weekdays.includes(idx)
        ? prev.weekdays.filter((d) => d !== idx)
        : [...prev.weekdays, idx].sort((a, b) => a - b)
      return { ...prev, weekdays: next }
    })
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    const payload = {
      ...form,
      weekdays: form.type === 'reminder' && form.repeat_mode === 'weekly'
        ? form.weekdays.join(',')
        : '',
    }
    onSubmit(payload)
  }

  const isReminder = form.type === 'reminder'

  return (
    <form id="task-form" onSubmit={handleSubmit} className={styles.form}>
      <div className={styles.field}>
        <label className={styles.label}>任务标题</label>
        <Input
          value={form.title}
          onChange={(e) => handleChange('title', e.target.value)}
          placeholder="请输入任务标题"
          required
        />
      </div>

      <div className={styles.field}>
        <label className={styles.label}>描述</label>
        <Textarea
          value={form.description}
          onChange={(e) => handleChange('description', e.target.value)}
          placeholder="补充任务详情..."
          rows={3}
        />
      </div>

      <div className={styles.row}>
        <div className={styles.field}>
          <label className={styles.label}>所属工作区</label>
          <Select
            value={form.workspace_id}
            onChange={(e) => handleChange('workspace_id', Number(e.target.value))}
          >
            {workspaces.map((ws) => (
              <option key={ws.id} value={ws.id}>{ws.name}</option>
            ))}
          </Select>
        </div>
        <div className={styles.field}>
          <label className={styles.label}>任务类型</label>
          <Select
            value={form.type}
            onChange={(e) => handleChange('type', e.target.value)}
          >
            <option value="todo">待办</option>
            <option value="reminder">定时提醒</option>
          </Select>
        </div>
      </div>

      <div className={styles.divider} />

      {!isReminder ? (
        <div className={styles.row}>
          <div className={styles.field}>
            <label className={styles.label}>截止时间</label>
            <Input
              type="datetime-local"
              value={form.due_at}
              onChange={(e) => handleChange('due_at', e.target.value)}
            />
          </div>
          <div className={styles.field}>
            <label className={styles.label}>提醒时间</label>
            <Input
              type="datetime-local"
              value={form.remind_at}
              onChange={(e) => handleChange('remind_at', e.target.value)}
            />
          </div>
        </div>
      ) : (
        <div className={styles.reminderFields}>
          <div className={styles.row}>
            <div className={styles.field}>
              <label className={styles.subLabel}>开始时间</label>
              <Input
                type="time"
                value={form.start_time}
                onChange={(e) => handleChange('start_time', e.target.value)}
              />
            </div>
            <div className={styles.field}>
              <label className={styles.subLabel}>结束时间</label>
              <Input
                type="time"
                value={form.end_time}
                onChange={(e) => handleChange('end_time', e.target.value)}
              />
            </div>
          </div>
          <div className={styles.field}>
            <label className={styles.subLabel}>提醒周期</label>
            <Select
              value={form.interval_value}
              onChange={(e) => handleChange('interval_value', Number(e.target.value))}
            >
              {INTERVAL_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </Select>
          </div>
          <div className={styles.field}>
            <label className={styles.subLabel}>重复模式</label>
            <Select
              value={form.repeat_mode}
              onChange={(e) => handleChange('repeat_mode', e.target.value)}
            >
              {REPEAT_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </Select>
          </div>

          {form.repeat_mode === 'weekly' && (
            <div className={styles.field}>
              <label className={styles.subLabel}>选择星期几</label>
              <div className={styles.weekdays}>
                {WEEKDAYS.map((day, idx) => (
                  <label key={day} className={styles.weekdayChip}>
                    <input
                      type="checkbox"
                      checked={form.weekdays.includes(idx + 1)}
                      onChange={() => toggleWeekday(idx + 1)}
                      className={styles.weekdayInput}
                    />
                    <span className={styles.weekdayLabel}>{day}</span>
                  </label>
                ))}
              </div>
            </div>
          )}

          {form.repeat_mode === 'monthly' && (
            <div className={styles.field}>
              <label className={styles.subLabel}>每月几号</label>
              <Select
                value={form.month_day}
                onChange={(e) => handleChange('month_day', Number(e.target.value))}
              >
                {Array.from({ length: 31 }, (_, i) => (
                  <option key={i + 1} value={i + 1}>{i + 1} 号</option>
                ))}
              </Select>
            </div>
          )}
        </div>
      )}

      {mode === 'edit' && onDelete && (
        <div className={styles.deleteArea}>
          <Button type="button" variant="danger" size="md" onClick={onDelete}>
            删除任务
          </Button>
        </div>
      )}
    </form>
  )
}

function toDatetimeLocal(value) {
  const d = new Date(value)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}
