import React from 'react'
import styles from './TaskCard.module.css'
import { IconButton } from '../common/IconButton'
import { Icon } from '../common/Icons'
import { Button } from '../common/Button'

export function TaskCard({
  task,
  onEdit,
  onDelete,
  onToggleComplete,
  onTogglePause,
}) {
  const isDone = task.is_completed
  const isReminder = task.type === 'reminder'

  const handleComplete = () => {
    if (onToggleComplete) onToggleComplete(task.id)
  }

  const handlePause = () => {
    if (onTogglePause) onTogglePause(task.id)
  }

  const renderRightAction = () => {
    if (isDone) {
      return (
        <div className={styles.statusBadgeCompleted}>
          已完成
        </div>
      )
    }
    if (isReminder) {
      const isPausedToday = task.pausedToday
      return (
        <Button
          variant={isPausedToday ? 'warning' : 'ghost'}
          size="sm"
          onClick={handlePause}
          className={isPausedToday ? '' : styles.pauseBtn}
        >
          {isPausedToday ? '已暂停' : '暂停'}
        </Button>
      )
    }
    return (
      <Button variant="success" size="sm" onClick={handleComplete}>
        完成
      </Button>
    )
  }

  const renderBadges = () => {
    const badges = []
    const deadlineText = formatDeadline(task)
    const isUrgent = task.due_at && new Date(task.due_at) < new Date(Date.now() + 24 * 60 * 60 * 1000) && !isDone

    badges.push(
      <span key="deadline" className={[styles.badge, isUrgent ? styles.badgeAmber : styles.badgeGray].join(' ')}>
        <Icon name="clock" size={14} />
        截止：{deadlineText}
      </span>
    )

    if (isReminder && task.remindText) {
      const isWeekly = task.remindText.includes('每周')
      badges.push(
        <span key="remind" className={[styles.badge, isWeekly ? styles.badgeEmerald : styles.badgeIndigo].join(' ')}>
          {isWeekly ? <Icon name="refresh" size={14} /> : <Icon name="bell" size={14} />}
          {task.remindText}
        </span>
      )
    } else if (task.remind_at && !isReminder) {
      badges.push(
        <span key="remind" className={[styles.badge, styles.badgeIndigo].join(' ')}>
          <Icon name="bell" size={14} />
          {formatTime(task.remind_at)} 提醒
        </span>
      )
    }

    return badges
  }

  return (
    <div className={[styles.card, isDone && styles.done].filter(Boolean).join(' ')}>
      <div className={styles.content}>
        {isDone ? (
          <>
            <h3 className={[styles.title, styles.titleDone].join(' ')}>{task.title}</h3>
            {task.description && <p className={[styles.desc, styles.descDone].join(' ')}>{task.description}</p>}
            <div className={styles.badges}>
              <span className={[styles.badge, styles.badgeGray].join(' ')}>
                <Icon name="clock" size={14} />
                截止：{formatDeadline(task)}
              </span>
            </div>
          </>
        ) : (
          <>
            <h3 className={styles.title}>{task.title}</h3>
            {task.description && <p className={styles.desc}>{task.description}</p>}
            <div className={styles.badges}>{renderBadges()}</div>
          </>
        )}
      </div>
      <div className={styles.actions}>
        <div className={styles.topActions}>
          <IconButton title="编辑" onClick={() => onEdit && onEdit(task)}>
            <Icon name="edit" size={16} />
          </IconButton>
          <IconButton title="删除" variant="danger" onClick={() => onDelete && onDelete(task.id)}>
            <Icon name="trash" size={16} />
          </IconButton>
        </div>
        <div className={styles.bottomAction}>
          {renderRightAction()}
        </div>
      </div>
    </div>
  )
}

function formatDeadline(task) {
  if (task.is_completed) return '已完成'
  if (task.type === 'reminder') {
    if (!task.start_time && !task.end_time) return '全天'
    if (task.start_time && task.end_time) return `${task.start_time} ~ ${task.end_time}`
    return task.start_time || task.end_time || '全天'
  }
  if (!task.due_at) return '无截止'
  const d = new Date(task.due_at)
  const now = new Date()
  const isToday = d.toDateString() === now.toDateString()
  if (isToday) return `今天 ${formatTime(d)}`
  const tomorrow = new Date(now)
  tomorrow.setDate(tomorrow.getDate() + 1)
  if (d.toDateString() === tomorrow.toDateString()) return `明天 ${formatTime(d)}`
  return `${d.getMonth() + 1}月${d.getDate()}日 ${formatTime(d)}`
}

function formatTime(dateOrString) {
  if (!dateOrString) return ''
  const d = typeof dateOrString === 'string' ? new Date(dateOrString) : dateOrString
  const hours = d.getHours().toString().padStart(2, '0')
  const minutes = d.getMinutes().toString().padStart(2, '0')
  return `${hours}:${minutes}`
}
