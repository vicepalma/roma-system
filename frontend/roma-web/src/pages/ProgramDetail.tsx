import { useParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getProgramDetail, addWeek, addDay, addPrescription } from '@/services/programs'
import AddWeekForm from '@/components/programs/AddWeekForm'
import AddDayForm from '@/components/programs/AddDayForm'
import AddPrescriptionForm from '@/components/programs/AddPrescriptionForm'
import { useState } from 'react'
import { useToast } from '@/components/toast/ToastProvider'

export default function ProgramDetail() {
  const { id = '' } = useParams()
  const { show } = useToast()
  const qc = useQueryClient()

  const detailQ = useQuery({
    queryKey: ['programs', id, 'detail'],
    queryFn: () => getProgramDetail(id),
    enabled: !!id,
    staleTime: 15_000,
  })

  const addWeekM = useMutation({
    mutationFn: (v: { index: number }) => addWeek({ program_id: id, index: v.index }),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['programs', id, 'detail'] }); show({ type: 'success', message: 'Semana creada' }) },
    onError: () => show({ type: 'error', message: 'No se pudo crear la semana' }),
  })

  const [selectedWeekId, setSelectedWeekId] = useState<string | null>(null)
  const addDayM = useMutation({
    mutationFn: (v: { day_index: number; notes?: string | null }) => {
      if (!selectedWeekId) throw new Error('Selecciona una semana')
      return addDay({ week_id: selectedWeekId, day_index: v.day_index, notes: v.notes })
    },
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['programs', id, 'detail'] }); show({ type: 'success', message: 'Día creado' }) },
    onError: (e: any) => show({ type: 'error', message: e?.message || 'No se pudo crear el día' }),
  })

  const [selectedDayId, setSelectedDayId] = useState<string | null>(null)
  const addPresM = useMutation({
    mutationFn: (v: any) => {
      if (!selectedDayId) throw new Error('Selecciona un día')
      return addPrescription({ day_id: selectedDayId, ...v })
    },
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['programs', id, 'detail'] }); show({ type: 'success', message: 'Prescripción creada' }) },
    onError: (e: any) => show({ type: 'error', message: e?.message || 'No se pudo crear la prescripción' }),
  })

  const prog = detailQ.data
  const weeks = prog?.weeks ?? []

  return (
    <div className="mx-auto max-w-6xl p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Program Builder</h2>
        <div className="text-sm text-gray-600 dark:text-neutral-300">
          Programa: <span className="font-mono">{id}</span> {prog?.title ? `— ${prog.title} v${prog.version}` : null}
        </div>
      </div>

      <div className="grid md:grid-cols-3 gap-4">
        {/* Columna 1: Semanas */}
        <div className="rounded border p-4 dark:bg-neutral-900 dark:border-neutral-800">
          <div className="font-semibold mb-2">Semanas</div>
          <AddWeekForm onSubmit={v => addWeekM.mutateAsync(v)} submitting={addWeekM.isPending} />
          <ul className="mt-3 space-y-2">
            {weeks.map((w: any) => (
              <li key={w.id} className={`rounded border px-3 py-2 cursor-pointer dark:border-neutral-800 ${selectedWeekId === w.id ? 'bg-gray-50 dark:bg-neutral-800' : ''}`}
                  onClick={() => { setSelectedWeekId(w.id); setSelectedDayId(null) }}>
                <div className="font-medium">Semana #{w.index}</div>
                <div className="text-xs text-gray-500 font-mono">{w.id}</div>
              </li>
            ))}
            {weeks.length === 0 && <li className="text-sm text-gray-500">Aún sin semanas</li>}
          </ul>
        </div>

        {/* Columna 2: Días */}
        <div className="rounded border p-4 dark:bg-neutral-900 dark:border-neutral-800">
          <div className="font-semibold mb-2">Días</div>
          {selectedWeekId ? (
            <>
              <AddDayForm onSubmit={v => addDayM.mutateAsync(v)} submitting={addDayM.isPending} />
              <ul className="mt-3 space-y-2">
                {(weeks.find((w: any) => w.id === selectedWeekId)?.days ?? []).map((d: any) => (
                  <li key={d.id}
                      className={`rounded border px-3 py-2 cursor-pointer dark:border-neutral-800 ${selectedDayId === d.id ? 'bg-gray-50 dark:bg-neutral-800' : ''}`}
                      onClick={() => setSelectedDayId(d.id)}>
                    <div className="font-medium">Día #{d.day_index}</div>
                    {d.notes?.length ? <div className="text-xs text-gray-600">{d.notes}</div> : null}
                    <div className="text-[11px] text-gray-500 font-mono">{d.id}</div>
                  </li>
                ))}
                {((weeks.find((w: any) => w.id === selectedWeekId)?.days ?? []).length === 0) && (
                  <li className="text-sm text-gray-500">Aún sin días</li>
                )}
              </ul>
            </>
          ) : <div className="text-sm text-gray-500">Selecciona una semana</div>}
        </div>

        {/* Columna 3: Prescripciones */}
        <div className="rounded border p-4 dark:bg-neutral-900 dark:border-neutral-800">
          <div className="font-semibold mb-2">Prescripciones</div>
          {selectedDayId ? (
            <>
              <AddPrescriptionForm onSubmit={v => addPresM.mutateAsync(v)} submitting={addPresM.isPending} />
              <div className="mt-4 space-y-2">
                {(weeks.flatMap((w: any) => w.days).find((d: any) => d.id === selectedDayId)?.prescriptions ?? []).map((p: any) => (
                  <div key={p.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                    <div className="font-medium">{p.exercise_name ?? p.exercise_id}</div>
                    <div className="text-sm text-gray-700">
                      Series: {p.series} · Reps: {p.reps} · Descanso: {p.rest_sec ?? '—'}s · Pos: {p.position ?? 1}
                    </div>
                    <div className="text-[11px] text-gray-500 font-mono">{p.id}</div>
                  </div>
                ))}
                {((weeks.flatMap((w: any) => w.days).find((d: any) => d.id === selectedDayId)?.prescriptions ?? []).length === 0) && (
                  <div className="text-sm text-gray-500">Aún sin prescripciones</div>
                )}
              </div>
            </>
          ) : <div className="text-sm text-gray-500">Selecciona un día</div>}
        </div>
      </div>
    </div>
  )
}
