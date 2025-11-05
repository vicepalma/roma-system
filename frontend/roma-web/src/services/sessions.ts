import api from '@/lib/axios'
import { SessionLog, SessionSet } from '@/types/sessions'


// Crear sesión para un assignment_id + day_id
// export async function startSession(args: {
//   assignment_id: string
//   day_id: string
//   performed_at?: string
//   notes?: string
// }): Promise<SessionLog> {
//   const { data } = await api.post('/api/sessions', args)
//   return data as SessionLog
// }

export async function getSession(id: string) {
  const { data } = await api.get(`/api/sessions/${id}`)
  return data
}

export async function getSessionSets(sessionId: string) {
  const { data } = await api.get(`/api/sessions/${sessionId}/sets`)
  return Array.isArray(data) ? data as SessionSet[] : []
}

export const patchSession = (
  id: string,
  body: {
    // cambiar el día base de la sesión
    day_id?: string

    // “lazy start”: marcar cuándo realmente empezó (al registrar el primer set)
    performed_at?: string | null

    // cierre explícito de la sesión
    status?: 'open' | 'closed'
    ended_at?: string | null

    // opcional: notas rápidas de la sesión
    notes?: string | null
  }
) => api.patch(`/api/sessions/${id}`, body).then(r => r.data)

export const addCardio = (sessionId: string, body: { minutes: number; distance_km?: number; notes?: string|null }) =>
  api.post(`/api/sessions/${sessionId}/cardio`, body).then(r => r.data)

export async function endSession(sessionId: string) {
  // usar el mismo PATCH existente para cerrar
  const { data } = await api.patch<SessionLog>(`/api/sessions/${sessionId}`, {
    status: 'closed',
    ended_at: new Date().toISOString(),
  })
  return data
}

export async function startSession(payload: { assignment_id: string; day_id: string }) {
  const { data } = await api.post<SessionLog>('/api/sessions', payload)
  return data // { id, ... }
}

export async function listSets(sessionId: string, prescriptionId?: string) {
  const url = prescriptionId
    ? `/api/sessions/${sessionId}/sets?prescription_id=${encodeURIComponent(prescriptionId)}`
    : `/api/sessions/${sessionId}/sets`
  const { data } = await api.get<{ items: SessionSet[] }>(url)
  return data.items ?? []
}

export async function addSet(sessionId: string, payload: Omit<SessionSet, 'id' | 'session_id' | 'created_at'>) {
  const { data } = await api.post<SessionSet>(`/api/sessions/${sessionId}/sets`, payload)
  return data
}

export async function deleteSet(sessionId: string, setId: string) {
  await api.delete(`/api/sessions/${sessionId}/sets/${setId}`)
  return true
}

export async function getMyActiveSession() {
  try {
    const { data } = await api.get('/api/me/session/active')
    if (!data || !data.id) return null
    return data as {
      id: string
      assignment_id: string
      disciple_id: string
      day_id: string
      performed_at: string
      notes?: string | null
    }
  }catch (err: any) {
    if (err?.response?.status === 404) return null
    throw err
  }
}