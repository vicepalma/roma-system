// src/types/sessions.ts
export type SessionLog = {
  id: string
  assignment_id: string
  disciple_id: string
  day_id: string
  performed_at: string
  notes?: string | null
}

export type SessionSet = {
  id: string
  session_id: string
  prescription_id: string
  set_index: number
  weight?: number | null
  reps: number
  rpe?: number | null
  to_failure: boolean
  created_at?: string
}


// Compatibilidad con componentes que esperan LoggedSet
export type LoggedSet = SessionSet
