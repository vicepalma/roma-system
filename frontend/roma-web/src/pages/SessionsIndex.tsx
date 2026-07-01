import { useEffect, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate, Link } from 'react-router-dom'
import axios from 'axios'
import { useToast } from '@/components/toast/ToastProvider'
import { getMyActiveAssignment, listAssignmentDays } from '@/services/assignments'
import type { MyActiveAssignment } from '@/services/assignments'
import { startSession } from '@/services/sessions'
import type { AssignmentDay } from '@/types/assignments'
import { getMyActiveSession } from '@/services/sessions'
import { getProgram, listPrescriptions } from '@/services/programs'

function dayLabel(day: AssignmentDay) {
  return day.title?.trim() ? day.title : `Día ${day.day_index}`
}

function restLabel(seconds?: number | null) {
  if (!seconds) return null
  if (seconds < 60) return `${seconds}s descanso`
  const minutes = Math.floor(seconds / 60)
  const rest = seconds % 60
  return rest ? `${minutes}m ${rest}s descanso` : `${minutes}m descanso`
}

function formatDate(value?: string | null) {
  if (!value) return 'Sin fecha'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('es-CL', { day: '2-digit', month: 'short', year: 'numeric' }).format(date)
}

function routineType(assignment?: MyActiveAssignment | null, kind?: string) {
  if (assignment?.assigned_by && assignment.assigned_by === assignment.disciple_id) return 'Rutina propia'
  return kind === 'self_training' ? 'Rutina propia' : 'Rutina asignada por maestro'
}

function exerciseSummary(names?: string[]) {
  if (!names?.length) return 'Sin ejercicios configurados'
  if (names.length <= 3) return names.join(', ')
  return `${names.slice(0, 3).join(', ')}…`
}

