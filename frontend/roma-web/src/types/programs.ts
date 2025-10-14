export type Program = {
  id: string
  title: string
  version: number
  created_at?: string
}

export type CreateProgramInput = { title: string }
export type CreateProgramResponse = Program

export type CreateWeekInput = { program_id: string; index: number }
export type CreateDayInput = { week_id: string; day_index: number; notes?: string | null }
export type CreatePrescriptionInput = {
  day_id: string
  exercise_id: string
  series: number
  reps: string
  rest_sec?: number | null
  to_failure?: boolean
  position?: number
}
