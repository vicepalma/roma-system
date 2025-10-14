import { useParams, NavLink } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getSessionSets } from '@/services/sessions'
import { useToast } from '@/components/toast/ToastProvider'

function fmt(n?: number | null) {
  return (n ?? '') === '' || n == null ? '—' : String(n)
}

export default function SessionView() {
  const { id = '' } = useParams()
  const { show } = useToast()

  const setsQ = useQuery({
    queryKey: ['session', id, 'sets'],
    queryFn: () => getSessionSets(id),
    enabled: !!id,
    staleTime: 10_000,
  })

  if (setsQ.isError) {
    show({ type: 'error', message: 'No se pudieron cargar los sets' })
  }

  const rows = setsQ.data ?? []

  return (
    <div className="mx-auto max-w-5xl p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Sesión: <span className="font-mono">{id}</span></h2>
        <NavLink to="/dashboard" className="text-blue-600 hover:underline text-sm">Volver</NavLink>
      </div>

      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800">
        <div className="p-4 font-semibold">Sets</div>
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-50 dark:bg-neutral-800 text-gray-600 dark:text-neutral-300">
              <tr>
                <th className="px-4 py-2 text-left">#</th>
                <th className="px-4 py-2 text-left">Prescripción</th>
                <th className="px-4 py-2 text-left">Reps</th>
                <th className="px-4 py-2 text-left">Peso</th>
                <th className="px-4 py-2 text-left">RPE</th>
                <th className="px-4 py-2 text-left">Al fallo</th>
                <th className="px-4 py-2 text-left">Creado</th>
              </tr>
            </thead>
            <tbody>
              {rows.map(s => (
                <tr key={s.id} className="border-t dark:border-neutral-800">
                  <td className="px-4 py-2 font-mono text-xs">{s.set_index}</td>
                  <td className="px-4 py-2 font-mono text-[11px]">{s.prescription_id}</td>
                  <td className="px-4 py-2">{s.reps}</td>
                  <td className="px-4 py-2">{fmt(s.weight)}</td>
                  <td className="px-4 py-2">{fmt(s.rpe)}</td>
                  <td className="px-4 py-2">{s.to_failure ? 'Sí' : 'No'}</td>
                  <td className="px-4 py-2 text-xs text-gray-500">{s.created_at ?? '—'}</td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr>
                  <td colSpan={7} className="px-4 py-6 text-center text-gray-500">Sin sets registrados</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      <div className="text-xs text-gray-500">
        Tip: desde el detalle del discípulo puedes seguir registrando sets; esta vista te sirve para auditar.
      </div>
    </div>
  )
}
