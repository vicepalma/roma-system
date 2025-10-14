import { useEffect, useMemo, useState } from 'react'
import { useParams, NavLink, useLocation } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useToast } from '@/components/toast/ToastProvider'

import { getDiscipleOverview, getDiscipleToday } from '@/services/disciples'
import { getCoachDisciples } from '@/services/coach'
import { startSession, addSet, listSets } from '@/services/sessions'

import type { Overview, Prescription } from '@/types/disciples'
import type { CoachDisciple } from '@/types/coach'

import OverviewVolumeChart from '@/components/charts/OverviewVolumeChart'
import { KpiCard } from '@/components/kpi/KpiCard'
import { RangePicker } from '@/components/charts/RangePicker'
import { pivotToChartData } from '@/lib/pivot'
import Modal from '@/components/ui/Modal'
import LogSetForm from '@/components/forms/LogSetForm'
import TodaySets from '@/components/sessions/TodaySets'
import { nextSetIndexForPrescription } from '@/lib/sets'
import { useSessionStore } from '@/store/session'

// ---------- Helpers ----------
function formatPercent(value: number, decimals = 1, locale = 'es-CL') {
  const v = value <= 1 ? value * 100 : value
  return new Intl.NumberFormat(locale, { minimumFractionDigits: decimals, maximumFractionDigits: decimals }).format(v)
}
function formatAdherence(ov?: Overview, decimals = 1) {
  const a = ov?.adherence
  if (typeof a === 'number') return `${formatPercent(a, decimals)}%`
  if (a && typeof a === 'object') {
    const days = a.days ?? 0
    const withSets = a.days_with_sets ?? 0
    const pct = formatPercent(Number(a.rate ?? 0), decimals)
    return `${pct}% (${withSets}/${days} días)`
  }
  return '-'
}

// UI: lista de prescripciones
function renderPrescriptions(list: Prescription[] | undefined, onLog: (p: Prescription) => void) {
  const rows = Array.isArray(list) ? list : []
  if (!rows.length) return <div className="text-gray-500">Sin sesiones para hoy</div>
  return (
    <ul className="mt-2 space-y-2">
      {rows.map((p) => (
        <li key={p.id} className="rounded border px-3 py-2">
          <div className="flex items-center justify-between">
            <div className="font-medium">{p.exercise_name}</div>
            <div className="flex items-center gap-2">
              <span className="text-xs text-gray-500">{p.equipment ?? ''}</span>
              <button
                onClick={() => onLog(p)}
                className="text-xs rounded border px-2 py-1 bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
              >
                Registrar set
              </button>
            </div>
          </div>
          <div className="text-sm text-gray-700">
            Series: {p.series ?? '-'} · Reps: {p.reps ?? '-'} · Descanso: {p.rest_sec ?? '-'}s
          </div>
          {p.primary_muscle && (
            <div className="text-xs text-gray-500">Músculo: {p.primary_muscle}</div>
          )}
        </li>
      ))}
    </ul>
  )
}

// ---------- Referencias estables fuera del componente ----------
const NOOP = () => {}

