import { useState } from 'react'

export default function AddWeekForm({ onSubmit, submitting }: {
  onSubmit: (v: { index: number }) => Promise<any> | void
  submitting?: boolean
}) {
  const [index, setIndex] = useState<number>(1)
  return (
    <form
      onSubmit={async (e) => { e.preventDefault(); await onSubmit({ index }) }}
      className="flex items-end gap-2"
    >
      <div>
        <label className="block text-xs text-gray-600 dark:text-neutral-300 mb-1">Índice de semana</label>
        <input
          type="number"
          min={1}
          value={index}
          onChange={(e) => setIndex(Number(e.target.value))}
          className="rounded border px-3 py-2 w-32 dark:bg-neutral-900 dark:border-neutral-800"
        />
      </div>
      <button
        disabled={submitting}
        className="rounded border px-3 py-2 text-sm dark:border-neutral-800"
      >
        Agregar semana
      </button>
    </form>
  )
}
