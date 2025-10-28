import { useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useMemo, useState } from 'react'
import { useToast } from '@/components/toast/ToastProvider'
import Modal from '@/components/ui/Modal'
import LogSetForm from '@/components/forms/LogSetForm'
import RestTimer from '@/components/timers/RestTimer'

import { getSession, addSet, deleteSet } from '@/services/sessions'
import { listPrescriptions } from '@/services/programs'
import { listExercises } from '@/services/exercises'

export default function SessionView() {
  const { id = '' } = useParams()
  const { show } = useToast()
  const qc = useQueryClient()

  // Sesión
  const sessQ = useQuery({
    queryKey: ['session', id],
    queryFn: () => getSession(id),
    enabled: !!id,
  })

  const dayId: string | undefined = sessQ.data?.session?.day_id

  // Prescripciones del día
  const prescQ = useQuery({
    queryKey: ['day', dayId, 'prescriptions'],
    queryFn: () => listPrescriptions(dayId!),
    enabled: !!dayId,
    staleTime: 30_000,
  })

  // Catálogo de ejercicios para mostrar nombre legible
  const exQ = useQuery({
    queryKey: ['exercises', 'dict'],
    queryFn: listExercises,
    staleTime: 5 * 60_000,
  })
  const exById = useMemo(
    () => Object.fromEntries((exQ.data ?? []).map((e: any) => [e.id, e])),
    [exQ.data]
  )

  const sets = sessQ.data?.sets ?? []
  const cardio = sessQ.data?.cardio ?? []
  const presc = prescQ.data ?? []

  const [openLog, setOpenLog] = useState<{ open: boolean; prescription_id?: string }>({ open: false })

  // Mutations
  const mAdd = useMutation({
    mutationFn: (vals: {
      prescription_id: string
      reps: number
      weight?: number | null
      rpe?: number | null
      to_failure?: boolean
      notes?: string | null
    }) => addSet(id, vals),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['session', id] })
      setOpenLog({ open: false })
      show({ type: 'success', message: 'Set registrado' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo registrar el set' }),
  })

  const mDel = useMutation({
    mutationFn: (setId: string) => deleteSet(id, setId),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['session', id] })
      show({ type: 'success', message: 'Set eliminado' })
    },
  })

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="text-sm text-gray-600 dark:text-neutral-300">
            {sessQ.isLoading
              ? 'Cargando…'
              : (sessQ.data?.session?.performed_at
                  ? new Date(sessQ.data.session.performed_at).toLocaleString()
                  : '—')}
          </div>
          <h1 className="text-xl font-semibold">Sesión</h1>
        </div>
        {/* Reinicia el timer cuando cambia la cantidad de sets */}
        <RestTimer initial={180} key={sets.length} />
      </div>

      {/* Prescripciones */}
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
        <div className="font-semibold mb-2">
          Prescripciones {prescQ.isLoading ? '…' : ''}
        </div>
        {presc.length === 0 && !prescQ.isLoading && (
          <div className="text-sm text-gray-500">Sin prescripciones en este día</div>
        )}
        <ul className="mt-2 space-y-2">
          {presc.map((p: any) => (
            <li key={p.id} className="rounded border px-3 py-2 dark:border-neutral-800">
              <div className="flex items-center justify-between">
                <div className="text-sm font-medium">
                  {exById[p.exercise_id]?.name || p.exercise_name || p.exercise_id}
                </div>
                <button
                  onClick={() => setOpenLog({ open: true, prescription_id: p.id })}
                  className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                >
                  + Set
                </button>
              </div>
              <div className="text-xs text-gray-600 mt-1">
                Series: {p.series} · Reps: {p.reps}
                {p.rest_sec ? ` · Descanso: ${p.rest_sec}s` : ''}
                {p.to_failure ? ' · A fallo' : ''} · Posición: {p.position}
              </div>
            </li>
          ))}
        </ul>
      </div>

      {/* Sets */}
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
        <div className="font-semibold mb-2">Sets de la sesión</div>
        {sets.length === 0 ? (
          <div className="text-sm text-gray-500">Aún no hay sets</div>
        ) : (
          <ul className="mt-2 space-y-2">
            {sets.map((s: any) => (
              <li key={s.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                <div className="flex items-center justify-between">
                  <div className="text-sm">
                    <span className="font-medium">Set {s.set_index}</span> — Reps: {s.reps}
                    {s.weight != null ? ` · Peso: ${s.weight}` : ''}
                    {s.rpe != null ? ` · RPE: ${s.rpe}` : ''}
                    {s.to_failure ? ' · A fallo' : ''}
                  </div>
                  <button
                    onClick={() => mDel.mutate(s.id)}
                    className="text-xs text-red-600 rounded px-2 py-0.5 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                  >
                    Eliminar
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* Cardio */}
      {!!cardio.length && (
        <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
          <div className="font-semibold mb-2">Cardio</div>
          <ul className="mt-2 space-y-1 text-sm">
            {cardio.map((c: any) => (
              <li key={c.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                {c.modality} — {c.minutes} min
                {c.target_hr_min ? ` · HR ${c.target_hr_min}-${c.target_hr_max ?? ''}` : ''}
                {c.notes ? ` · ${c.notes}` : ''}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Modal registrar set */}
      <Modal open={openLog.open} onClose={() => setOpenLog({ open: false })} title="Registrar set">
        {openLog.prescription_id ? (
          <LogSetForm
            defaultValues={{ reps: 10 }}
            onCancel={() => setOpenLog({ open: false })}
            onSubmit={(vals) =>
              mAdd.mutateAsync({
                prescription_id: openLog.prescription_id!,
                reps: vals.reps,
                weight: vals.weight ?? null,
                rpe: vals.rpe ?? null,
                to_failure: !!vals.to_failure,
                notes: vals.notes ?? null,
              })
            }
          />
        ) : (
          <div className="text-sm text-gray-500">Selecciona una prescripción</div>
        )}
      </Modal>
    </div>
  )
}
