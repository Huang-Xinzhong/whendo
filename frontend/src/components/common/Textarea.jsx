import React from 'react'
import styles from './Input.module.css'

export function Textarea({
  value,
  defaultValue,
  placeholder,
  onChange,
  disabled = false,
  className = '',
  id,
  name,
  rows = 3,
  required = false,
}) {
  return (
    <textarea
      id={id}
      name={name}
      value={value}
      defaultValue={defaultValue}
      placeholder={placeholder}
      onChange={onChange}
      disabled={disabled}
      rows={rows}
      required={required}
      className={[styles.input, className].filter(Boolean).join(' ')}
    />
  )
}
