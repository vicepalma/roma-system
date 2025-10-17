import api from '@/lib/axios'
import type { SearchExercisesResponse } from '@/types/exercises'


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

export const listExercises = (p?: any) => api.get('/api/exercises', { params: p }).then(r => r.data)
export const getExercise = (id: string) => api.get(`/api/exercises/${id}`).then(r => r.data)
export const createExercise = (b: any) => api.post('/api/exercises', b).then(r => r.data)
export const updateExercise = (id: string, b: any) => api.put(`/api/exercises/${id}`, b).then(r => r.data)
export const deleteExercise = (id: string) => api.delete(`/api/exercises/${id}`).then(r => r.data)
