import api from '@/lib/axios'

export type Checkin = {
  id: string
  disciple_id: string
  checked_at: string
  created_at: string
  weight_kg?: number | null
  notes?: string | null
}

export type CheckinListParams = {
  from?: string
  to?: string
  limit?: number
  offset?: number
}

export async function listCheckins(params: CheckinListParams = {}) {
  const { data } = await api.get<{ items?: Checkin[]; total?: number; limit?: number; offset?: number }>('/api/checkins', {
    params: {
      from: params.from || undefined,
      to: params.to || undefined,
      limit: params.limit,
      offset: params.offset,
    },
  })
  return {
    items: data.items ?? [],
    total: data.total ?? 0,
    limit: data.limit ?? params.limit ?? 50,
    offset: data.offset ?? params.offset ?? 0,
  }
}

export async function listCoachDiscipleCheckins(discipleId: string, params: number | CheckinListParams = 5) {
  const queryParams = typeof params === 'number' ? { limit: params } : params
  const { data } = await api.get<{ items?: Checkin[]; total?: number; limit?: number; offset?: number }>(`/api/coach/disciples/${discipleId}/checkins`, {
    params: {
      from: queryParams.from || undefined,
      to: queryParams.to || undefined,
      limit: queryParams.limit,
      offset: queryParams.offset,
    },
  })
  return {
    items: data.items ?? [],
    total: data.total ?? 0,
    limit: data.limit ?? queryParams.limit ?? 50,
    offset: data.offset ?? queryParams.offset ?? 0,
  }
}

export async function createCheckin(payload: {
  checked_at?: string
  weight_kg?: number | null
  notes?: string | null
}) {
  const { data } = await api.post<Checkin>('/api/checkins', {
    checked_at: payload.checked_at || undefined,
    weight_kg: payload.weight_kg ?? null,
    notes: payload.notes?.trim() || null,
  })
  return data
}
