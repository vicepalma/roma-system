import { useState } from 'react'

type Props = {
  disciples: { id: string; name: string; email: string }[]
  onSubmit: (values: {
    disciple_id: string
    program_id: string
    program_version: number
    start_date: string
    end_date?: string | null
  }) => Promise<any> | any
  submitting?: boolean
}

export default function AssignmentForm({ disciples, onSubmit, submitting }: Props) {
  const [disciple_id, setDisciple] = useState('')
  const [program_id, setProgram] = useState('')
  const [program_version, setVersion] = useState<number>(1)
  const [start_date, setStart] = useState<string>(new Date().toISOString().slice(0, 10))
  const [end_date, setEnd] = useState<string>('')

  return (
    <form
      className="space-y-3"
      onSubmit={async (e) => {
        e.preventDefault()
        await onSubmit({
          disciple_id,
          program_id,
          program_version: Number(program_version) || 1,
          start_date,
          end_date: end_date ? end_date : null,
        })
      }}
    >
      <div className="grid gap-3 md:grid-cols-2">
        <label className="text-sm">
          <div className="mb-1">Discípulo</div>
          <select
            className="w-full rounded border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-2"
            value={disciple_id}
            onChange={(e) => setDisciple(e.target.value)}
            required
          >
            <option value="">Selecciona…</option>
            {disciples.map(d => (
              <option key={d.id} value={d.id}>{d.name} — {d.email}</option>
            ))}
          </select>
        </label>

        <label className="text-sm">
          <div className="mb-1">Programa ID</div>
          <input
            className="w-full rounded border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-2"
            value={program_id}
            onChange={(e) => setProgram(e.target.value)}
            placeholder="uuid del programa"
            required
          />
        </label>

        <label className="text-sm">
          <div className="mb-1">Versión</div>
          <input
            type="number"
            min={1}
            className="w-full rounded border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-2"
            value={program_version}
            onChange={(e) => setVersion(Number(e.target.value))}
            required
          />
        </label>

        <label className="text-sm">
          <div className="mb-1">Inicio</div>
          <input
            type="date"
            className="w-full rounded border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-2"
            value={start_date}
            onChange={(e) => setStart(e.target.value)}
            required
          />
        </label>

        <label className="text-sm">
          <div className="mb-1">Fin (opcional)</div>
          <input
            type="date"
            className="w-full rounded border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-2"
            value={end_date}
            onChange={(e) => setEnd(e.target.value)}
          />
        </label>
      </div>

      <div className="pt-2 flex justify-end gap-2">
        <button
          type="submit"
          disabled={submitting}
          className="rounded border px-3 py-2 bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800 text-sm"
        >
          {submitting ? 'Creando…' : 'Crear asignación'}
        </button>
      </div>
    </form>
  )
}
