import api from '@/lib/axios'
import type { CoachDisciple, CoachLink } from '@/types/coach'

// GET /api/coach/disciples  -> { items: [...] }
export async function getCoachDisciples(): Promise<CoachDisciple[]> {
  const { data } = await api.get('/api/coach/disciples')
  return Array.isArray(data?.items) ? data.items : []
}

// GET /api/coach/links -> { incoming: [], outgoing: [...] }
export async function getCoachLinks(): Promise<CoachLink[]> {
  const { data } = await api.get('/api/coach/links')
  if (data?.outgoing && Array.isArray(data.outgoing)) return data.outgoing as CoachLink[]
  return []
}

// POST /api/coach/assignments
export async function assignProgram(body: {
  disciple_id: string
  program_id: string
  program_version: number
  start_date: string // YYYY-MM-DD
}): Promise<{ id: string }> {
  const { data } = await api.post('/api/coach/assignments', body)
  return data
}
