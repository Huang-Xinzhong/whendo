import React from 'react'
import styles from './Input.module.css'

export function Input({
  type = 'text',
  value,
  defaultValue,
  placeholder,
  onChange,
  disabled = false,
  className = '',
  id,
  name,
  required = false,
}) {
  return (
    <input
      id={id}
      name={name}
      type={type}
      value={value}
      defaultValue={defaultValue}
      placeholder={placeholder}
      onChange={onChange}
      disabled={disabled}
      required={required}
      className={[styles.input, className].filter(Boolean).join(' ')}
    />
  )
}
