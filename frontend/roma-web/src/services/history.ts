import api from '@/lib/axios'

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
