import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate, Link } from 'react-router-dom'
import { useToast } from '@/components/toast/ToastProvider'
import Modal from '@/components/ui/Modal'
import { getMyActiveAssignment, listAssignmentDays } from '@/services/assignments'
import { startSession } from '@/services/sessions'
import type { AssignmentDay } from '@/types/assignments'
import { getMyActiveSession } from '@/services/sessions'


export default function SessionsIndex() {
  const { show } = useToast()
  const navigate = useNavigate()
  const qc = useQueryClient()
  const [openPicker, setOpenPicker] = useState(false)

  // 1) Sesión activa (si existe). 404 => null
  const activeSessQ = useQuery({
    queryKey: ['me', 'session', 'active'],
    queryFn: getMyActiveSession,
    retry: false,
    staleTime: 10_000,
  })

  // 2) Programa activo del discípulo logueado
  const activeAssignQ = useQuery({
    queryKey: ['me', 'assignment', 'active'],
    queryFn: getMyActiveAssignment,
    staleTime: 30_000,
  })

  // 3) Días del programa activo — solo cuando abras el modal
  const daysQ = useQuery({
    queryKey: ['assignment', activeAssignQ.data?.id, 'days'],
    queryFn: () => listAssignmentDays(activeAssignQ.data!.id),
    enabled: openPicker && !!activeAssignQ.data?.id,
    staleTime: 30_000,
  })

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
                onClick={() => setOpenPicker(true)} // <- abre modal; ahí se piden días del programa activo
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

      {/* Modal: elegir día del programa activo */}
      <Modal open={openPicker} onClose={() => setOpenPicker(false)} title="Seleccionar día del programa activo">
        {!activeAssignQ.data ? (
          <div className="text-sm text-gray-500">No tienes un programa activo.</div>
        ) : daysQ.isLoading ? (
          <div className="text-sm text-gray-500">Cargando días…</div>
        ) : !(daysQ.data ?? []).length ? (
          <div className="text-sm text-gray-500">Este programa no tiene días configurados.</div>
        ) : (
          <ul className="space-y-2">
            {(daysQ.data as AssignmentDay[]).map((d) => (
              <li key={d.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                <div className="flex items-center justify-between">
                  <div>
                    <div className="text-sm font-medium">
                      {d.title?.trim() ? d.title : `Día ${d.day_index}`}
                    </div>
                    <div className="text-xs text-gray-600 dark:text-neutral-300">
                      {d.exercise_names?.length ? d.exercise_names.join(', ') : '—'}
                    </div>
                  </div>
                  <button
                    onClick={() => mStart.mutate({ assignment_id: activeAssignQ.data!.id, day_id: d.id })}
                    disabled={mStart.isPending}
                    className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                  >
                    {mStart.isPending ? 'Iniciando…' : 'Entrenar aquí'}
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </Modal>
    </div>
  )
}
