export type SessionLog = {
  id: string
  assignment_id: string
  disciple_id: string
  day_id: string
  performed_at: string
  notes?: string | null
}

export type SetLog = {
  id: string
  session_id: string
  prescription_id: string
  set_index: number
  weight?: number | null
  reps: number
  rpe?: number | null
  to_failure: boolean
}

export type LoggedSet = {
  id: string
  session_id: string
  prescription_id: string
  set_index: number
  reps: number
  weight?: number | null
  rpe?: number | null
  to_failure: boolean
  created_at?: string | null
}

