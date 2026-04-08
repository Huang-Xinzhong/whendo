import React from 'react'
import styles from './WorkspaceItem.module.css'

export function WorkspaceItem({
  name,
  count = 0,
  color = 'blue',
  isActive = false,
  onClick,
}) {
  const colorMap = {
    blue: styles.colorBlue,
    green: styles.colorGreen,
    purple: styles.colorPurple,
  }

  return (
    <button
      type="button"
      onClick={onClick}
      className={[styles.item, isActive && styles.active, colorMap[color]].filter(Boolean).join(' ')}
    >
      <span className={styles.dot} />
      <span className={styles.name}>{name}</span>
      <span className={styles.count}>{count}</span>
    </button>
  )
}
