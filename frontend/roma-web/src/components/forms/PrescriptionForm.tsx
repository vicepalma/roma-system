import { useState } from 'react'

type Values = {
  series: number
  reps: string
  rest_sec?: number | null
  to_failure?: boolean
  position: number
  tempo?: string | null
  rir?: number | null
  rpe?: number | null
  notes?: string | null
}

export default function PrescriptionForm({
  exerciseId,
  defaultValues,
  submitting,
  onCancel,
  onSubmit,
}: {
  exerciseId: string
  defaultValues?: Partial<Values>
  submitting?: boolean
  onCancel?: () => void
  onSubmit: (v: { exercise_id: string } & Values & { method_id: null }) => Promise<void> | void
}) {
  const [v, setV] = useState<Values>({
    series: Number(defaultValues?.series ?? 3),
    reps: defaultValues?.reps ?? '10-12',
    rest_sec: defaultValues?.rest_sec ?? 90,
    to_failure: defaultValues?.to_failure ?? false,
    position: Number(defaultValues?.position ?? 1),
    tempo: defaultValues?.tempo ?? null,
    rir: defaultValues?.rir ?? null,
    rpe: defaultValues?.rpe ?? null,
    notes: defaultValues?.notes ?? null,
  })

  const toNullIfEmpty = (s: string) => (s.trim() === '' ? null : s.trim())
  const toNumOrNull = (s: string) => (s === '' ? null : Number(s))

  return (
    <form
      className="space-y-3"
      onSubmit={async (e) => {
        e.preventDefault()
        await onSubmit({
          exercise_id: exerciseId,
          series: Number(v.series),
          reps: v.reps.trim(),
          rest_sec: v.rest_sec == null ? null : Number(v.rest_sec),
          to_failure: !!v.to_failure,
          position: Number(v.position),
          tempo: v.tempo ? v.tempo.trim() : null,
          rir: v.rir == null ? null : Number(v.rir),
          rpe: v.rpe == null ? null : Number(v.rpe),
          notes: v.notes ? v.notes.trim() : null,
          method_id: null, // <- SIEMPRE null (oculto al usuario)
        })
      }}
    >
      <div className="grid sm:grid-cols-2 gap-3">
        <label className="text-sm">
          <div className="mb-1">Series *</div>
          <input
            type="number" min={1}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.series}
            required
            onChange={(e) => setV(s => ({ ...s, series: Number(e.target.value) }))}
          />
        </label>

        <label className="text-sm">
          <div className="mb-1">Reps *</div>
          <input
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.reps}
            required
            onChange={(e) => setV(s => ({ ...s, reps: e.target.value }))}
            placeholder="10-12"
          />
        </label>

        <label className="text-sm">
          <div className="mb-1">Descanso (seg)</div>
          <input
            type="number" min={0}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.rest_sec ?? ''}
            onChange={(e) => setV(s => ({ ...s, rest_sec: toNumOrNull(e.target.value) }))}
            placeholder="90"
          />
        </label>

        <label className="text-sm">
          <div className="mb-1">Posición *</div>
          <input
            type="number" min={1}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.position}
            required
            onChange={(e) => setV(s => ({ ...s, position: Number(e.target.value) }))}
          />
        </label>

        <label className="text-sm">
          <div className="mb-1">Tempo</div>
          <input
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.tempo ?? ''}
            onChange={(e) => setV(s => ({ ...s, tempo: toNullIfEmpty(e.target.value) }))}
            placeholder="p. ej. 3-1-1"
          />
        </label>

        <label className="text-sm">
          <div className="mb-1">RIR</div>
          <input
            type="number" min={0}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.rir ?? ''}
            onChange={(e) => setV(s => ({ ...s, rir: toNumOrNull(e.target.value) }))}
          />
        </label>

        <label className="text-sm">
          <div className="mb-1">RPE</div>
          <input
            type="number" step="0.5" min={0} max={10}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.rpe ?? ''}
            onChange={(e) => setV(s => ({ ...s, rpe: toNumOrNull(e.target.value) }))}
          />
        </label>

        <label className="text-sm sm:col-span-2">
          <div className="mb-1">Notas</div>
          <textarea
            rows={2}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.notes ?? ''}
            onChange={(e) => setV(s => ({ ...s, notes: toNullIfEmpty(e.target.value) }))}
          />
        </label>

        <label className="text-sm inline-flex items-center gap-2">
          <input
            type="checkbox"
            checked={!!v.to_failure}
            onChange={(e) => setV(s => ({ ...s, to_failure: e.target.checked }))}
          />
          A fallo
        </label>
      </div>

      <div className="flex items-center gap-2 pt-1">
        <button
          type="submit"
          disabled={submitting}
          className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
        >
          {submitting ? 'Guardando…' : 'Guardar'}
        </button>
        {onCancel && (
          <button type="button" onClick={onCancel} className="text-sm text-gray-600 hover:underline">
            Cancelar
          </button>
        )}
      </div>
    </form>
  )
}
