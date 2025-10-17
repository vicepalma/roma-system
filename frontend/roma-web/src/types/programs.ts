export type Program = {
  id: string
  owner_id: string
  title: string
  notes?: string | null
  visibility: string
  version: number
  created_at?: string
  updated_at?: string
}

export type ProgramWeek = {
  id: string
  program_id: string
  week_index: number
}

export type ProgramDay = {
  id: string
  week_id: string
  day_index: number
  notes?: string | null
}

// Lite (repo)
export type ProgramDayLite = {
  id: string
  week_id: string
  day_index: number
  notes?: string | null
}

export type Prescription = {
  id: string
  day_id: string
  exercise_id: string
  series: number
  reps: string
  rest_sec?: number | null
  to_failure: boolean
  tempo?: string | null
  rir?: number | null
  rpe?: number | null
  method_id?: string | null
  notes?: string | null
  position: number
}

// payloads
export type CreateProgram = {
  title: string
  notes?: string | null
  visibility?: string | null
}

export type UpdateProgram = {
  title?: string
  notes?: string | null
  visibility?: string | null
}

export type TodayPrescription = {
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
}