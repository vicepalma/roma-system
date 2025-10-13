// Convierte overview.pivot -> datos aptos para Recharts
// - rows: [{ date: 'YYYY-MM-DD', 'Ejercicio A': 100, 'Ejercicio B': 50, ... }]
// - columns: ['date', 'Ejercicio A', 'Ejercicio B', ...]
export type Pivot = {
  columns: string[]
  rows: Array<Record<string, number | string>>
}

export function pivotToChartData(pivot?: Pivot, limitDays = 14) {
  if (!pivot?.rows?.length || !pivot?.columns?.length) return { data: [], seriesKeys: [] as string[] }
  const rows = pivot.rows.slice(-limitDays)
  const seriesKeys = pivot.columns.filter((c) => c !== 'date')
  const data = rows.map((r) => {
    const date = String(r['date'] ?? '')
    const obj: Record<string, any> = { date }
    for (const k of seriesKeys) obj[k] = Number(r[k] ?? 0) || 0
    return obj
  })
  return { data, seriesKeys }
}
