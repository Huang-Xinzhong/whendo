import React from 'react'
import styles from './TodoReminderModal.module.css'
import { Modal } from '../common/Modal'
import { Button } from '../common/Button'
import { Icon } from '../common/Icons'

export function TodoReminderModal({ isOpen, onClose, task, onComplete, onSnooze }) {
  if (!isOpen || !task) return null

  const deadlineText = task.due_at
    ? formatTime(task.due_at)
    : ''

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="sm">
      <div className={styles.body}>
        <div className={styles.iconWrap}>
          <Icon name="bell" size={28} className={styles.icon} />
        </div>
        <div className={styles.content}>
          <h2 className={styles.title}>待办提醒</h2>
          <p className={styles.subtitle}>您有一个待办任务需要处理</p>
          <div className={styles.taskCard}>
            <div className={styles.taskTitle}>{task.title}</div>
            {task.description && <div className={styles.taskDesc}>{task.description}</div>}
            {deadlineText && (
              <div className={styles.taskDeadline}>
                <Icon name="clock" size={14} />
                截止：{deadlineText}
              </div>
            )}
          </div>
          <div className={styles.actions}>
            <Button variant="success" size="md" onClick={onComplete}>完成</Button>
            <Button variant="ghost" size="md" onClick={onSnooze}>稍后提醒 (10分钟)</Button>
            <Button variant="primary" size="md" onClick={onClose}>知道了</Button>
          </div>
        </div>
      </div>
    </Modal>
  )
}

function formatTime(dateOrString) {
  if (!dateOrString) return ''
  const d = typeof dateOrString === 'string' ? new Date(dateOrString) : dateOrString
  const now = new Date()
  const isToday = d.toDateString() === now.toDateString()
  const timeStr = `${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
  return isToday ? `今天 ${timeStr}` : `${d.getMonth() + 1}月${d.getDate()}日 ${timeStr}`
}
