// src/store/session.ts
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { immer } from 'zustand/middleware/immer'

type SessionState = {
  currentSessionId: string | null
  currentDiscipleId: string | null
  currentDiscipleName: string | null
  setCurrentSessionId: (id: string | null) => void
  setCurrentDisciple: (id: string | null, name: string | null) => void
  clear: () => void
}

export const useSessionStore = create<SessionState>()(
  persist(
    immer((set) => ({
      currentSessionId: null,
      currentDiscipleId: null,
      currentDiscipleName: null,
      setCurrentSessionId: (id) =>
        set((s) => {
          s.currentSessionId = id
        }),
      setCurrentDisciple: (id, name) =>
        set((s) => {
          s.currentDiscipleId = id
          s.currentDiscipleName = name
        }),
      clear: () =>
        set((s) => {
          s.currentSessionId = null
          s.currentDiscipleId = null
          s.currentDiscipleName = null
        }),
    })),
    { name: 'roma/session' }
  )
)

// Hook utilitario para saber cuándo terminó de hidratar el persist
import { useEffect, useState } from 'react'

export function useSessionHydrated() {
  const [hydrated, setHydrated] = useState(
    // si ya está hidratado al montar, lo marcamos
    (useSessionStore as any).persist?.hasHydrated?.() ?? false
  )
  useEffect(() => {
    const api = (useSessionStore as any).persist
    if (!api) return
    // marca true cuando termine de hidratar
    const unsub = api.onFinishHydration?.(() => setHydrated(true))
    // por si hidrató justo antes
    if (api.hasHydrated?.()) setHydrated(true)
    return () => { unsub?.() }
  }, [])
  return hydrated
}
