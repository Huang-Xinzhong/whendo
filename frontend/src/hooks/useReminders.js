import { useState, useEffect } from 'react'
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime.js'

export function useReminders() {
  const [todoReminder, setTodoReminder] = useState(null)
  const [timerReminder, setTimerReminder] = useState(null)

  useEffect(() => {
    if (typeof window === 'undefined' || !window.runtime) {
      return
    }
    const handler = (task) => {
      if (!task) return
      if (task.type === 'todo') {
        setTodoReminder(task)
      } else if (task.type === 'reminder') {
        setTimerReminder(task)
      }
    }

    EventsOn('reminder:triggered', handler)
    return () => {
      EventsOff('reminder:triggered')
    }
  }, [])

  const dismissTodoReminder = () => setTodoReminder(null)
  const dismissTimerReminder = () => setTimerReminder(null)

  return {
    todoReminder,
    timerReminder,
    dismissTodoReminder,
    dismissTimerReminder,
  }
}
