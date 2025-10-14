import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getHistoryPivot } from '@/services/history'
import OverviewVolumeChart from '@/components/charts/OverviewVolumeChart'

export default function History() {
  const [days, setDays] = useState(14)
  const [mode, setMode] = useState<'by_exercise'|'by_muscle'>('by_exercise')

  const q = useQuery({
    queryKey: ['history','pivot', mode, days],
    queryFn: () => getHistoryPivot({ days, mode, metric: 'total_volume', includeCatalog: true }),
    staleTime: 30_000,
  })

  return (
    <div className="mx-auto max-w-6xl p-6 space-y-4">
      <h2 className="text-xl font-semibold">Historia</h2>
      <div className="flex flex-wrap items-center gap-3">
        <label className="text-sm">Días</label>
        <select value={days} onChange={(e) => setDays(Number(e.target.value))}
          className="rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800">
          <option value={7}>7</option>
          <option value={14}>14</option>
          <option value={30}>30</option>
        </select>

        <label className="text-sm ml-4">Modo</label>
        <select value={mode} onChange={(e) => setMode(e.target.value as any)}
          className="rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800">
          <option value="by_exercise">Por ejercicio</option>
          <option value="by_muscle">Por músculo</option>
        </select>
      </div>

      {q.isLoading && <div>Cargando…</div>}
      {q.isError && <div className="text-red-600">No se pudo cargar el pivot</div>}

      {q.data && (
        <div className="rounded border p-4 dark:bg-neutral-900 dark:border-neutral-800">
          <OverviewVolumeChart overview={{ pivot: { ...q.data, days } }} />
        </div>
      )}
    </div>
  )
}
