import api from '@/lib/axios'

export type SessionLog = {
  id: string
  assignment_id: string
  disciple_id: string
  day_id: string
  performed_at: string
  notes?: string | null
}

export type SetLog = {
  id: string
  session_id: string
  prescription_id: string
  set_index: number
  weight?: number | null
  reps: number
  rpe?: number | null
  to_failure: boolean
}

export type SessionSet = {
  id: string
  session_id: string
  prescription_id: string
  set_index: number
  reps: number
  weight?: number | null
  rpe?: number | null
  to_failure?: boolean
  created_at?: string
}

export type LoggedSet = {
  id: string
  session_id: string
  prescription_id: string
  set_index: number
  reps: number
  weight?: number | null
  rpe?: number | null
  to_failure: boolean
  created_at?: string | null
}

// Crear sesión para un assignment_id + day_id
export async function startSession(args: {
  assignment_id: string
  day_id: string
  performed_at?: string
  notes?: string
}): Promise<SessionLog> {
  const { data } = await api.post('/api/sessions', args)
  return data as SessionLog
}

// GET /api/sessions/:id/sets -> { items: LoggedSet[], ... }
export async function listSets(sessionId: string): Promise<LoggedSet[]> {
  const { data } = await api.get(`/api/sessions/${sessionId}/sets`)
  if (Array.isArray(data)) return data as LoggedSet[]
  if (data && Array.isArray(data.items)) return data.items as LoggedSet[]
  return []
}

// Agregar un set a la sesión
export async function addSet(sessionId: string, body: {
  prescription_id: string
  set_index: number
  weight?: number | null
  reps: number
  rpe?: number | null
  to_failure?: boolean
}) {
  const { data } = await api.post(`/api/sessions/${sessionId}/sets`, body)
  return data
}

export async function getSession(id: string) {
  const { data } = await api.get(`/api/sessions/${id}`)
  return data
}

export async function getSessionSets(sessionId: string) {
  const { data } = await api.get(`/api/sessions/${sessionId}/sets`)
  return Array.isArray(data) ? data as SessionSet[] : []
}