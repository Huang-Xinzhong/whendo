import { useState, useEffect } from 'react'

export function useReminders() {
  const [todoReminder, setTodoReminder] = useState(null)
  const [timerReminder, setTimerReminder] = useState(null)

  // Simulated event listeners for demo purposes
  useEffect(() => {
    // In real app, this would use EventsOn("reminder:triggered", ...)
    return () => {}
  }, [])

  const showTodoReminder = (task) => setTodoReminder(task)
  const dismissTodoReminder = () => setTodoReminder(null)

  const showTimerReminder = (task) => setTimerReminder(task)
  const dismissTimerReminder = () => setTimerReminder(null)

  return {
    todoReminder,
    timerReminder,
    showTodoReminder,
    dismissTodoReminder,
    showTimerReminder,
    dismissTimerReminder,
  }
}
