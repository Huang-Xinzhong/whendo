import React, { useEffect, useState } from 'react'
import styles from './TimerReminderModal.module.css'
import { Modal } from '../common/Modal'
import { Button } from '../common/Button'
import { Icon } from '../common/Icons'

export function TimerReminderModal({ isOpen, onClose, task }) {
  const [seconds, setSeconds] = useState(60)

  useEffect(() => {
    if (!isOpen) {
      setSeconds(60)
      return
    }
    setSeconds(60)
    const timer = setInterval(() => {
      setSeconds((s) => {
        if (s <= 1) {
          clearInterval(timer)
          onClose()
          return 0
        }
        return s - 1
      })
    }, 1000)
    return () => clearInterval(timer)
  }, [isOpen, onClose])

  if (!isOpen || !task) return null

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="sm">
      <div className={styles.body}>
        <div className={styles.iconWrap}>
          <Icon name="clock" size={28} className={styles.icon} />
        </div>
        <h2 className={styles.title}>定时提醒</h2>
        <p className={styles.taskTitle}>{task.title}</p>
        {task.description && <p className={styles.taskDesc}>{task.description}</p>}
        <div className={styles.countdown}>
          <div className={styles.countdownNumber}>{seconds}</div>
          <div className={styles.countdownLabel}>秒后自动关闭</div>
        </div>
        <div className={styles.actions}>
          <Button variant="warning" size="md" onClick={onClose}>知道了</Button>
        </div>
      </div>
    </Modal>
  )
}
