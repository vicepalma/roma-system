// Para la UI: estado derivado
export type AssignmentUI = Assignment & {
  status: 'upcoming' | 'active' | 'finished' | 'inactive'
}

// src/types/assignments.ts
export type Assignment = {
  id: string
  program_id: string
  program_version: number
  disciple_id: string
  assigned_by: string
  start_date: string
  end_date?: string | null
  is_active: boolean
  created_at: string
}

export type AssignmentMinimal = {
  id: string
  program_id: string
  program_version: number
  disciple_id: string
  assigned_by: string
  start_date: string
  is_active: boolean
  created_at: string
}

export type AssignmentListRow = {
  id: string
  disciple_id: string
  disciple_name: string
  disciple_email: string
  program_id: string
  program_title: string
  program_version: number
  start_date: string
  end_date?: string | null
  is_active: boolean
  created_at: string
}
