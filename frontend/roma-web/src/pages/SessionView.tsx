import { useEffect, useState } from 'react'
import { useParams, NavLink } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useToast } from '@/components/toast/ToastProvider'
import { getSession, listSets, addSet, deleteSet, endSession } from '@/services/sessions'
import type { SessionSet } from '@/types/sessions'
import LogSetForm, { type Values as LogValues } from '@/components/forms/LogSetForm'

export default function SessionView() {
  const { id = '' } = useParams()
  const { show } = useToast()
  const qc = useQueryClient()
  const [open, setOpen] = useState(false)
  const [selectedPrescriptionId, setSelectedPrescriptionId] = useState<string | null>(null)

  const sessionQ = useQuery({
    queryKey: ['session', id, 'meta'],
    queryFn: () => getSession(id),
    enabled: !!id,
  })

  const setsQ = useQuery({
    queryKey: ['session', id, 'sets'],
    queryFn: () => listSets(id),
    enabled: !!id,
    staleTime: 10_000,
  })

  useEffect(() => {
    if (sessionQ.isError) show({ type: 'error', message: 'Error al cargar sesión' })
  }, [sessionQ.isError, show])
  useEffect(() => {
    if (setsQ.isError) show({ type: 'error', message: 'Error al cargar sets' })
  }, [setsQ.isError, show])

  const mAdd = useMutation({
    mutationFn: async (vals: { prescription_id: string; reps: number; weight?: number | null; rpe?: number | null; to_failure?: boolean }) => {
      // …tu llamada a addSet(...) aquí
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['session'] })
      setOpen(false)
    },
  })

  const mDelete = useMutation({
    mutationFn: (setId: string) => deleteSet(id, setId),
    onSuccess: async () => {
      show({ type: 'success', message: 'Set eliminado' })
      await qc.invalidateQueries({ queryKey: ['session', id, 'sets'] })
    },
    onError: () => show({ type: 'error', message: 'No se pudo eliminar el set' }),
  })

  const mEnd = useMutation({
    mutationFn: () => endSession(id, { ended_at: new Date().toISOString() }),
    onSuccess: async () => {
      show({ type: 'success', message: 'Sesión finalizada' })
      await qc.invalidateQueries({ queryKey: ['session', id, 'meta'] })
      await qc.invalidateQueries({ queryKey: ['session', id, 'sets'] })
    },
    onError: () => show({ type: 'error', message: 'No se pudo finalizar la sesión' }),
  })

  const totalSets = (setsQ.data ?? []).length
  const started = sessionQ.data?.started_at
  const ended = sessionQ.data?.ended_at

  // UI mínima para agregar set: exige prescription_id
  const [prescriptionId, setPrescriptionId] = useState('')

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold">Sesión: {id}</h2>
          <div className="text-xs text-gray-600 dark:text-neutral-300">
            {started ? `Inicio: ${new Date(started).toLocaleString('es-CL')}` : 'Sin inicio'}
            {ended ? ` · Fin: ${new Date(ended).toLocaleString('es-CL')}` : ''}
            {' · '}Sets: {totalSets}
          </div>
        </div>
        <div className="flex items-center gap-2">
          <NavLink to="/dashboard" className="text-sm text-blue-600 hover:underline">Volver</NavLink>
          {!ended && (
            <button
              onClick={() => mEnd.mutate()}
              className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
            >
              Finalizar sesión
            </button>
          )}
        </div>
      </div>

      {/* Agregar set rápido */}
      {!ended && (
        <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4 space-y-3">
          <div className="font-semibold">Agregar set</div>
          <div className="grid sm:grid-cols-3 gap-2">
            <input
              className="border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
              placeholder="prescription_id"
              value={prescriptionId}
              onChange={(e) => setPrescriptionId(e.target.value)}
            />
            {/* Form nativo mínimo con LogSetForm */}
            <LogSetForm
              defaultValues={{ reps: 10 }}
              onCancel={() => setOpen(false)}
              onSubmit={(vals: LogValues) => {
                if (!selectedPrescriptionId) return
                return mAdd.mutateAsync({
                  prescription_id: selectedPrescriptionId,
                  reps: vals.reps,
                  weight: vals.weight ?? null,
                  rpe: vals.rpe ?? null,
                  to_failure: false,
                })
              }}
            />
          </div>
        </div>
      )}

      {/* Listado de sets */}
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800">
        <div className="p-4 font-semibold">Sets</div>
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-50 dark:bg-neutral-800 text-gray-600 dark:text-neutral-300">
              <tr>
                <th className="text-left px-4 py-2">#</th>
                <th className="text-left px-4 py-2">Prescription</th>
                <th className="text-left px-4 py-2">Reps</th>
                <th className="text-left px-4 py-2">Peso</th>
                <th className="text-left px-4 py-2">RPE</th>
                <th className="text-left px-4 py-2">Fail</th>
                <th className="px-4 py-2"></th>
              </tr>
            </thead>
            <tbody>
              {(setsQ.data ?? []).map((s: SessionSet) => (
                <tr key={s.id} className="border-t dark:border-neutral-800">
                  <td className="px-4 py-2">{s.set_index}</td>
                  <td className="px-4 py-2 font-mono text-[11px]">{s.prescription_id}</td>
                  <td className="px-4 py-2">{s.reps}</td>
                  <td className="px-4 py-2">{s.weight ?? '—'}</td>
                  <td className="px-4 py-2">{s.rpe ?? '—'}</td>
                  <td className="px-4 py-2">{s.to_failure ? 'Sí' : 'No'}</td>
                  <td className="px-4 py-2 text-right">
                    {!ended && (
                      <button
                        onClick={() => mDelete.mutate(s.id)}
                        className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                      >
                        Eliminar
                      </button>
                    )}
                  </td>
                </tr>
              ))}
              {(setsQ.data ?? []).length === 0 && (
                <tr><td colSpan={7} className="px-4 py-6 text-center text-gray-500">Sin sets</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