export default function SessionsIndex() {
  const { show } = useToast()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const [selectedDay, setSelectedDay] = useState<AssignmentDay | null>(null)

  // 1) Sesión activa (si existe). 404 => null
  const activeSessQ = useQuery({
    queryKey: ['me', 'session', 'active'],
    queryFn: getMyActiveSession,
    retry: false,
    staleTime: 0,
    refetchOnMount: 'always',
  })

  // 2) Programa activo del discípulo logueado
  const activeAssignQ = useQuery({
    queryKey: ['me', 'assignment', 'active'],
    queryFn: getMyActiveAssignment,
    retry: false,
    staleTime: 0,
    refetchOnMount: 'always',
  })

  const activeProgramQ = useQuery({
    queryKey: ['programs', activeAssignQ.data?.program_id, 'active-summary'],
    queryFn: () => getProgram(activeAssignQ.data!.program_id),
    enabled: !!activeAssignQ.data?.program_id,
    staleTime: 30_000,
  })

  // 3) Días del assignment activo, fuente de verdad para iniciar sesión.
  const daysQ = useQuery({
    queryKey: ['assignment', activeAssignQ.data?.id, 'days'],
    queryFn: () => listAssignmentDays(activeAssignQ.data!.id),
    enabled: !!activeAssignQ.data?.id,
    staleTime: 0,
  })

  const prescQ = useQuery({
    queryKey: ['days', selectedDay?.id, 'prescriptions', 'start-session'],
    queryFn: () => listPrescriptions(selectedDay!.id),
    enabled: !!selectedDay?.id,
    staleTime: 0,
  })

  useEffect(() => {
    setSelectedDay(null)
    if (activeAssignQ.data?.id) {
      qc.invalidateQueries({ queryKey: ['assignment', activeAssignQ.data.id, 'days'] })
    }
  }, [activeAssignQ.data?.id, qc])

  // Crear sesión
  const mStart = useMutation({
    mutationFn: (v: { assignment_id: string; day_id: string }) => startSession(v),
    onSuccess: async (sess) => {
      await qc.invalidateQueries({ queryKey: ['me', 'session', 'active'] })
      navigate(`/sessions/${sess.id}`)
    },
    onError: (err) => {
      if (axios.isAxiosError(err) && err.response?.status === 409 && err.response.data?.error === 'assignment_inactive') {
        show({ type: 'error', message: 'Esta rutina ya no está activa. Activa la rutina antes de entrenar.' })
        qc.invalidateQueries({ queryKey: ['me', 'assignment', 'active'] })
        return
      }
      show({ type: 'error', message: 'No se pudo iniciar la sesión' })
    },
  })

  const haveActiveSession = !!activeSessQ.data
  const activeAssignment = activeAssignQ.data
  const days = daysQ.data ?? []
  const prescriptions = prescQ.data ?? []
  const activeProgramTitle = activeProgramQ.data?.title?.trim() || 'Rutina sin título'
  const activeProgramNotes = activeProgramQ.data?.notes?.trim()
  const activeRoutineType = routineType(activeAssignment, activeProgramQ.data?.kind)
  const activeSessionDay = activeSessQ.data?.day_id
    ? days.find((d) => d.id === activeSessQ.data?.day_id)
    : null
  const activeSessionDayLabel = activeSessionDay
    ? `Semana ${activeSessionDay.week_index} · Día ${activeSessionDay.day_index}`
    : 'Día de entrenamiento'

  return (
    <div className="space-y-4">
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
        <div className="flex items-start justify-between gap-3">
          <div>
            <div className="text-xs uppercase text-gray-500 dark:text-neutral-400">
              {haveActiveSession ? 'Sesión activa' : 'Entrenar'}
            </div>
            {haveActiveSession ? (
              <>
                <h1 className="text-xl font-semibold">
                  {activeProgramQ.isLoading ? 'Cargando rutina…' : activeProgramTitle || 'Rutina activa'}
                </h1>
                <div className="text-sm text-gray-600 dark:text-neutral-300">{activeSessionDayLabel}</div>
              </>
            ) : (
              <>
                <h1 className="text-xl font-semibold">Inicia una nueva sesión</h1>
                <div className="text-sm text-gray-600 dark:text-neutral-300">
                  Selecciona un día de tu rutina activa para comenzar.
                </div>
              </>
            )}
          </div>

          <div className="flex items-center gap-2">
            {haveActiveSession ? (
              <Link
                to={`/sessions/${activeSessQ.data?.id}`}
                className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
              >
                Continuar sesión
              </Link>
            ) : (
              <span className="text-xs rounded border px-2 py-1 text-gray-600 dark:text-neutral-300 dark:border-neutral-800">
                Sin sesión activa
              </span>
            )}
          </div>
        </div>
      </div>

      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
        <div className="flex items-start justify-between gap-3">
          <div>
            <div className="text-xs uppercase text-gray-500 dark:text-neutral-400">Rutina activa</div>
            {activeAssignQ.isLoading ? (
              <div className="text-sm text-gray-500">Cargando rutina…</div>
            ) : activeAssignment ? (
              <>
                <div className="font-semibold text-lg">
                  {activeProgramQ.isLoading ? 'Cargando rutina…' : activeProgramTitle}
                </div>
                <div className="mt-1 text-sm text-gray-600 dark:text-neutral-300">
                  {activeProgramNotes || 'Sin descripción'}
                </div>
                <div className="mt-2 flex flex-wrap gap-2 text-xs text-gray-600 dark:text-neutral-300">
                  <span className="rounded border px-2 py-1 dark:border-neutral-800">{activeRoutineType}</span>
                  <span className="rounded border px-2 py-1 dark:border-neutral-800">
                    Inicio: {formatDate(activeAssignment.start_date)}
                  </span>
                  <span className="rounded border px-2 py-1 dark:border-neutral-800">
                    {daysQ.isLoading ? 'Cargando días' : `${days.length} día${days.length === 1 ? '' : 's'} disponible${days.length === 1 ? '' : 's'}`}
                  </span>
                </div>
              </>
            ) : (
              <div className="text-sm text-gray-500">No tienes rutina activa.</div>
            )}
          </div>
          {activeAssignment && <span className="text-xs rounded border px-2 py-1 dark:border-neutral-800">Activa</span>}
        </div>
      </div>

      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
        <div className="mb-3">
          <h2 className="font-semibold">Días disponibles</h2>
          <div className="text-sm text-gray-600 dark:text-neutral-300">
            Elige un día para revisar ejercicios e iniciar entrenamiento.
          </div>
        </div>

        {activeAssignQ.isLoading ? (
          <div className="text-sm text-gray-500">Cargando rutina activa…</div>
        ) : !activeAssignment ? (
          <div className="text-sm text-gray-500">No tienes una rutina activa. Activa una desde Mis rutinas.</div>
        ) : daysQ.isLoading ? (
          <div className="text-sm text-gray-500">Cargando días de la rutina activa…</div>
        ) : daysQ.isError ? (
          <div className="text-sm text-red-600">No pudimos cargar los días de la rutina activa.</div>
        ) : days.length === 0 ? (
          <div className="text-sm text-gray-500">Esta rutina activa todavía no tiene días configurados.</div>
        ) : (
          <div className="grid gap-3 lg:grid-cols-[minmax(0,1fr)_minmax(280px,420px)]">
            <div className="space-y-2">
              {days.map((d) => {
                const selected = selectedDay?.id === d.id
                return (
                  <button
                    key={d.id}
                    type="button"
                    onClick={() => setSelectedDay(d)}
                    className={`w-full rounded border px-3 py-3 text-left dark:border-neutral-800 ${selected ? 'border-black bg-black text-white dark:border-white' : 'bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:hover:bg-neutral-800'}`}
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div>
                        <div className={`text-xs ${selected ? 'text-gray-200' : 'text-gray-500 dark:text-neutral-400'}`}>
                          Semana {d.week_index}
                        </div>
                        <div className="text-sm font-medium">Día {d.day_index}</div>
                        {d.notes?.trim() && (
                          <div className={`mt-1 text-xs ${selected ? 'text-gray-200' : 'text-gray-600 dark:text-neutral-300'}`}>
                            {d.notes.trim()}
                          </div>
                        )}
                        <div className={`mt-2 text-xs ${selected ? 'text-gray-200' : 'text-gray-500 dark:text-neutral-400'}`}>
                          {d.prescriptions_count} ejercicio{d.prescriptions_count === 1 ? '' : 's'}
                        </div>
                      </div>
                      {activeSessionDay?.id === d.id && (
                        <span className="text-xs rounded border px-2 py-1">En sesión</span>
                      )}
                    </div>
                    <div className={`mt-2 text-xs ${selected ? 'text-gray-200' : 'text-gray-600 dark:text-neutral-300'}`}>
                      {exerciseSummary(d.exercise_names)}
                    </div>
                  </button>
                )
              })}
            </div>

            <div className="rounded border p-3 dark:border-neutral-800">
              {selectedDay ? (
                <>
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="font-medium text-sm">{dayLabel(selectedDay)}</div>
                      <div className="text-xs text-gray-500">
                        Semana {selectedDay.week_index} · Día {selectedDay.day_index}
                      </div>
                      <div className="mt-1 text-xs text-gray-600 dark:text-neutral-300">
                        {selectedDay.notes?.trim() || 'Sin descripción'}
                      </div>
                    </div>
                    {!haveActiveSession && (
                      <button
                        onClick={() => mStart.mutate({ assignment_id: activeAssignment.id, day_id: selectedDay.id })}
                        disabled={mStart.isPending || prescQ.isLoading || !prescriptions.length}
                        className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 disabled:opacity-50 dark:bg-neutral-900 dark:border-neutral-800"
                      >
                        {mStart.isPending ? 'Iniciando…' : 'Iniciar sesión'}
                      </button>
                    )}
                  </div>

                  {haveActiveSession && (
                    <div className="mt-3 text-sm text-gray-500">
                      Ya tienes una sesión activa. Continúala antes de iniciar otra.
                    </div>
                  )}

                  {prescQ.isLoading ? (
                    <div className="text-sm text-gray-500 mt-3">Cargando ejercicios…</div>
                  ) : prescriptions.length ? (
                    <ul className="mt-3 divide-y dark:divide-neutral-800">
                      {prescriptions.map((p) => {
                        const rest = restLabel(p.rest_sec)
                        return (
                          <li key={p.id} className="py-2 text-sm">
                            <div className="font-medium text-gray-800 dark:text-neutral-100">
                              {p.exercise_name ?? p.exercise_id}
                            </div>
                            <div className="text-xs text-gray-600 dark:text-neutral-300">
                              {p.series} series · {p.reps} reps{rest ? ` · ${rest}` : ''}
                              {p.to_failure ? ' · al fallo' : ''}
                            </div>
                          </li>
                        )
                      })}
                    </ul>
                  ) : (
                    <div className="text-sm text-gray-500 mt-3">Este día todavía no tiene ejercicios.</div>
                  )}
                </>
              ) : (
                <div className="text-sm text-gray-500">Selecciona un día para ver sus ejercicios.</div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
