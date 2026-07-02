import api from '@/lib/axios'

export type HistorySession = {
  session_id: string
  assignment_id: string
  program_id?: string
  program_title?: string
  day_id: string
  week_index?: number
  day_index?: number
  day_title?: string | null
  performed_at: string
  status: string
  ended_at?: string | null
  sets: number
  exercises_count?: number
  volume: number
}

export async function getHistorySessions(params: {
  from?: string
  to?: string
  status?: '' | 'open' | 'closed'
  programId?: string
  discipleId?: string
  limit?: number
}) {
  const { data } = await api.get<{ items: HistorySession[]; total: number }>('/api/history', {
    params: {
      group: 'session',
      disciple_id: params.discipleId || undefined,
      from: params.from || undefined,
      to: params.to || undefined,
      status: params.status || undefined,
      program_id: params.programId || undefined,
      limit: params.limit ?? 100,
    },
  })
  return {
    items: data.items ?? [],
    total: data.total ?? 0,
  }
}

export async function getHistoryPivot(params: {
  days: number
  mode: 'by_exercise' | 'by_muscle'
  metric?: 'total_volume' | 'sets' | 'reps'
  includeCatalog?: boolean
}) {
  const { days, mode, metric = 'total_volume', includeCatalog = true } = params
  const { data } = await api.get('/api/history/summary/pivot', {
    params: {
      mode,
      days,
      metric,
      include: includeCatalog ? 'catalog' : undefined,
    },
  })
  return data // { columns, rows, catalog, days, mode }
}
