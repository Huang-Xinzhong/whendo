import React from 'react'
import styles from './Input.module.css'

export function Select({
  value,
  defaultValue,
  onChange,
  disabled = false,
  className = '',
  id,
  name,
  children,
  required = false,
}) {
  return (
    <select
      id={id}
      name={name}
      value={value}
      defaultValue={defaultValue}
      onChange={onChange}
      disabled={disabled}
      required={required}
      className={[styles.input, className].filter(Boolean).join(' ')}
    >
      {children}
    </select>
  )
}
