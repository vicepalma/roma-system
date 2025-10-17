import { useMemo } from 'react'
import type { SessionSet  } from '@/types/sessions'
import { groupSetsByPrescription } from '@/lib/sets'

type Props = {
  sets: SessionSet[]
  prescriptions: { id: string; exercise_name: string }[]
}

export default function TodaySets({ sets, prescriptions }: Props) {
  const nameByPid = (pid: string) => prescriptions.find(p => p.id === pid)?.exercise_name ?? pid
  const groups = useMemo(() => groupSetsByPrescription(sets, nameByPid), [sets, prescriptions])

  if (!sets.length) {
    return <div className="text-gray-500 text-sm">Aún no hay sets registrados hoy</div>
  }

  return (
    <div className="space-y-3">
      {groups.map(g => (
        <div key={g.prescription_id} className="rounded border dark:border-neutral-800">
          <div className="px-3 py-2 font-medium bg-gray-50 dark:bg-neutral-800">
            {g.name}
            <span className="ml-2 text-xs text-gray-500">({g.items.length} set{g.items.length !== 1 ? 's' : ''})</span>
          </div>
          <div className="px-3 py-2">
            <div className="grid grid-cols-4 md:grid-cols-6 text-xs text-gray-500 mb-1">
              <div>Set</div><div>Reps</div><div>Peso</div><div>RPE</div><div className="hidden md:block col-span-2">ID</div>
            </div>
            <ul className="space-y-1">
              {g.items.map(s => (
                <li key={s.id} className="grid grid-cols-4 md:grid-cols-6 text-sm">
                  <div>{s.set_index}</div>
                  <div>{s.reps}</div>
                  <div>{s.weight ?? '—'}</div>
                  <div>{s.rpe ?? '—'}</div>
                  <div className="hidden md:block col-span-2 font-mono text-[11px] text-gray-500">{s.id}</div>
                </li>
              ))}
            </ul>
          </div>
        </div>
      ))}
          <ul className="space-y-1">
      {sets.map(s => (
        <li key={s.id} className="text-sm">
          #{s.set_index} — {s.reps} reps {s.weight ? `@ ${s.weight}kg` : ''} {s.rpe ? `· RPE ${s.rpe}` : ''}{s.to_failure ? ' · a fallo' : ''}
        </li>
      ))}
    </ul>
    </div>
  )
}
