import { createContext, useContext, useState, useCallback } from 'react'

const WorkspaceContext = createContext(null)

export function WorkspaceProvider({ children }) {
  const [activeWorkspaceId, setActiveWorkspaceId] = useState(null)

  const selectWorkspace = useCallback((id) => {
    setActiveWorkspaceId(id)
  }, [])

  return (
    <WorkspaceContext.Provider value={{ activeWorkspaceId, selectWorkspace }}>
      {children}
    </WorkspaceContext.Provider>
  )
}

export function useWorkspaceStore() {
  const ctx = useContext(WorkspaceContext)
  if (!ctx) throw new Error('useWorkspaceStore must be used within WorkspaceProvider')
  return ctx
}
