export type HistorySessionRow = {
  session_id: string
  assignment_id: string
  day_id: string
  performed_at: string
  status: string
  ended_at?: string | null
  sets: number
  volume: number
}

export type HistoryDayAgg = {
  day_date: string
  sessions: number
  sets: number
  volume: number
}

export type PlanVsDoneRow = {
  day_date: string
  day_id: string
  planned_sets: number
  done_sets: number
}
