import { useEffect, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate, Link } from 'react-router-dom'
import { useToast } from '@/components/toast/ToastProvider'
import Modal from '@/components/ui/Modal'
import { getMyActiveAssignment, listAssignmentDays } from '@/services/assignments'
import { startSession } from '@/services/sessions'
import type { AssignmentDay } from '@/types/assignments'
import { getMyActiveSession } from '@/services/sessions'
import { getProgram, listPrescriptions } from '@/services/programs'


export default function SessionsIndex() {
  const { show } = useToast()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const [openPicker, setOpenPicker] = useState(false)
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

  // 3) Días del programa activo — solo cuando abras el modal
  const daysQ = useQuery({
    queryKey: ['assignment', activeAssignQ.data?.id, 'days'],
    queryFn: () => listAssignmentDays(activeAssignQ.data!.id),
    enabled: openPicker && !!activeAssignQ.data?.id,
    staleTime: 0,
  })

  const prescQ = useQuery({
    queryKey: ['days', selectedDay?.id, 'prescriptions', 'start-session'],
    queryFn: () => listPrescriptions(selectedDay!.id),
    enabled: openPicker && !!selectedDay?.id,
    staleTime: 0,
  })

  useEffect(() => {
    setSelectedDay(null)
    if (openPicker && activeAssignQ.data?.id) {
      qc.invalidateQueries({ queryKey: ['assignment', activeAssignQ.data.id, 'days'] })
    }
  }, [activeAssignQ.data?.id, openPicker, qc])

  // Crear sesión
  const mStart = useMutation({
    mutationFn: (v: { assignment_id: string; day_id: string }) => startSession(v),
    onSuccess: async (sess) => {
      await qc.invalidateQueries({ queryKey: ['me', 'session', 'active'] })
      navigate(`/sessions/${sess.id}`)
    },
    onError: () => show({ type: 'error', message: 'No se pudo iniciar la sesión' }),
  })

  const haveActiveSession = !!activeSessQ.data
  const activeAssignment = activeAssignQ.data
  const days = daysQ.data ?? []
  const prescriptions = prescQ.data ?? []

  return (
    <div className="space-y-4">
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-semibold">Sesiones</h1>
            <div className="text-sm text-gray-600 dark:text-neutral-300">
              {haveActiveSession
                ? 'Tienes una sesión activa.'
                : 'No hay sesión activa. Puedes iniciar una nueva.'}
            </div>
          </div>

          <div className="flex items-center gap-2">
            {haveActiveSession ? (
              <Link
                to={`/sessions/${activeSessQ.data?.id}`}
                className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
              >
                Ir a mi sesión
              </Link>
            ) : (
              <button
                onClick={() => {
                  setSelectedDay(null)
                  setOpenPicker(true)
                }}
                disabled={activeAssignQ.isLoading || !activeAssignQ.data}
                className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                title={!activeAssignQ.data ? 'No tienes programa activo' : ''}
              >
                Iniciar nueva sesión
              </button>
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
                <div className="font-semibold">
                  {activeProgramQ.data?.title ?? activeAssignment.program_id}
                </div>
                <div className="text-xs text-gray-600 dark:text-neutral-300">
                  Inicio: {activeAssignment.start_date}
                </div>
              </>
            ) : (
              <div className="text-sm text-gray-500">No tienes rutina activa.</div>
            )}
          </div>
          {activeAssignment && (
            <span className="text-xs rounded border px-2 py-1 dark:border-neutral-800">Rutina activa</span>
          )}
        </div>
      </div>

      {/* Modal: elegir día del programa activo */}
      <Modal open={openPicker} onClose={() => setOpenPicker(false)} title="Seleccionar día de la rutina activa">
        {!activeAssignment ? (
          <div className="text-sm text-gray-500">No tienes rutina activa.</div>
        ) : daysQ.isLoading ? (
          <div className="text-sm text-gray-500">Cargando días…</div>
        ) : daysQ.isError ? (
          <div className="text-sm text-red-600">No se pudieron cargar los días de la rutina activa.</div>
        ) : !days.length ? (
          <div className="text-sm text-gray-500">Este programa no tiene días configurados.</div>
        ) : (
          <div className="space-y-4">
            <ul className="space-y-2">
              {days.map((d) => (
                <li key={d.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                  <div className="flex items-center justify-between gap-3">
                    <button
                      type="button"
                      onClick={() => setSelectedDay(d)}
                      className={`text-left flex-1 rounded px-2 py-1 ${selectedDay?.id === d.id ? 'bg-black text-white' : 'hover:bg-gray-50 dark:hover:bg-neutral-800'}`}
                    >
                      <div className="text-sm font-medium">
                        {d.title?.trim() ? d.title : `Día ${d.day_index}`}
                      </div>
                      <div className="text-xs opacity-80">
                        {d.exercise_names?.length ? d.exercise_names.join(', ') : 'Sin ejercicios en resumen'}
                      </div>
                    </button>
                    <span className="text-xs text-gray-500">Semana {d.week_index}</span>
                  </div>
                </li>
              ))}
            </ul>

            {selectedDay && (
              <div className="rounded border p-3 dark:border-neutral-800">
                <div className="font-medium text-sm">
                  {selectedDay.title?.trim() ? selectedDay.title : `Día ${selectedDay.day_index}`}
                </div>
                {prescQ.isLoading ? (
                  <div className="text-sm text-gray-500 mt-2">Cargando ejercicios…</div>
                ) : prescriptions.length ? (
                  <ul className="mt-2 space-y-1">
                    {prescriptions.map((p) => (
                      <li key={p.id} className="text-sm text-gray-700 dark:text-neutral-200">
                        {p.exercise_name ?? p.exercise_id} · {p.series}x{p.reps}
                      </li>
                    ))}
                  </ul>
                ) : (
                  <div className="text-sm text-gray-500 mt-2">Este día no tiene ejercicios configurados.</div>
                )}

                <div className="mt-3">
                  <button
                    onClick={() => mStart.mutate({ assignment_id: activeAssignment.id, day_id: selectedDay.id })}
                    disabled={mStart.isPending || prescQ.isLoading || !prescriptions.length}
                    className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 disabled:opacity-50 dark:bg-neutral-900 dark:border-neutral-800"
                  >
                    {mStart.isPending ? 'Iniciando…' : 'Iniciar sesión'}
                  </button>
                </div>
              </div>
            )}

            {!selectedDay && (
              <div className="text-sm text-gray-500">Selecciona un día para ver sus ejercicios.</div>
            )}
          </div>
        )}
      </Modal>
    </div>
  )
}
