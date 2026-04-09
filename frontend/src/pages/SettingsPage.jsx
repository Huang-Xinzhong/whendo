import React, { useEffect, useState } from 'react'
import styles from './SettingsPage.module.css'
import { Card } from '../components/common/Card'
import { Select } from '../components/common/Select'
import { Button } from '../components/common/Button'
import * as api from '../services/api'

const SORT_OPTIONS = [
  { value: 'due_at_desc', label: '按截止时间排序' },
  { value: 'created_at_desc', label: '按创建时间排序' },
]

export function SettingsPage({ workspaces }) {
  const [settings, setSettings] = useState({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.SettingsGet().then((s) => {
      setSettings(s)
      setLoading(false)
    })
  }, [])

  const updateSetting = async (key, value) => {
    const next = { ...settings, [key]: value }
    setSettings(next)
    await api.SettingsUpdate(next)
  }

  const handleExport = async () => {
    const data = await api.DataExport()
    const blob = new Blob([data], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'whendo-backup.txt'
    a.click()
    URL.revokeObjectURL(url)
  }

  const handleImport = async () => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = '.txt'
    input.onchange = async (e) => {
      const file = e.target.files?.[0]
      if (!file) return
      const text = await file.text()
      try {
        await api.DataImport(text)
        window.location.reload()
      } catch (err) {
        alert('导入失败')
      }
    }
    input.click()
  }

  const handleClearCompleted = async () => {
    if (window.confirm('确定要永久删除所有已完成的任务吗？')) {
      await api.DataClearCompleted()
      window.location.reload()
    }
  }

  const defaultWorkspaceId = String(settings.default_workspace_id || workspaces[0]?.id || '')
  const defaultSort = settings.default_sort || 'due_at_desc'

  return (
    <>
      <header className={styles.header}>
        <h1 className={styles.title}>设置</h1>
      </header>
      <div className={styles.content}>
        <div className={styles.section}>
          <h3 className={styles.sectionTitle}>默认偏好</h3>
          <Card padding="md">
            <div className={styles.settingRow}>
              <div>
                <div className={styles.settingName}>默认工作区</div>
                <div className={styles.settingDesc}>启动应用时默认显示的工作区</div>
              </div>
              <Select
                value={defaultWorkspaceId}
                onChange={(e) => updateSetting('default_workspace_id', e.target.value)}
                className={styles.settingSelect}
                disabled={loading}
              >
                {(workspaces || []).map((ws) => (
                  <option key={ws.id} value={ws.id}>{ws.name}</option>
                ))}
              </Select>
            </div>
            <div className={styles.divider} />
            <div className={styles.settingRow}>
              <div>
                <div className={styles.settingName}>任务列表默认排序</div>
                <div className={styles.settingDesc}>新建任务时的默认排序方式</div>
              </div>
              <Select
                value={defaultSort}
                onChange={(e) => updateSetting('default_sort', e.target.value)}
                className={styles.settingSelect}
                disabled={loading}
              >
                {SORT_OPTIONS.map((opt) => (
                  <option key={opt.value} value={opt.value}>{opt.label}</option>
                ))}
              </Select>
            </div>
          </Card>
        </div>

        <div className={styles.section}>
          <h3 className={styles.sectionTitle}>数据管理</h3>
          <Card padding="md">
            <div className={styles.settingRow}>
              <div>
                <div className={styles.settingName}>导出数据</div>
                <div className={styles.settingDesc}>将所有任务和提醒导出为 SQLite 备份文件</div>
              </div>
              <Button variant="secondary" size="md" onClick={handleExport}>导出</Button>
            </div>
            <div className={styles.divider} />
            <div className={styles.settingRow}>
              <div>
                <div className={styles.settingName}>导入数据</div>
                <div className={styles.settingDesc}>从 SQLite 备份文件恢复数据</div>
              </div>
              <Button variant="secondary" size="md" onClick={handleImport}>导入</Button>
            </div>
            <div className={styles.divider} />
            <div className={styles.settingRow}>
              <div>
                <div className={styles.settingName}>清除已完成任务</div>
                <div className={styles.settingDesc}>永久删除所有已标记为完成的任务</div>
              </div>
              <Button variant="danger" size="md" onClick={handleClearCompleted}>清除</Button>
            </div>
          </Card>
        </div>

        <div className={styles.version}>WhenDo v0.1.0</div>
      </div>
    </>
  )
}
