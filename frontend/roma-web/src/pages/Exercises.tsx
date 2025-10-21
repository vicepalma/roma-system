import { useEffect } from 'react'
import { useToast } from '@/components/toast/ToastProvider'
import Modal from '@/components/ui/Modal'
import ExerciseForm from '@/components/forms/ExerciseForm'
import { useMemo, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listExercises, createExercise, updateExercise, deleteExercise } from '@/services/exercises'
import type { Exercise } from '@/types/exercises'

export default function Exercises() {
  const { show } = useToast()
  const qc = useQueryClient()
  const [openCreate, setOpenCreate] = useState(false)
  const [edit, setEdit] = useState<Exercise | null>(null)
  const [q, setQ] = useState('')

  const listQ = useQuery({
    queryKey: ['exercises', 'list'],
    queryFn: listExercises,
    staleTime: 60_000,
  })

  const raw = listQ.data as any
  const all: Exercise[] = Array.isArray(raw) ? raw : (raw?.items ?? [])
  
  useEffect(() => {
    if (listQ.isError) show({ type: 'error', message: 'No se pudieron cargar los ejercicios' })
  }, [listQ.isError, show])

  const createM = useMutation({
    mutationFn: createExercise,
    onSuccess: async () => {
      show({ type: 'success', message: 'Ejercicio creado' })
      setOpenCreate(false)
      await qc.invalidateQueries({ queryKey: ['exercises', 'list'] })
    },
    onError: () => show({ type: 'error', message: 'No se pudo crear el ejercicio' }),
  })

  const updateM = useMutation({
    mutationFn: ({ id, patch }: { id: string; patch: Partial<Omit<Exercise, 'id'>> }) =>
      updateExercise(id, patch),
      onSuccess: async () => {
    setEdit(null); // cierra el modal
    await qc.invalidateQueries({ queryKey: ['exercises', 'list'] });
  },
  })

  const deleteM = useMutation({
    mutationFn: (id: string) => deleteExercise(id),
    onSuccess: async () => {
      show({ type: 'success', message: 'Ejercicio eliminado' })
      await qc.invalidateQueries({ queryKey: ['exercises', 'list'] })
    },
    onError: () => show({ type: 'error', message: 'No se pudo eliminar' }),
  })

  const rows = useMemo(() => {
    const term = q.trim().toLowerCase()
    if (!term) return all

    return all.filter((e: Exercise) => {
      const name = (e.name ?? '').toLowerCase()
      const muscle = (e.primary_muscle ?? '').toLowerCase()
      const equip = (e.equipment ?? '').toLowerCase()
      const notes = (e.notes ?? '').toLowerCase()
      const tags = Array.isArray(e.tags) ? e.tags.join(' ').toLowerCase() : ''
      return (
        name.includes(term) ||
        muscle.includes(term) ||
        equip.includes(term) ||
        notes.includes(term) ||
        tags.includes(term)
      )
    })
  }, [all, q])

    if (listQ.isLoading) {
    return (
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-semibold">Ejercicios</h2>
          <div className="h-8 w-32 rounded bg-gray-100 dark:bg-neutral-800" />
        </div>
        <div className="h-9 rounded bg-gray-100 dark:bg-neutral-800" />
        <div className="h-64 rounded border bg-white dark:bg-neutral-900 dark:border-neutral-800" />
      </div>
    )
  }
  
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Ejercicios</h2>
        <button
          onClick={() => setOpenCreate(true)}
          className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
        >
          Nuevo ejercicio
        </button>
      </div>

      <div className="flex items-center gap-2">
        <input
          value={q}
          onChange={(e) => setQ(e.target.value)}
          placeholder="Buscar por nombre / músculo / equipo"
          className="border rounded px-2 py-1 text-sm w-full md:w-72 dark:bg-neutral-900 dark:border-neutral-800"
        />
      </div>

      <div className="overflow-x-auto rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800">
        <table className="min-w-full text-sm">
          <thead className="bg-gray-50 dark:bg-neutral-800 text-gray-600 dark:text-neutral-300">
            <tr>
              <th className="text-left px-4 py-2">Nombre</th>
              <th className="text-left px-4 py-2">Músculo</th>
              <th className="text-left px-4 py-2">Equipo</th>
              <th className="text-left px-4 py-2">Notas</th>
              <th className="px-4 py-2"></th>
            </tr>
          </thead>
          <tbody>
            {rows.map((e: Exercise) => (
              <tr key={e.id} className="border-t dark:border-neutral-800">
                <td className="px-4 py-2 font-medium">{e.name}</td>
                <td className="px-4 py-2">{e.primary_muscle ?? '—'}</td>
                <td className="px-4 py-2">{e.equipment ?? '—'}</td>
                <td className="px-4 py-2 text-xs text-gray-600 dark:text-neutral-300">
                  {e.notes ?? '—'}
                </td>
                <td className="px-4 py-2 text-right">
                  <div className="inline-flex items-center gap-2">
                    <button
                      onClick={() => setEdit(e)}
                      className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                    >
                      Editar
                    </button>
                    <button
                      onClick={() => deleteM.mutate(e.id)}
                      className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800 text-red-600"
                    >
                      Eliminar
                    </button>
                  </div>
                </td>
              </tr>
            ))}
            {rows.length === 0 && (
              <tr><td colSpan={5} className="px-4 py-6 text-center text-gray-500">Sin ejercicios</td></tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Modal crear */}
      <Modal open={openCreate} onClose={() => setOpenCreate(false)} title="Nuevo ejercicio">
        <ExerciseForm
          onCancel={() => setOpenCreate(false)}
          onSubmit={(vals) => createM.mutateAsync(vals)}
          submitting={createM.isPending}
        />
      </Modal>

      {/* Modal editar */}
      <Modal open={!!edit} onClose={() => setEdit(null)} title={`Editar — ${edit?.name ?? ''}`}>
        {edit && (
          <ExerciseForm
            defaultValues={edit}
            onCancel={() => setEdit(null)}
            onSubmit={(vals) =>
              updateM.mutateAsync({
                id: edit.id,
                patch: {
                  ...vals,
                  notes: vals.notes ?? undefined,
                  primary_muscle: vals.primary_muscle ?? undefined,
                  equipment: vals.equipment ?? undefined,
                },
              })
            }
            submitting={updateM.isPending}
          />
        )}
      </Modal>
    </div>
  )
}
