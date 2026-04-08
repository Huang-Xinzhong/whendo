import React from 'react'
import styles from './Card.module.css'

export function Card({ children, className = '', padding = 'md', hover = false }) {
  const classes = [
    styles.card,
    styles[padding],
    hover && styles.hover,
    className,
  ].filter(Boolean).join(' ')

  return <div className={classes}>{children}</div>
}
