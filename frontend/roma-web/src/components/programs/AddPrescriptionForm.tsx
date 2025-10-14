import { useState } from 'react'

export default function AddPrescriptionForm({ onSubmit, submitting }: {
  onSubmit: (v: {
    exercise_id: string
    series: number
    reps: string
    rest_sec?: number | null
    to_failure?: boolean
    position?: number
  }) => Promise<any> | void
  submitting?: boolean
}) {
  const [exerciseId, setExerciseId] = useState('')
  const [series, setSeries] = useState(3)
  const [reps, setReps] = useState('10-12')
  const [rest, setRest] = useState<number | ''>(90)
  const [pos, setPos] = useState<number | ''>(1)
  const [toFailure, setToFailure] = useState(false)

  return (
    <form
      onSubmit={async (e) => {
        e.preventDefault()
        await onSubmit({
          exercise_id: exerciseId,
          series,
          reps,
          rest_sec: rest === '' ? null : Number(rest),
          to_failure: toFailure,
          position: pos === '' ? undefined : Number(pos),
        })
      }}
      className="grid md:grid-cols-3 gap-2"
    >
      <div className="md:col-span-3">
        <label className="block text-xs mb-1">Exercise ID</label>
        <input
          value={exerciseId}
          onChange={(e) => setExerciseId(e.target.value)}
          placeholder="Selecciona desde /exercises y pega el id (por ahora)"
          className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
        />
      </div>
      <div>
        <label className="block text-xs mb-1">Series</label>
        <input type="number" min={1} value={series} onChange={(e) => setSeries(Number(e.target.value))}
          className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800" />
      </div>
      <div>
        <label className="block text-xs mb-1">Reps</label>
        <input value={reps} onChange={(e) => setReps(e.target.value)}
          className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800" />
      </div>
      <div>
        <label className="block text-xs mb-1">Descanso (seg)</label>
        <input type="number" min={0} value={rest} onChange={(e) => setRest(e.target.value === '' ? '' : Number(e.target.value))}
          className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800" />
      </div>
      <div>
        <label className="block text-xs mb-1">Posición</label>
        <input type="number" min={1} value={pos} onChange={(e) => setPos(e.target.value === '' ? '' : Number(e.target.value))}
          className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800" />
      </div>
      <label className="inline-flex items-center gap-2">
        <input type="checkbox" checked={toFailure} onChange={(e) => setToFailure(e.target.checked)} />
        To Failure
      </label>
      <div className="md:col-span-3">
        <button disabled={submitting} className="rounded border px-3 py-2 text-sm dark:border-neutral-800">
          Agregar prescripción
        </button>
      </div>
    </form>
  )
}
