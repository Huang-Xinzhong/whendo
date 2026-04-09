import React from 'react'
import styles from './Sidebar.module.css'
import { WorkspaceItem } from './WorkspaceItem'
import { Icon } from '../common/Icons'
import { Button } from '../common/Button'

export function Sidebar({
  workspaces,
  activeWorkspaceId,
  onSelectWorkspace,
  onAddWorkspace,
  onOpenSettings,
  isSettingsActive = false,
}) {
  return (
    <aside className={styles.sidebar}>
      <div className={styles.logo}>
        <div className={styles.logoIcon}>
          <Icon name="clock" size={20} className={styles.logoSvg} />
        </div>
        <span className={styles.logoText}>WhenDo</span>
      </div>

      <nav className={styles.nav}>
        <div className={styles.navTitle}>工作区</div>
        <div className={styles.list}>
          {workspaces.map((ws) => (
            <WorkspaceItem
              key={ws.id}
              name={ws.name}
              count={ws.taskCount || 0}
              color={ws.color}
              isActive={ws.id === activeWorkspaceId && !isSettingsActive}
              onClick={() => onSelectWorkspace(ws.id)}
            />
          ))}
        </div>
        <Button variant="ghost" size="sm" className={styles.addBtn} onClick={onAddWorkspace} iconLeft={<Icon name="plus" size={16} />}>
          新建工作区
        </Button>
      </nav>

      <div className={styles.footer}>
        <button
          type="button"
          className={[styles.settingsBtn, isSettingsActive && styles.settingsActive].filter(Boolean).join(' ')}
          onClick={onOpenSettings}
        >
          <Icon name="settings" size={20} />
          设置
        </button>
      </div>
    </aside>
  )
}
