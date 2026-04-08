import React from 'react'
import styles from './IconButton.module.css'

export function IconButton({
  children,
  onClick,
  variant = 'default',
  size = 'md',
  title,
  className = '',
}) {
  const classes = [
    styles.button,
    styles[variant],
    styles[size],
    className,
  ].filter(Boolean).join(' ')

  return (
    <button type="button" onClick={onClick} title={title} className={classes}>
      {children}
    </button>
  )
}
