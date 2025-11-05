// src/components/forms/LogSetForm.tsx
import { useState } from 'react'

type Values = {
  reps: number
  weight?: number | null
  rpe?: number | null
  notes?: string | null
  to_failure?: boolean
}

export default function LogSetForm({
  defaultValues,
  onSubmit,
  onCancel,
}: {
  defaultValues?: Partial<Values>
  onSubmit: (v: Values) => void | Promise<unknown>
  onCancel?: () => void
}) {
  const [v, setV] = useState<Values>({
    reps: defaultValues?.reps ?? 10,
    weight: defaultValues?.weight ?? null,
    rpe: defaultValues?.rpe ?? null,
    to_failure: defaultValues?.to_failure ?? false,
  })

  return (
    <form
      className="space-y-3"
      onSubmit={async (e) => {
        e.preventDefault()
        await onSubmit({
          reps: Number(v.reps),
          weight: v.weight != null ? Number(v.weight) : null,
          rpe: v.rpe != null ? Number(v.rpe) : null,
          notes: v.notes?.trim() ? v.notes : null,
          to_failure: !!v.to_failure,
        })
      }}
    >
      <div className="grid sm:grid-cols-2 gap-3">
        <label className="text-sm">
          <div className="mb-1">Reps *</div>
          <input
            type="number"
            min={1}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.reps}
            onChange={(e) => setV(s => ({ ...s, reps: Number(e.target.value) }))}
            required
          />
        </label>
        <label className="text-sm">
          <div className="mb-1">Peso</div>
          <input
            type="number"
            step="0.5"
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.weight ?? ''}
            onChange={(e) => setV(s => ({ ...s, weight: e.target.value === '' ? null : Number(e.target.value) }))}
          />
        </label>
        <label className="text-sm">
          <div className="mb-1">RPE</div>
          <input
            type="number"
            step="0.5"
            min={1}
            max={10}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.rpe ?? ''}
            onChange={(e) => setV(s => ({ ...s, rpe: e.target.value === '' ? null : Number(e.target.value) }))}
          />
        </label>

        <label className="text-sm inline-flex items-center gap-2 mt-1 sm:col-span-2">
          <input
            type="checkbox"
            checked={!!v.to_failure}
            onChange={(e) => setV(s => ({ ...s, to_failure: e.target.checked }))}
          />
          Llegó a fallo
        </label>
      </div>

      <div className="flex items-center gap-2 pt-1">
        <button type="submit" className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800">
          Guardar
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
