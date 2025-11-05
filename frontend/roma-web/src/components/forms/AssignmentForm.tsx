import { useState } from 'react'
import type { ProgramOption } from '@/services/programs'

type FormValues = {
  disciple_id: string
  program_id: string
  program_version: number
  start_date: string
  end_date?: string | null
}

type Props = {
  disciples: { id: string; name: string; email: string }[]
  programs: ProgramOption[]
  submitting?: boolean
  onSubmit: (values: FormValues) => Promise<unknown> | void
}

export default function AssignmentForm({ disciples, programs, submitting, onSubmit }: Props) {
  const [v, setV] = useState<FormValues>({
    disciple_id: '',
    program_id: '',
    program_version: 1,
    start_date: new Date().toISOString().slice(0, 10),
    end_date: null,
  })

  const onProgramChange = (programId: string) => {
    const p = programs.find(x => x.id === programId)
    setV(s => ({
      ...s,
      program_id: programId,
      program_version: p?.version ?? 1,
    }))
  }

  return (
    <form
      className="grid sm:grid-cols-2 gap-3"
      onSubmit={async (e) => {
        e.preventDefault()
        await onSubmit({
          ...v,
          end_date: v.end_date && v.end_date.trim() !== '' ? v.end_date : null,
        })
      }}
    >
      {/* Discípulo */}
      <label className="text-sm">
        <div className="mb-1">Discípulo *</div>
        <select
          required
          className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
          value={v.disciple_id}
          onChange={(e) => setV(s => ({ ...s, disciple_id: e.target.value }))}
        >
          <option value="">Selecciona…</option>
          {disciples.map(d => (
            <option key={d.id} value={d.id}>{d.name} — {d.email}</option>
          ))}
        </select>
      </label>

      {/* Programa */}
      <label className="text-sm">
        <div className="mb-1">Programa *</div>
        <select
          required
          className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
          value={v.program_id}
          onChange={(e) => onProgramChange(e.target.value)}
        >
          <option value="">Selecciona…</option>
          {programs.map(p => (
            <option key={p.id} value={p.id}>
              {p.title} (v{p.version})
            </option>
          ))}
        </select>
        {/* Mostrar versión seleccionada (solo lectura) */}
        <div className="text-[11px] text-gray-500 mt-1">
          Versión seleccionada: v{v.program_version}
        </div>
      </label>

      {/* Fecha inicio */}
      <label className="text-sm">
        <div className="mb-1">Fecha de inicio *</div>
        <input
          type="date"
          required
          className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
          value={v.start_date}
          onChange={(e) => setV(s => ({ ...s, start_date: e.target.value }))}
        />
      </label>

      {/* Fecha fin (opcional) */}
      <label className="text-sm">
        <div className="mb-1">Fecha de término</div>
        <input
          type="date"
          className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
          value={v.end_date ?? ''}
          onChange={(e) => setV(s => ({ ...s, end_date: e.target.value || null }))}
        />
      </label>

      <div className="sm:col-span-2 flex items-center gap-2 pt-1">
        <button
          type="submit"
          disabled={submitting}
          className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
        >
          {submitting ? 'Creando…' : 'Crear asignación'}
        </button>
      </div>
    </form>
  )
}
