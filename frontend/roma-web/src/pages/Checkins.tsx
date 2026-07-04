import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createCheckin, listCheckins, updateCheckin, type Checkin } from '@/services/checkins'
import { useToast } from '@/components/toast/ToastProvider'

function todayISO() {
  return new Date().toISOString().slice(0, 10)
}

function formatDate(value: string) {
  if (!value) return '-'
  return new Date(value).toLocaleDateString('es-CL', { year: 'numeric', month: 'short', day: 'numeric' })
}

function dateInputValue(value: string) {
  return value ? value.slice(0, 10) : ''
}

export default function Checkins() {
  const qc = useQueryClient()
  const { show } = useToast()
  const [checkedAt, setCheckedAt] = useState(todayISO())
  const [weight, setWeight] = useState('')
  const [notes, setNotes] = useState('')
  const [from, setFrom] = useState('')
  const [to, setTo] = useState('')
  const [limit, setLimit] = useState(10)
  const [page, setPage] = useState(1)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editCheckedAt, setEditCheckedAt] = useState('')
  const [editWeight, setEditWeight] = useState('')
  const [editNotes, setEditNotes] = useState('')
  const offset = (page - 1) * limit

  const q = useQuery({
    queryKey: ['checkins', from, to, limit, offset],
    queryFn: () => listCheckins({ from, to, limit, offset }),
    staleTime: 15_000,
  })
  const total = q.data?.total ?? 0
  const totalPages = Math.max(1, Math.ceil(total / limit))
  const canPrev = page > 1
  const canNext = page < totalPages

  const canSubmit = useMemo(() => {
    if (!checkedAt) return false
    if (weight.trim() === '') return true
    const parsed = Number(weight)
    return Number.isFinite(parsed) && parsed > 0
  }, [checkedAt, weight])
  const canSaveEdit = useMemo(() => {
    if (!editCheckedAt) return false
    if (editWeight.trim() === '') return true
    const parsed = Number(editWeight)
    return Number.isFinite(parsed) && parsed > 0
  }, [editCheckedAt, editWeight])

  function startEdit(item: Checkin) {
    setEditingId(item.id)
    setEditCheckedAt(dateInputValue(item.checked_at))
    setEditWeight(item.weight_kg ? String(item.weight_kg) : '')
    setEditNotes(item.notes ?? '')
  }

  function cancelEdit() {
    setEditingId(null)
    setEditCheckedAt('')
    setEditWeight('')
    setEditNotes('')
  }

  const createM = useMutation({
    mutationFn: () => createCheckin({
      checked_at: checkedAt,
      weight_kg: weight.trim() === '' ? null : Number(weight),
      notes,
    }),
    onSuccess: async () => {
      setWeight('')
      setNotes('')
      setPage(1)
      show({ type: 'success', message: 'Check-in guardado' })
      await qc.invalidateQueries({ queryKey: ['checkins'] })
    },
    onError: () => show({ type: 'error', message: 'No se pudo guardar el check-in' }),
  })
  const updateM = useMutation({
    mutationFn: (id: string) => updateCheckin(id, {
      checked_at: editCheckedAt,
      weight_kg: editWeight.trim() === '' ? null : Number(editWeight),
      notes: editNotes,
    }),
    onSuccess: async () => {
      cancelEdit()
      show({ type: 'success', message: 'Check-in actualizado' })
      await qc.invalidateQueries({ queryKey: ['checkins'] })
    },
    onError: () => show({ type: 'error', message: 'No se pudo actualizar el check-in' }),
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
        <div className="mb-3 flex flex-wrap items-center justify-between gap-3">
          <div className="font-semibold">Mis check-ins</div>
          <div className="text-sm text-gray-600 dark:text-neutral-300">
            {total} resultado{total === 1 ? '' : 's'}
          </div>
        </div>
        <div className="mb-4 grid gap-3 sm:grid-cols-3">
          <label className="text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Desde</span>
            <input
              type="date"
              value={from}
              onChange={(e) => { setFrom(e.target.value); setPage(1) }}
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
            />
          </label>
          <label className="text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Hasta</span>
            <input
              type="date"
              value={to}
              onChange={(e) => { setTo(e.target.value); setPage(1) }}
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
            />
          </label>
          <label className="text-sm">
            <span className="mb-1 block text-gray-600 dark:text-neutral-300">Por página</span>
            <select
              value={limit}
              onChange={(e) => { setLimit(Number(e.target.value)); setPage(1) }}
              className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
            >
              <option value={5}>5</option>
              <option value={10}>10</option>
              <option value={25}>25</option>
              <option value={50}>50</option>
            </select>
          </label>
        </div>
        <div className="mb-3 flex flex-wrap gap-2 text-xs text-gray-600 dark:text-neutral-300">
          <span className="rounded border px-2 py-1 dark:border-neutral-800">
            {from ? `Desde ${from}` : 'Sin fecha inicial'}
          </span>
          <span className="rounded border px-2 py-1 dark:border-neutral-800">
            {to ? `Hasta ${to}` : 'Sin fecha final'}
          </span>
          <span className="rounded border px-2 py-1 dark:border-neutral-800">
            Página {page} de {totalPages}
          </span>
        </div>
        {q.isLoading && <div className="text-sm text-gray-500">Cargando check-ins...</div>}
        {q.isError && <div className="text-sm text-red-600">No se pudieron cargar los check-ins.</div>}
        {!q.isLoading && !q.isError && (q.data?.items ?? []).length === 0 && (
          <div className="text-sm text-gray-500">Aún no tienes check-ins.</div>
        )}
        <ul className="space-y-2">
          {(q.data?.items ?? []).map((item) => (
            <li key={item.id} className="rounded border px-3 py-3 dark:border-neutral-800">
              {editingId === item.id ? (
                <div className="space-y-3">
                  <div className="grid gap-3 sm:grid-cols-2">
                    <label className="text-sm">
                      <span className="mb-1 block text-gray-600 dark:text-neutral-300">Fecha</span>
                      <input
                        type="date"
                        value={editCheckedAt}
                        onChange={(e) => setEditCheckedAt(e.target.value)}
                        className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
                      />
                    </label>
                    <label className="text-sm">
                      <span className="mb-1 block text-gray-600 dark:text-neutral-300">Peso kg</span>
                      <input
                        type="number"
                        min="0"
                        step="0.1"
                        value={editWeight}
                        onChange={(e) => setEditWeight(e.target.value)}
                        placeholder="Opcional"
                        className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
                      />
                    </label>
                    <label className="sm:col-span-2 text-sm">
                      <span className="mb-1 block text-gray-600 dark:text-neutral-300">Notas</span>
                      <textarea
                        value={editNotes}
                        onChange={(e) => setEditNotes(e.target.value)}
                        rows={3}
                        className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
                      />
                    </label>
                  </div>
                  {!canSaveEdit && <div className="text-sm text-red-600">Revisa la fecha o el peso ingresado.</div>}
                  <div className="flex flex-wrap gap-2">
                    <button
                      type="button"
                      disabled={!canSaveEdit || updateM.isPending}
                      onClick={() => updateM.mutate(item.id)}
                      className="rounded bg-black px-3 py-2 text-sm text-white disabled:cursor-not-allowed disabled:opacity-50"
                    >
                      {updateM.isPending ? 'Guardando...' : 'Guardar cambios'}
                    </button>
                    <button
                      type="button"
                      disabled={updateM.isPending}
                      onClick={cancelEdit}
                      className="rounded border px-3 py-2 text-sm bg-white hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-neutral-900 dark:border-neutral-800"
                    >
                      Cancelar
                    </button>
                  </div>
                </div>
              ) : (
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <div className="text-sm font-medium">{formatDate(item.checked_at)}</div>
                    <div className="mt-1 text-sm text-gray-700 dark:text-neutral-300">
                      {item.notes?.trim() || 'Sin notas'}
                    </div>
                  </div>
                  <div className="flex flex-wrap items-start gap-2">
                    <div className="rounded border px-2 py-1 text-xs dark:border-neutral-800">
                      {item.weight_kg ? `${item.weight_kg} kg` : 'Sin peso'}
                    </div>
                    <button
                      type="button"
                      onClick={() => startEdit(item)}
                      className="rounded border px-2 py-1 text-xs bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                    >
                      Editar
                    </button>
                  </div>
                </div>
              )}
            </li>
          ))}
        </ul>
        <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
          <button
            type="button"
            disabled={!canPrev || q.isFetching}
            onClick={() => setPage((current) => Math.max(1, current - 1))}
            className="rounded border px-3 py-2 text-sm bg-white hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-neutral-900 dark:border-neutral-800"
          >
            Anterior
          </button>
          <div className="text-sm text-gray-600 dark:text-neutral-300">
            Mostrando {total === 0 ? 0 : offset + 1}-{Math.min(offset + limit, total)} de {total}
          </div>
          <button
            type="button"
            disabled={!canNext || q.isFetching}
            onClick={() => setPage((current) => current + 1)}
            className="rounded border px-3 py-2 text-sm bg-white hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-neutral-900 dark:border-neutral-800"
          >
            Siguiente
          </button>
        </div>
      </section>
    </div>
  )
}
