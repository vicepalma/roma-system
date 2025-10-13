export type NullableString = { String?: string; Valid?: boolean } | string | null | undefined
export type NullableInt = { Int32?: number; Valid?: boolean } | number | null | undefined

export type AdherenceObj = {
  days: number
  days_with_sets: number
  rate: number
}

export type Overview = {
  disciple_id?: string
  adherence?: AdherenceObj | number
  me_today?: {
    day?: {
      id: string
      week_id: string
      day_index: number
      notes?: string
    }
    prescriptions?: Array<Prescription>
  }
  pivot?: {
    columns: string[]
    rows: Array<Record<string, number | string>>
    catalog: Array<{ id: string; name: string }>
    mode: string
    days: number
  }
  // compat previo
  last7Days?: Array<{ date: string; completed: number }> | Record<string, number>
  summary?: string | Record<string, unknown>
}

export type Prescription = {
  id: string
  day_id: string
  exercise_id: string
  exercise_name: string
  series: number
  reps: string
  rest_sec?: number | null
  to_failure: boolean
  position: number
  primary_muscle?: string | null
  equipment?: string | null
}


export type Today = {
  assignment_id?: string
  current_session_id?: string | null
  day?: {
    id: string
    week_id: string
    day_index: number
    notes?: string
  }
  prescriptions?: Prescription[]
}