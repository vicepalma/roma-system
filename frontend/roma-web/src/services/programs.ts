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
  series: number
  reps: string
  rest_sec?: number | null
  to_failure: boolean
  position: number
  exercise_name?: string
  primary_muscle?: string
  equipment?: string | null
}

function pickStr(o: any, ...keys: string[]): string | undefined {
  for (const k of keys) if (o?.[k] != null && String(o[k]).trim() !== '') return String(o[k])
}

function pickNum(o: any, ...keys: string[]): number | undefined {
  for (const k of keys) {
    const v = o?.[k]
    if (v !== undefined && v !== null && !Number.isNaN(Number(v))) return Number(v)
  }
}

function pickBool(o: any, ...keys: string[]): boolean {
  for (const k of keys) {
    const v = o?.[k]
    if (typeof v === 'boolean') return v
    if (v === 0 || v === 1) return Boolean(v)
  }
  return false
}

function normalizePresc(raw: any): DayPrescription {
  return {
    id: pickStr(raw, 'id', 'ID')!,
    day_id: pickStr(raw, 'day_id', 'DayID')!,
    exercise_id: pickStr(raw, 'exercise_id', 'ExerciseID')!,
    series: pickNum(raw, 'series', 'Series') ?? 1,
    reps: pickStr(raw, 'reps', 'Reps') ?? '',
    rest_sec: pickNum(raw, 'rest_sec', 'RestSec'),
    to_failure: pickBool(raw, 'to_failure', 'ToFailure'),
    position: pickNum(raw, 'position', 'Position') ?? 1,
    exercise_name: pickStr(raw, 'exercise_name', 'ExerciseName'),
    primary_muscle: pickStr(raw, 'primary_muscle', 'PrimaryMuscle'),
    equipment: pickStr(raw, 'equipment', 'Equipment'),
  }
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
export async function listDays(programId: string, weekId: string): Promise<ProgramDay[]> {
  const { data } = await api.get(`/api/programs/${programId}/weeks/${weekId}/days`)
  const items = Array.isArray(data?.items) ? data.items : []
  return items.map((d: any) => ({
    id: d.id ?? d.ID,
    week_id: d.week_id ?? d.WeekID,
    day_index: d.day_index ?? d.DayIndex,
    notes: d.notes ?? d.Notes ?? null,
  }))
}

export async function addDay(
  programId: string,
  weekId: string,
  payload: { day_index: number; notes?: string | null }
) {
  const body = { ...payload, week_id: weekId } // <-- redundante pero seguro
  const { data } = await api.post(`/api/programs/${programId}/weeks/${weekId}/days`, body)
  return {
    id: data.id ?? data.ID,
    week_id: data.week_id ?? data.WeekID,
    day_index: data.day_index ?? data.DayIndex,
    notes: data.notes ?? data.Notes ?? null,
  }
}

export async function updateDay(programId: string, weekId: string, dayId: string, payload: { title?: string|null; notes?: string|null; day_index?: number }) {
  const { data } = await api.put(`/api/programs/${programId}/weeks/${weekId}/days/${dayId}`, payload)
  return data
}
export async function deleteDay(programId: string, weekId: string, dayId: string) {
  await api.delete(`/api/programs/${programId}/weeks/${weekId}/days/${dayId}`)
  return true
}

// ---- Prescriptions ----
export async function listPrescriptions(dayId: string): Promise<DayPrescription[]> {
  const res = await api.get(`/api/programs/days/${dayId}/prescriptions`)
  const items = Array.isArray(res.data?.items) ? res.data.items : res.data
  return (items || []).map(normalizePresc)
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
    tempo?: string | null
    rir?: number | null
    rpe?: number | null
    notes?: string | null
    method_id: null
  }
) {
  return api.post(`/api/programs/days/${dayId}/prescriptions`, payload)
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