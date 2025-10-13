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

// Listar sets de una sesión (probamos endpoint dedicado y fallback)
export async function listSets(sessionId: string, prescriptionId?: string): Promise<SetLog[]> {
  const { data } = await api.get(`/api/sessions/${sessionId}/sets`, {
    params: { prescription_id: prescriptionId }
  })
  if (Array.isArray(data?.items)) return data.items as SetLog[]
  if (Array.isArray(data)) return data as SetLog[]
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
