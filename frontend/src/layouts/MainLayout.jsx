import React from 'react'
import styles from './MainLayout.module.css'
import { Sidebar } from '../components/business/Sidebar'

export function MainLayout({
  children,
  workspaces,
  activeWorkspaceId,
  onSelectWorkspace,
  onAddWorkspace,
  onOpenSettings,
  isSettingsActive,
}) {
  return (
    <div className={styles.layout}>
      <Sidebar
        workspaces={workspaces}
        activeWorkspaceId={activeWorkspaceId}
        onSelectWorkspace={onSelectWorkspace}
        onAddWorkspace={onAddWorkspace}
        onOpenSettings={onOpenSettings}
        isSettingsActive={isSettingsActive}
      />
      <main className={styles.main}>{children}</main>
    </div>
  )
}
