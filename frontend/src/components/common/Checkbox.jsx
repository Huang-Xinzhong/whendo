import React from 'react'
import styles from './Checkbox.module.css'

export function Checkbox({
  id,
  name,
  checked,
  defaultChecked,
  onChange,
  disabled = false,
  label,
  className = '',
}) {
  return (
    <label className={[styles.wrapper, className].filter(Boolean).join(' ')}>
      <input
        id={id}
        name={name}
        type="checkbox"
        checked={checked}
        defaultChecked={defaultChecked}
        onChange={onChange}
        disabled={disabled}
        className={styles.input}
      />
      <span className={styles.box}>
        <svg className={styles.check} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3" strokeLinecap="round" strokeLinejoin="round">
          <polyline points="20 6 9 17 4 12" />
        </svg>
      </span>
      {label && <span className={styles.label}>{label}</span>}
    </label>
  )
}
