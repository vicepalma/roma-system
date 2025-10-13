export type Assignment = {
  id: string
  program_id: string
  program_version: number
  disciple_id: string
  assigned_by: string
  start_date: string            // ISO (date)
  end_date?: string | null      // ISO (date) | null
  is_active: boolean
  created_at: string            // ISO datetime
}

// Para la UI: estado derivado
export type AssignmentUI = Assignment & {
  status: 'upcoming' | 'active' | 'finished' | 'inactive'
}

export type AssignmentListRow = {
  id: string
  disciple_id: string
  disciple_name: string
  disciple_email: string
  program_id: string
  program_title: string
  program_version: number
  start_date: string        // vendrá como RFC3339 desde Go; lo manejamos como string
  end_date?: string | null
  is_active: boolean
  created_at: string
}
