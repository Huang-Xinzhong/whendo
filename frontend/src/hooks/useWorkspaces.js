import { useState, useEffect, useCallback } from 'react'
import * as api from '../services/api'

export function useWorkspaces() {
  const [workspaces, setWorkspaces] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const refresh = useCallback(async () => {
    try {
      setLoading(true)
      const list = await api.WorkspaceList()
      setWorkspaces(list || [])
      setError(null)
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
  }, [refresh])

  const create = useCallback(async (name) => {
    const ws = await api.WorkspaceCreate({ name })
    await refresh()
    return ws
  }, [refresh])

  const update = useCallback(async (id, name) => {
    const ws = await api.WorkspaceUpdate({ id, name })
    await refresh()
    return ws
  }, [refresh])

  const remove = useCallback(async (id) => {
    await api.WorkspaceDelete(id)
    await refresh()
  }, [refresh])

  return { workspaces, loading, error, refresh, create, update, remove }
}
