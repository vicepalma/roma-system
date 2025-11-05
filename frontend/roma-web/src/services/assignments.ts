import api from '@/lib/axios'
import type { AssignmentListRow, AssignmentDay } from '@/types/assignments'

export type MyActiveAssignment = {
  id: string
  program_id: string
  program_version: number
  disciple_id: string
  start_date: string
}

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

export async function listAssignmentDays(assignmentId: string): Promise<AssignmentDay[]> {
  const { data } = await api.get<{ items: AssignmentDay[] }>(`/api/assignments/${assignmentId}/days`)
  return data.items ?? []
}

export async function activateAssignment(assignmentId: string, discipleId: string): Promise<void> {
  await api.post(`/api/coach/assignments/${assignmentId}/activate`, null, {
    params: { disciple_id: discipleId },
  })
}

export async function getMyActiveAssignment(): Promise<MyActiveAssignment | null> {
  const { data } = await api.get('/api/me/assignment/active')
  // Si el usuario no tiene activo, backend podría devolver 204/404/{}.
  return data?.id ? (data as MyActiveAssignment) : null
}