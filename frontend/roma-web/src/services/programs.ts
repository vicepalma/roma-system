import api from '@/lib/axios'

export type Program = {
  id: string
  title: string
  description?: string | null
  version?: number
}

export type ProgramWeek = {
  id: string
  program_id: string
  week_index: number
  title?: string | null
}

export type ProgramDay = {
  id: string
  week_id: string
  day_index: number
  notes?: string | null
}

export type DayPrescription = {
  id: string
  day_id: string
  exercise_id: string
  exercise_name?: string
  series: number
  reps: string
  rest_sec?: number | null
  to_failure?: boolean
  position: number
  primary_muscle?: string
  equipment?: string | null
}

export async function listPrograms(): Promise<Program[]> {
  const { data } = await api.get('/api/programs')
  return data?.items ?? data ?? []
}

// Listar detalle del programa con semanas y días
export async function getProgramDetail(id: string) {
  const { data } = await api.get(`/api/programs/${id}`)
  return data
}

// ---- Programs ----
export async function listMyPrograms() {
  const { data } = await api.get<{ items?: Program[] }>('/api/programs')
  return data.items ?? []
}
export async function createProgram(payload: { title: string; description?: string | null }) {
  const { data } = await api.post<Program>('/api/programs', payload)
  return data
}
export async function getProgram(id: string) {
  const { data } = await api.get<Program>(`/api/programs/${id}`)
  return data
}

// ---- Weeks ----
export async function listWeeks(programId: string): Promise<ProgramWeek[]> {
  const { data } = await api.get(`/api/programs/${programId}/weeks`)
  // asumiendo { items: [...] }
  return (data.items ?? data ?? []).map((w: any) => ({
    id: w.id ?? w.ID,
    week_index: w.week_index ?? w.WeekIndex,
    title: w.title ?? w.Title ?? null,
  }))
}

export async function addWeek(programId: string, payload: { week_index: number; title?: string | null }) {
  const { data } = await api.post(`/api/programs/${programId}/weeks`, payload)
  return data
}
export async function deleteWeek(programId: string, weekId: string): Promise<void> {
  console.log(programId)
  console.log(weekId)
  await api.delete(`/api/programs/${programId}/weeks/${weekId}`)
}

// ---- Days ----
export async function listDays(programId: string, weekId: string) {
  const { data } = await api.get<{ items?: ProgramDay[] }>(`/api/programs/${programId}/weeks/${weekId}/days`)
  return data.items ?? []
}
export async function addDay(programId: string, weekId: string, payload: { day_index: number; notes?: string | null }) {
  const { data } = await api.post(`/api/programs/${programId}/weeks/${weekId}/days`, payload)
  return data
}
export async function updateDay(programId: string, weekId: string, dayId: string, patch: Partial<{ notes: string | null }>) {
  const { data } = await api.put<ProgramDay>(`/api/programs/${programId}/weeks/${weekId}/days/${dayId}`, patch)
  return data
}
export async function deleteDay(programId: string, weekId: string, dayId: string) {
  await api.delete(`/api/programs/${programId}/weeks/${weekId}/days/${dayId}`)
  return true
}

// ---- Prescriptions ----
export async function listPrescriptions(dayId: string) {
  const { data } = await api.get<{ items?: DayPrescription[] }>(`/api/programs/days/${dayId}/prescriptions`)
  return data.items ?? []
}
export async function addPrescription(
  dayId: string,
  payload: {
    exercise_id: string
    series: number
    reps: string
    rest_sec?: number | null
    to_failure?: boolean
    position: number
  }
) {
  const { data } = await api.post<DayPrescription>(`/api/programs/days/${dayId}/prescriptions`, payload)
  return data
}
export async function updatePrescription(id: string, patch: Partial<Omit<DayPrescription, 'id' | 'day_id'>>) {
  const { data } = await api.put<DayPrescription>(`/api/programs/prescriptions/${id}`, patch)
  return data
}
export async function deletePrescription(id: string) {
  await api.delete(`/api/programs/prescriptions/${id}`)
  return true
}
export async function reorderPrescriptions(payload: { day_id: string; order: string[] }) {
  // order = array de prescription_id en el nuevo orden
  await api.patch('/api/programs/prescriptions/reorder', payload)
  return true
}

export async function deleteProgram(id: string): Promise<void> {
  await api.delete(`/api/programs/${id}`)
}