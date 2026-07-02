import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getHistoryPivot, getHistorySessions } from '@/services/history'
import OverviewVolumeChart from '@/components/charts/OverviewVolumeChart'

export default function History() {
  const [days, setDays] = useState(14)
  const [mode, setMode] = useState<'by_exercise'|'by_muscle'>('by_exercise')

  const q = useQuery({
    queryKey: ['history','pivot', mode, days],
    queryFn: () => getHistoryPivot({ days, mode, metric: 'total_volume', includeCatalog: true }),
    staleTime: 30_000,
  })
  const sessionsQ = useQuery({
    queryKey: ['history', 'sessions', days],
    queryFn: () => getHistorySessions({ days }),
    staleTime: 15_000,
  })

  return (
    <div className="mx-auto max-w-6xl p-6 space-y-4">
      <h2 className="text-xl font-semibold">Historial</h2>
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

      <div className="rounded border bg-white p-4 dark:bg-neutral-900 dark:border-neutral-800">
        <div className="font-semibold mb-3">Sesiones</div>
        {sessionsQ.isLoading && <div className="text-sm text-gray-500">Cargando sesiones…</div>}
        {sessionsQ.isError && <div className="text-sm text-red-600">No se pudo cargar el historial de sesiones</div>}
        {!sessionsQ.isLoading && !sessionsQ.isError && (sessionsQ.data?.items ?? []).length === 0 && (
          <div className="text-sm text-gray-500">Sin sesiones en el rango seleccionado.</div>
        )}
        <ul className="space-y-2">
          {(sessionsQ.data?.items ?? []).map((session) => {
            const finished = session.status === 'closed'
            return (
              <li key={session.session_id} className="rounded border px-3 py-3 dark:border-neutral-800">
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <div className="text-sm font-medium">
                      {session.program_title?.trim() || 'Rutina'}
                    </div>
                    <div className="text-xs text-gray-600 dark:text-neutral-300">
                      {new Date(session.performed_at).toLocaleString()} · {finished ? 'Finalizada' : 'Abierta'}
                    </div>
                    <div className="mt-1 text-xs text-gray-600 dark:text-neutral-300">
                      {session.week_index && session.day_index
                        ? `Semana ${session.week_index} · Día ${session.day_index}`
                        : 'Día de entrenamiento'}
                      {session.day_title?.trim() ? ` · ${session.day_title}` : ''}
                    </div>
                    <div className="mt-2 flex flex-wrap gap-2 text-xs text-gray-600 dark:text-neutral-300">
                      <span className="rounded border px-2 py-1 dark:border-neutral-800">
                        {session.exercises_count ?? 0} ejercicio{session.exercises_count === 1 ? '' : 's'}
                      </span>
                      <span className="rounded border px-2 py-1 dark:border-neutral-800">
                        {session.sets} set{session.sets === 1 ? '' : 's'}
                      </span>
                      <span className="rounded border px-2 py-1 dark:border-neutral-800">
                        Volumen {session.volume ? session.volume.toLocaleString() : '—'}
                      </span>
                    </div>
                  </div>
                  <Link
                    to={`/sessions/${session.session_id}`}
                    className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                  >
                    Ver resumen
                  </Link>
                </div>
              </li>
            )
          })}
        </ul>
      </div>

      {q.data && (
        <div className="rounded border p-4 dark:bg-neutral-900 dark:border-neutral-800">
          <OverviewVolumeChart overview={{ pivot: { ...q.data, days } }} />
        </div>
      )}
    </div>
  )
}
