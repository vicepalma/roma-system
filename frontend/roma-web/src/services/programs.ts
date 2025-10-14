import api from '@/lib/axios'
import type {
  Program, CreateProgramInput, CreateProgramResponse,
  CreateWeekInput, CreateDayInput, CreatePrescriptionInput
} from '@/types/programs'

export async function listPrograms(): Promise<Program[]> {
  const { data } = await api.get('/api/programs')
  return data?.items ?? data ?? []
}
export async function createProgram(input: CreateProgramInput) {
  const { data } = await api.post('/api/programs', input)
  return data as CreateProgramResponse
}
export async function addWeek(input: CreateWeekInput) {
  const { data } = await api.post(`/api/programs/${input.program_id}/weeks`, { index: input.index })
  return data
}
export async function addDay(input: CreateDayInput) {
  const { data } = await api.post(`/api/weeks/${input.week_id}/days`, {
    day_index: input.day_index,
    notes: input.notes ?? null,
  })
  return data
}
export async function addPrescription(input: CreatePrescriptionInput) {
  const { data } = await api.post(`/api/days/${input.day_id}/prescriptions`, {
    exercise_id: input.exercise_id,
    series: input.series,
    reps: input.reps,
    rest_sec: input.rest_sec ?? null,
    to_failure: !!input.to_failure,
    position: input.position ?? 1,
  })
  return data
}

// Listar detalle del programa con semanas y días
export async function getProgramDetail(id: string) {
  const { data } = await api.get(`/api/programs/${id}`)
  return data
}

