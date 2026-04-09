import { useState, useEffect, useCallback } from 'react'
import * as api from '../services/api'

export function useTasks(workspaceId, filter = 'all') {
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const refresh = useCallback(async () => {
    if (!workspaceId) {
      setTasks([])
      setLoading(false)
      return
    }
    try {
      setLoading(true)
      const list = await api.TaskList(workspaceId, filter)
      setTasks(list || [])
      setError(null)
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }, [workspaceId, filter])

  useEffect(() => {
    refresh()
  }, [refresh])

  const create = useCallback(async (req) => {
    const t = await api.TaskCreate(req)
    await refresh()
    return t
  }, [refresh])

  const update = useCallback(async (req) => {
    const t = await api.TaskUpdate({ ...req, id: req.id })
    await refresh()
    return t
  }, [refresh])

  const remove = useCallback(async (id) => {
    await api.TaskDelete(id)
    await refresh()
  }, [refresh])

  const toggleComplete = useCallback(async (id) => {
    await api.TaskToggleCompleted(id)
    await refresh()
  }, [refresh])

  const togglePause = useCallback(async (id) => {
    await api.TaskTogglePause(id)
    await refresh()
  }, [refresh])

  return { tasks, loading, error, refresh, create, update, remove, toggleComplete, togglePause }
}
