import type { NullInt32, NullString } from './backend.common'

// Formas crudas que vienen del repo para "today"
export type MeTodayDay = {
  id: string
  week_id: string
  day_index: number
  notes: NullString
}

export type MeTodayPrescription = {
  id: string
  day_id: string
  exercise_id: string
  series: number
  reps: string
  rest_sec: NullInt32
  to_failure: boolean
  position: number
  exercise_name: string
  primary_muscle: string
  equipment: NullString
}

// Forma “bonita” que usa la UI (normalizada)
export type Today = {
  assignment_id: string
  current_session_id?: string | null
  current_session_sets_count?: number
  current_session_started_at?: string | null
  day?: {
    id: string
    week_id: string
    day_index: number
    notes?: string | null
  }
  prescriptions?: Array<{
    id: string
    day_id: string
    exercise_id: string
    series: number
    reps: string
    rest_sec?: number | null
    to_failure: boolean
    position: number
    exercise_name: string
    primary_muscle: string
    equipment?: string | null
  }>
}

// Si ya tienes Overview, mantenlo; si no, un mínimo:
export type Pivot = {
  days: number
  series: string[]
  data: any[]
  // si el backend más adelante entrega rows/columns, ajustamos
}

export type Overview = {
  adherence?: number | {
    rate?: number
    days?: number
    days_with_sets?: number
  }
  pivot?: Pivot
}
