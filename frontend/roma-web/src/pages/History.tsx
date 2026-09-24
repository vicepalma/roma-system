import { useState } from 'react'
import { Link, useLocation, useSearchParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getHistoryPivot, getHistorySessions } from '@/services/history'
import { listMyPrograms } from '@/services/programs'
import { getCoachDisciples } from '@/services/coach'
import OverviewVolumeChart from '@/components/charts/OverviewVolumeChart'
import { QueryState } from '@/components/ui/query-state'
import useAuth from '@/store/auth'
import type { CoachDisciple } from '@/types/coach'

function dateISOOffset(daysAgo: number) {
  const d = new Date()
  d.setDate(d.getDate() - daysAgo)
  return d.toISOString().slice(0, 10)
}

function todayISO() {
  return new Date().toISOString().slice(0, 10)
}

export default function History() {
  const [searchParams] = useSearchParams()
  const location = useLocation() as { state?: { discipleName?: string; discipleEmail?: string } }
  const user = useAuth((s) => s.user)
  const discipleId = searchParams.get('disciple_id') || ''
  const [days, setDays] = useState(14)
  const [mode, setMode] = useState<'by_exercise'|'by_muscle'>('by_exercise')
  const [from, setFrom] = useState(dateISOOffset(14))
  const [to, setTo] = useState(todayISO())
  const [status, setStatus] = useState<'' | 'open' | 'closed'>('')
  const [programId, setProgramId] = useState('')

  const q = useQuery({
    queryKey: ['history','pivot', mode, days],
    queryFn: () => getHistoryPivot({ days, mode, metric: 'total_volume', includeCatalog: true }),
    staleTime: 30_000,
  })
  const sessionsQ = useQuery({
    queryKey: ['history', 'sessions', discipleId, from, to, status, programId],
    queryFn: () => getHistorySessions({ discipleId, from, to, status, programId }),
    staleTime: 15_000,
  })
  const programsQ = useQuery({
    queryKey: ['programs', 'history-filter'],
    queryFn: listMyPrograms,
    staleTime: 30_000,
  })
  const disciplesQ = useQuery({
    queryKey: ['coach', 'disciples'],
    queryFn: getCoachDisciples,
    enabled: user?.role === 'coach' && !!discipleId && !location.state?.discipleName,
    staleTime: 5 * 60 * 1000,
  })
  const selectedDisciple = ((disciplesQ.data ?? []) as CoachDisciple[]).find((d) => String(d.id) === String(discipleId))
  const discipleName = location.state?.discipleName || selectedDisciple?.name || ''
  const discipleEmail = location.state?.discipleEmail || selectedDisciple?.email || ''
  const selectedProgram = (programsQ.data ?? []).find((program) => String(program.id) === String(programId))
  const activeFilters = [
    from ? `Desde ${from}` : 'Sin fecha inicial',
    to ? `Hasta ${to}` : 'Sin fecha final',
    status ? (status === 'closed' ? 'Finalizadas' : 'Abiertas') : 'Todos los estados',
    programId ? `Rutina: ${selectedProgram?.title?.trim() || 'seleccionada'}` : 'Todas las rutinas',
  ]

  return (
    <div className="mx-auto max-w-6xl space-y-4 p-3 sm:p-6">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h2 className="text-xl font-semibold">Historial</h2>
          <div className="mt-1 text-sm text-gray-600 dark:text-neutral-300">
            {discipleId ? 'Vista coach del historial de un discípulo.' : 'Historial propio de entrenamiento.'}
          </div>
        </div>
        {discipleId && (
          <Link
            to={`/disciples/${discipleId}`}
            state={discipleName ? { name: discipleName, email: discipleEmail } : undefined}
            className="text-sm rounded border px-3 py-2 bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800 text-blue-600"
          >
            Volver al discípulo
          </Link>
        )}
      </div>
      {discipleId && (
        <div className="rounded border bg-white p-4 dark:bg-neutral-900 dark:border-neutral-800">
          <div className="flex flex-wrap items-start justify-between gap-3">
            <div>
              <div className="text-xs uppercase tracking-wide text-gray-500 dark:text-neutral-400">Discípulo seleccionado</div>
              <div className="mt-1 text-lg font-semibold">
                {disciplesQ.isLoading ? 'Cargando alumno...' : discipleName || `ID ${discipleId}`}
              </div>
              <div className="mt-1 text-sm text-gray-600 dark:text-neutral-300">
                {discipleEmail || 'Email no disponible'}
              </div>
            </div>
            <div className="text-xs text-gray-500 dark:text-neutral-400">
              ID: {discipleId}
            </div>
          </div>
          {disciplesQ.isError && (
            <div className="mt-3 text-sm text-red-600">
              No se pudo cargar el nombre del discípulo, pero el historial se mantiene filtrado por ID.
            </div>
          )}
        </div>
      )}
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

      <div className="rounded border bg-white p-4 dark:bg-neutral-900 dark:border-neutral-800">
        <div className="mb-3 flex flex-wrap items-center justify-between gap-3">
          <div className="font-semibold">Filtros de sesiones</div>
          <div className="text-sm text-gray-600 dark:text-neutral-300">
            {sessionsQ.data?.total ?? 0} resultado{sessionsQ.data?.total === 1 ? '' : 's'}
          </div>
        </div>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          <label className="text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Desde</span>
            <input
              type="date"
              value={from}
              onChange={(e) => setFrom(e.target.value)}
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
            />
          </label>
          <label className="text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Hasta</span>
            <input
              type="date"
              value={to}
              onChange={(e) => setTo(e.target.value)}
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
            />
          </label>
          <label className="text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Estado</span>
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value as '' | 'open' | 'closed')}
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
            >
              <option value="">Todos</option>
              <option value="open">Abiertas</option>
              <option value="closed">Cerradas</option>
            </select>
          </label>
          <label className="text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Rutina</span>
            <select
              value={programId}
              onChange={(e) => setProgramId(e.target.value)}
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
            >
              <option value="">Todas</option>
              {(programsQ.data ?? []).map((program) => (
                <option key={program.id} value={program.id}>
                  {program.title?.trim() || 'Rutina sin título'}
                </option>
              ))}
            </select>
          </label>
        </div>
        <div className="mt-3 flex flex-wrap gap-2 text-xs text-gray-600 dark:text-neutral-300">
          {activeFilters.map((filter) => (
            <span key={filter} className="rounded border px-2 py-1 dark:border-neutral-800">
              {filter}
            </span>
          ))}
        </div>
      </div>

      {q.isLoading && <QueryState title="Cargando resumen" detail="Calculando volumen y actividad reciente." />}
      {q.isError && <QueryState tone="error" title="No se pudo cargar el resumen" detail="El historial de sesiones sigue disponible abajo." />}

      <div className="rounded border bg-white p-4 dark:bg-neutral-900 dark:border-neutral-800">
        <div className="font-semibold mb-3">Sesiones</div>
        {sessionsQ.isLoading && <QueryState title="Cargando sesiones" detail="Buscando entrenamientos con los filtros actuales." />}
        {sessionsQ.isError && <QueryState tone="error" title="No se pudo cargar el historial de sesiones" detail="Revisa los filtros o intenta nuevamente." />}
        {!sessionsQ.isLoading && !sessionsQ.isError && (sessionsQ.data?.items ?? []).length === 0 && (
          <QueryState title="No hay sesiones para estos filtros" detail="Ajusta fechas, estado o rutina para ampliar la busqueda." />
        )}
        <ul className="space-y-2">
          {(sessionsQ.data?.items ?? []).map((session) => {
            const finished = session.status === 'closed'
            return (
              <li key={session.session_id} className="rounded border px-3 py-3 dark:border-neutral-800">
                <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
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
                    className="min-h-11 sm:min-h-0 text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
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
