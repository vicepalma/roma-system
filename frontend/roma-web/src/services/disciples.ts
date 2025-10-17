import api from '@/lib/axios'
import type { Overview, Today } from '@/types/disciples'
import { normalizeToday } from '@/lib/normalize'

export async function getDiscipleOverview(id: string): Promise<Overview> {
  const { data } = await api.get(`/api/coach/disciples/${id}/overview`)
  return data
}

export async function getDiscipleToday(id: string): Promise<Today> {
  const { data } = await api.get(`/api/coach/disciples/${id}/today`)
  return normalizeToday(data) // devuelve Today con week_id incluido
}
