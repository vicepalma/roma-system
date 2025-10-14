import { useState } from 'react'

export default function AddDayForm({ onSubmit, submitting }: {
  onSubmit: (v: { day_index: number; notes?: string | null }) => Promise<any> | void
  submitting?: boolean
}) {
  const [dayIndex, setDayIndex] = useState<number>(1)
  const [notes, setNotes] = useState<string>('')

  return (
    <form
      onSubmit={async (e) => { e.preventDefault(); await onSubmit({ day_index: dayIndex, notes: notes || null }) }}
      className="flex flex-wrap items-end gap-2"
    >
      <div>
        <label className="block text-xs text-gray-600 dark:text-neutral-300 mb-1">Índice de día</label>
        <input
          type="number"
          min={1}
          value={dayIndex}
          onChange={(e) => setDayIndex(Number(e.target.value))}
          className="rounded border px-3 py-2 w-32 dark:bg-neutral-900 dark:border-neutral-800"
        />
      </div>
      <div className="flex-1 min-w-[220px]">
        <label className="block text-xs text-gray-600 dark:text-neutral-300 mb-1">Notas</label>
        <input
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          placeholder="Opcional"
          className="rounded border px-3 py-2 w-full dark:bg-neutral-900 dark:border-neutral-800"
        />
      </div>
      <button disabled={submitting} className="rounded border px-3 py-2 text-sm dark:border-neutral-800">
        Agregar día
      </button>
    </form>
  )
}
