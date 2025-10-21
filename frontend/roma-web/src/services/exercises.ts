import api from '@/lib/axios'
import type { SearchExercisesResponse } from '@/types/exercises'
import type { Exercise } from '@/types/exercises'

type ExerciseCreate = {
  name: string
  primary_muscle: string
  equipment?: string | null
  notes?: string | null
  tags?: string[]
}

type ExerciseUpdate = Partial<Omit<Exercise, 'id'>>

export async function searchExercises(params?: {
  q?: string
  tags?: string[]
  match?: 'any' | 'all'
  limit?: number
  offset?: number
}) {
  const { q = '', tags = [], match = 'any', limit = 20, offset = 0 } = params || {}
  const { data } = await api.get('/api/exercises', {
    params: {
      q,
      tags: tags.join(','),
      match,
      limit,
      offset,
    },
  })
  return data as SearchExercisesResponse
}

export async function listExercises(): Promise<Exercise[]> {
  const { data } = await api.get('/api/exercises')
  return Array.isArray(data) ? data : (data?.items ?? [])
}

export const getExercise = (id: string) => api.get(`/api/exercises/${id}`).then(r => r.data)

export async function createExercise(payload: {
  name: string
  primary_muscle: string
  equipment?: string | null
  notes?: string | null
  tags?: string[]
}) {
  const body: Record<string, any> = {
    name: payload.name,
    primary_muscle: payload.primary_muscle,
  }

  if (payload.equipment && payload.equipment.trim() !== '') {
    body.equipment = payload.equipment.trim()
  }
  if (payload.notes && payload.notes.trim() !== '') {
    body.notes = payload.notes.trim()
  }
  if (Array.isArray(payload.tags) && payload.tags.length > 0) {
    body.tags = payload.tags
  }

  const { data } = await api.post('/api/exercises', body)
  return data
}

export async function updateExercise(id: string, patch: ExerciseUpdate): Promise<Exercise> {
  const body = {
    ...(patch.name ? { name: patch.name } : {}),
    ...(patch.primary_muscle ? { primary_muscle: patch.primary_muscle } : {}),
    ...(patch.equipment ? { equipment: patch.equipment } : { equipment: undefined }),
    ...(patch.notes ? { notes: patch.notes } : { notes: undefined }),
    ...(Array.isArray(patch.tags) ? { tags: patch.tags } : {}),
  }
  const { data } = await api.put(`/api/exercises/${id}`, body)
  return data
}

export async function deleteExercise(id: string): Promise<void> {
  await api.delete(`/api/exercises/${id}`)
}