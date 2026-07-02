import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createCheckin, listCheckins } from '@/services/checkins'
import { useToast } from '@/components/toast/ToastProvider'

function todayISO() {
  return new Date().toISOString().slice(0, 10)
}

function formatDate(value: string) {
  if (!value) return '-'
  return new Date(value).toLocaleDateString('es-CL', { year: 'numeric', month: 'short', day: 'numeric' })
}

export default function Checkins() {
  const qc = useQueryClient()
  const { show } = useToast()
  const [checkedAt, setCheckedAt] = useState(todayISO())
  const [weight, setWeight] = useState('')
  const [notes, setNotes] = useState('')

  const q = useQuery({
    queryKey: ['checkins'],
    queryFn: listCheckins,
    staleTime: 15_000,
  })

  const canSubmit = useMemo(() => {
    if (!checkedAt) return false
    if (weight.trim() === '') return true
    const parsed = Number(weight)
    return Number.isFinite(parsed) && parsed > 0
  }, [checkedAt, weight])

  const createM = useMutation({
    mutationFn: () => createCheckin({
      checked_at: checkedAt,
      weight_kg: weight.trim() === '' ? null : Number(weight),
      notes,
    }),
    onSuccess: async () => {
      setWeight('')
      setNotes('')
      show({ type: 'success', message: 'Check-in guardado' })
      await qc.invalidateQueries({ queryKey: ['checkins'] })
    },
    onError: () => show({ type: 'error', message: 'No se pudo guardar el check-in' }),
  })

  return (
    <div className="mx-auto max-w-4xl space-y-4">
      <div>
        <h2 className="text-xl font-semibold">Check-ins</h2>
        <p className="text-sm text-gray-600 dark:text-neutral-300">
          Registra peso y notas de seguimiento.
        </p>
      </div>

      <section className="rounded border bg-white p-4 dark:bg-neutral-900 dark:border-neutral-800">
        <div className="mb-3 font-semibold">Nuevo check-in</div>
        <div className="grid gap-3 sm:grid-cols-2">
          <label className="text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Fecha</span>
            <input
              type="date"
              value={checkedAt}
              onChange={(e) => setCheckedAt(e.target.value)}
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
            />
          </label>
          <label className="text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Peso kg</span>
            <input
              type="number"
              min="0"
              step="0.1"
              value={weight}
              onChange={(e) => setWeight(e.target.value)}
              placeholder="Opcional"
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
            />
          </label>
          <label className="sm:col-span-2 text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Notas</span>
            <textarea
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              rows={4}
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
              placeholder="Cómo te sentiste, sueño, energía o comentarios para tu seguimiento."
            />
          </label>
        </div>
        {!canSubmit && <div className="mt-2 text-sm text-red-600">Revisa la fecha o el peso ingresado.</div>}
        <button
          type="button"
          disabled={!canSubmit || createM.isPending}
          onClick={() => createM.mutate()}
          className="mt-3 rounded bg-black px-3 py-2 text-sm text-white disabled:cursor-not-allowed disabled:opacity-50"
        >
          {createM.isPending ? 'Guardando...' : 'Guardar check-in'}
        </button>
      </section>

      <section className="rounded border bg-white p-4 dark:bg-neutral-900 dark:border-neutral-800">
        <div className="mb-3 font-semibold">Mis check-ins</div>
        {q.isLoading && <div className="text-sm text-gray-500">Cargando check-ins...</div>}
        {q.isError && <div className="text-sm text-red-600">No se pudieron cargar los check-ins.</div>}
        {!q.isLoading && !q.isError && (q.data?.items ?? []).length === 0 && (
          <div className="text-sm text-gray-500">Aún no tienes check-ins.</div>
        )}
        <ul className="space-y-2">
          {(q.data?.items ?? []).map((item) => (
            <li key={item.id} className="rounded border px-3 py-3 dark:border-neutral-800">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <div className="text-sm font-medium">{formatDate(item.checked_at)}</div>
                  <div className="mt-1 text-sm text-gray-700 dark:text-neutral-300">
                    {item.notes?.trim() || 'Sin notas'}
                  </div>
                </div>
                <div className="rounded border px-2 py-1 text-xs dark:border-neutral-800">
                  {item.weight_kg ? `${item.weight_kg} kg` : 'Sin peso'}
                </div>
              </div>
            </li>
          ))}
        </ul>
      </section>
    </div>
  )
}
