import api from '@/lib/axios'

export type RoutineTemplate = {
  id: string
  owner_id: string
  title: string
  notes?: string | null
  status: 'draft' | 'published'
  created_at?: string
  published_at?: string | null
}

export async function listTemplates(): Promise<RoutineTemplate[]> {
  const { data } = await api.get<{ items?: RoutineTemplate[] }>('/api/templates')
  return data.items ?? []
}
export async function createTemplate(payload: { program_id: string; title: string; notes?: string | null }) {
  const { data } = await api.post<RoutineTemplate>('/api/templates', payload); return data
}
export async function publishTemplate(id: string) { await api.post(`/api/templates/${id}/publish`) }
export async function importTemplate(id: string) { const { data } = await api.post<{ program_id: string; kind: string }>(`/api/templates/${id}/import`); return data }
