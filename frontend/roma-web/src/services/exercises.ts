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
