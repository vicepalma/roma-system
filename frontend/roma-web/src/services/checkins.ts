import api from '@/lib/axios'

export type Checkin = {
  id: string
  disciple_id: string
  checked_at: string
  created_at: string
  weight_kg?: number | null
  notes?: string | null
}

export async function listCheckins() {
  const { data } = await api.get<{ items?: Checkin[]; total?: number }>('/api/checkins')
  return {
    items: data.items ?? [],
    total: data.total ?? 0,
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
