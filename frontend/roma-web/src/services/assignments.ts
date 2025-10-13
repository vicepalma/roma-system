import api from '@/lib/axios'
import type { AssignmentListRow } from '@/types/assignments'

// GET /api/coach/assignments
export async function getCoachAssignments(): Promise<AssignmentListRow[]> {
  const { data } = await api.get('/api/coach/assignments')
  if (Array.isArray(data?.items)) return data.items as AssignmentListRow[]
  if (Array.isArray(data)) return data as AssignmentListRow[]
  return []
}

// POST existente (sin cambios)
export async function createAssignment(input: {
  disciple_id: string
  program_id: string
  program_version: number
  start_date: string
  end_date?: string | null
}) {
  const { data } = await api.post('/api/coach/assignments', input)
  return data as { id: string }
}
