import { useState } from 'react'

type Values = {
  exercise_id: string
  series: number
  reps: string
  rest_sec?: number | null
  to_failure?: boolean
  position: number
}

export default function PrescriptionForm({
  defaultValues,
  submitting,
  onCancel,
  onSubmit,
}: {
  defaultValues?: Partial<Values>
  submitting?: boolean
  onCancel?: () => void
  onSubmit: (v: Values) => Promise<void> | void
}) {
  const [v, setV] = useState<Values>({
    exercise_id: defaultValues?.exercise_id ?? '',
    series: Number(defaultValues?.series ?? 3),
    reps: defaultValues?.reps ?? '10-12',
    rest_sec: defaultValues?.rest_sec ?? 90,
    to_failure: defaultValues?.to_failure ?? false,
    position: Number(defaultValues?.position ?? 1),
  })

  return (
    <form
      className="space-y-3"
      onSubmit={async (e) => {
        e.preventDefault()
        await onSubmit({
          exercise_id: v.exercise_id.trim(),
          series: Number(v.series),
          reps: v.reps.trim(),
          rest_sec: v.rest_sec ? Number(v.rest_sec) : null,
          to_failure: !!v.to_failure,
          position: Number(v.position),
        })
      }}
    >
      <div className="grid sm:grid-cols-2 gap-3">
        <label className="text-sm">
          <div className="mb-1">Exercise ID *</div>
          <input
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.exercise_id}
            required
            onChange={(e) => setV(s => ({ ...s, exercise_id: e.target.value }))}
            placeholder="UUID del ejercicio"
          />
        </label>
        <label className="text-sm">
          <div className="mb-1">Series *</div>
          <input
            type="number"
            min={1}
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
            type="number"
            min={0}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.rest_sec ?? 0}
            onChange={(e) => setV(s => ({ ...s, rest_sec: Number(e.target.value) }))}
          />
        </label>
        <label className="text-sm">
          <div className="mb-1">Posición *</div>
          <input
            type="number"
            min={1}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.position}
            required
            onChange={(e) => setV(s => ({ ...s, position: Number(e.target.value) }))}
          />
        </label>
        <label className="text-sm inline-flex items-center gap-2 mt-6">
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
