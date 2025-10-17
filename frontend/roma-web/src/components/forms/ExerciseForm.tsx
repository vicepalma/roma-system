import { useState } from 'react'

type Values = {
  name: string
  primary_muscle?: string | null
  equipment?: string | null
  notes?: string | null
}

export default function ExerciseForm({
  defaultValues,
  submitting,
  onCancel,
  onSubmit,
}: {
  defaultValues?: Values
  submitting?: boolean
  onCancel?: () => void
  onSubmit: (v: Values) => Promise<void> | void
}) {
  const [v, setV] = useState<Values>({
    name: defaultValues?.name ?? '',
    primary_muscle: defaultValues?.primary_muscle ?? '',
    equipment: defaultValues?.equipment ?? '',
    notes: defaultValues?.notes ?? '',
  })

  return (
    <form
      className="space-y-3"
      onSubmit={async (e) => {
        e.preventDefault()
        await onSubmit({
          name: v.name.trim(),
          primary_muscle: v.primary_muscle || null,
          equipment: v.equipment || null,
          notes: v.notes || null,
        })
      }}
    >
      <div className="grid sm:grid-cols-2 gap-3">
        <label className="text-sm">
          <div className="mb-1">Nombre *</div>
          <input
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.name}
            required
            onChange={(e) => setV((s) => ({ ...s, name: e.target.value }))}
          />
        </label>
        <label className="text-sm">
          <div className="mb-1">Músculo principal</div>
          <input
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.primary_muscle ?? ''}
            onChange={(e) => setV((s) => ({ ...s, primary_muscle: e.target.value }))}
          />
        </label>
        <label className="text-sm">
          <div className="mb-1">Equipo</div>
          <input
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.equipment ?? ''}
            onChange={(e) => setV((s) => ({ ...s, equipment: e.target.value }))}
          />
        </label>
        <label className="text-sm sm:col-span-2">
          <div className="mb-1">Notas</div>
          <textarea
            rows={3}
            className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            value={v.notes ?? ''}
            onChange={(e) => setV((s) => ({ ...s, notes: e.target.value }))}
          />
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
          <button
            type="button"
            onClick={onCancel}
            className="text-sm text-gray-600 hover:underline"
          >
            Cancelar
          </button>
        )}
      </div>
    </form>
  )
}