export default function DiscipleDetail() {
  const { id = '' } = useParams()
  const location = useLocation() as { state?: { name?: string; email?: string } }
  const { show } = useToast()
  const qc = useQueryClient()

  // Store setters (no crean nuevas referencias en cada render)
  const setCurrentSessionId = useSessionStore(s => s.setCurrentSessionId)
  const setCurrentDisciple = useSessionStore(s => s.setCurrentDisciple)
  // Si tu store aún no define setSessionForDisciple, usamos NOOP estable
  const setSessionForDisciple = useMemo(
    () => ((useSessionStore.getState() as any).setSessionForDisciple ?? NOOP),
    []
  )

  // Sesión activa local
  const [sessionId, setSessionId] = useState<string | null>(null)

  // Nombre del discípulo (cache o state de navegación)
  const disciplesQ = useQuery({
    queryKey: ['coach', 'disciples'],
    queryFn: getCoachDisciples,
    staleTime: 5 * 60 * 1000,
  })
  const displayName = useMemo(() => {
    if (location.state?.name) return location.state.name
    const list = (disciplesQ.data ?? []) as CoachDisciple[]
    const found = list.find(d => String(d.id) === String(id))
    return found?.name ?? id
  }, [location.state?.name, disciplesQ.data, id])

  // Datos principales
  const todayQ = useQuery({
    queryKey: ['disciple', id, 'today'],
    queryFn: () => getDiscipleToday(id),
    enabled: !!id,
    retry: 1,
    retryDelay: 800,
  })
  const overviewQ = useQuery({
    queryKey: ['disciple', id, 'overview'],
    queryFn: () => getDiscipleOverview(id),
    enabled: !!id,
    retry: 1,
    retryDelay: 800,
  })
  useEffect(() => { if (overviewQ.isError) show({ type: 'error', message: 'Error al cargar overview' }) }, [overviewQ.isError, show])
  useEffect(() => { if (todayQ.isError) show({ type: 'error', message: 'Error al cargar datos de hoy' }) }, [todayQ.isError, show])

  const ov = overviewQ.data
  const today = todayQ.data
  const sidFromToday = today?.current_session_id ?? null

  // Sincroniza local + store con el current_session_id que viene del backend (efecto único)
  useEffect(() => {
    // siempre propaga el discípulo visible
    setCurrentDisciple(id || null, displayName || null)

    if (sidFromToday && sidFromToday !== sessionId) {
      setSessionId(sidFromToday)
      setCurrentSessionId(sidFromToday)
      try { setSessionForDisciple(id!, sidFromToday) } catch {}

      // Prefetch sets de la sesión activa para suavizar
      qc.prefetchQuery({
        queryKey: ['session', sidFromToday, 'sets'],
        queryFn: () => listSets(sidFromToday),
        staleTime: 30_000,
      })
    }

    if (!sidFromToday && sessionId) {
      setSessionId(null)
      setCurrentSessionId(null)
      try { setSessionForDisciple(id!, null) } catch {}
    }
    // Importante: NO dependas de sessionId para evitar loops
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sidFromToday, id, displayName])

  // SID efectivo: usa el local si existe, si no el del backend
  const effectiveSid = useMemo(
    () => sessionId ?? sidFromToday ?? null,
    [sessionId, sidFromToday]
  )

  // Rango del gráfico
  const [range, setRange] = useState(ov?.pivot?.days ?? 14)
  useEffect(() => { if (ov?.pivot?.days) setRange(ov.pivot.days) }, [ov?.pivot?.days])

  const chartPreviewTotal = useMemo(() => {
    const p = ov?.pivot
    if (!p) return 0
    const { data } = pivotToChartData(p, range)
    return data.reduce((acc, r) => {
      return acc + Object.entries(r)
        .filter(([k]) => k !== 'date')
        .reduce((s, [, v]) => s + (Number(v) || 0), 0)
    }, 0)
  }, [ov?.pivot, range])

  // Modal set
  const [open, setOpen] = useState(false)
  const [selected, setSelected] = useState<Prescription | null>(null)

  // Mutación: asegura sesión, calcula next set_index y registra set
  const mLog = useMutation({
    mutationFn: async (vals: { reps: number; weight?: number | null; rpe?: number | null; notes?: string | null }) => {
      if (!selected || !today?.day?.id || !today?.assignment_id) {
        throw new Error('Faltan datos de sesión (assignment/day)')
      }

      // Asegura sesión
      let sid = effectiveSid
      if (!sid) {
        const created = await startSession({
          assignment_id: today.assignment_id,
          day_id: today.day.id as string,
        })
        sid = created.id
        setSessionId(sid)
        setCurrentSessionId(sid)
        setCurrentDisciple(id || null, displayName || null)
        try { setSessionForDisciple(id!, sid) } catch {}
      }

      // Calcula índice leyendo los sets actuales de la sesión
      const existing = await listSets(sid!)
      const nextIdx = nextSetIndexForPrescription(existing, selected.id)

      await addSet(sid!, {
        prescription_id: selected.id,
        set_index: nextIdx,
        reps: vals.reps,
        weight: vals.weight ?? null,
        rpe: vals.rpe ?? null,
        to_failure: false,
      })

      return { sid: sid! }
    },
    onSuccess: async ({ sid }) => {
      show({ type: 'success', message: 'Set registrado' })
      await qc.invalidateQueries({ queryKey: ['disciple', id, 'today'] })
      await qc.invalidateQueries({ queryKey: ['session', sid, 'sets'] })
      await qc.refetchQueries({ queryKey: ['session', sid, 'sets'] })
      setOpen(false)
      setSelected(null)
    },
    onError: () => show({ type: 'error', message: 'No se pudo registrar el set' }),
  })

  // Sets del día (si hay sesión efectiva)
  const setsQ = useQuery({
    queryKey: ['session', effectiveSid, 'sets'],
    queryFn: () => listSets(effectiveSid as string),
    enabled: !!effectiveSid,
    staleTime: 15_000,
  })

  const isLoadingAll = overviewQ.isLoading || todayQ.isLoading

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold">Discípulo: {displayName}</h2>
          {sidFromToday && (
            <div className="mt-1 text-xs text-gray-600 dark:text-neutral-300">
              Sesión activa · sets: {today?.current_session_sets_count ?? 0}
            </div>
          )}
        </div>

        <div className="flex items-center gap-2">
          {sidFromToday ? (
            <NavLink
              to={`/sessions/${sidFromToday}`}
              className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800 text-blue-600"
            >
              Ir a sesión
            </NavLink>
          ) : (
            <span className="text-sm text-gray-500">Sin sesión activa</span>
          )}
          <NavLink to="/dashboard" className="text-sm text-blue-600 hover:underline">
            Volver
          </NavLink>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
          <div className="font-semibold mb-2">Overview</div>
          {overviewQ.isError && <div className="text-red-600">Error al cargar overview</div>}
          {isLoadingAll && <div className="h-24 bg-gray-100 dark:bg-neutral-800 rounded" />}
          {!isLoadingAll && ov && (
            <>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                <KpiCard label="Adherencia" value={formatAdherence(ov, 1)} sub="Últimos 14 días" />
                <KpiCard
                  label="Días con sets"
                  value={`${ov.adherence && typeof ov.adherence === 'object' ? ov.adherence.days_with_sets : 0}`}
                  sub={`de ${ov.adherence && typeof ov.adherence === 'object' ? ov.adherence.days : 0}`}
                />
                <KpiCard label="Volumen total" value={`${chartPreviewTotal}`} sub={`${range} días`} />
              </div>

              <div className="mt-4 space-y-2">
                <div className="flex items-center justify-between">
                  <div className="font-medium">Últimos {range} días</div>
                  <RangePicker value={range} onChange={setRange} />
                </div>
                <OverviewVolumeChart
                  overview={{ ...ov, pivot: ov.pivot ? { ...ov.pivot, days: range } : ov.pivot }}
                />
              </div>
            </>
          )}
        </div>

        <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
          <div className="font-semibold mb-2">Hoy</div>
          {todayQ.isError && <div className="text-red-600">Error al cargar hoy</div>}
          {isLoadingAll && <div className="h-24 bg-gray-100 dark:bg-neutral-800 rounded" />}
          {!isLoadingAll && today && (
            <div className="text-sm space-y-2">
              <div>Índice de día: {today.day?.day_index ?? '-'}</div>
              {today.day?.notes && <div>Notas: {today.day?.notes}</div>}
              <div className="font-medium mt-2">Prescripciones</div>
              {renderPrescriptions(today.prescriptions, (p) => { setSelected(p); setOpen(true) })}
            </div>
          )}

          {effectiveSid && (
            <div className="mt-4">
              <div className="font-medium mb-1">Sets de hoy</div>
              {setsQ.isLoading && <div className="text-sm">Cargando sets…</div>}
              {setsQ.isError && <div className="text-sm text-red-600">No se pudieron cargar los sets</div>}
              {Array.isArray(setsQ.data) && setsQ.data.length > 0 ? (
                <TodaySets sets={setsQ.data} prescriptions={today?.prescriptions ?? []} />
              ) : (
                <div className="text-xs text-gray-500">Aún no hay sets en esta sesión</div>
              )}
            </div>
          )}
        </div>
      </div>

      <Modal
        open={open}
        onClose={() => { if (!mLog.isPending) { setOpen(false); setSelected(null) } }}
        title={selected ? `Registrar set — ${selected.exercise_name}` : 'Registrar set'}
      >
        <LogSetForm
          defaultValues={{ reps: 10 }}
          onCancel={() => { if (!mLog.isPending) { setOpen(false); setSelected(null) } }}
          onSubmit={(vals) => mLog.mutateAsync(vals)}
        />
      </Modal>
    </div>
  )
}
