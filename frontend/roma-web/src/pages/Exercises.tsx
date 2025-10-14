import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import ExerciseFilters from '@/components/exercises/ExerciseFilters'
import { searchExercises } from '@/services/exercises'
import type { Exercise } from '@/types/exercises'
import { useToast } from '@/components/toast/ToastProvider'

export default function Exercises() {
  const { show } = useToast()
  const [query, setQuery] = useState({ q: '', tags: [] as string[], match: 'any' as 'any'|'all' })
  const [page, setPage] = useState(0)
  const limit = 20
  const offset = page * limit

  const qKey = ['exercises', query.q, query.tags.join(','), query.match, limit, offset] as const

  const listQ = useQuery({
    queryKey: qKey,
    queryFn: () => searchExercises({ ...query, limit, offset }),
    staleTime: 30_000,
  })

  const total = listQ.data?.total ?? 0
  const items: Exercise[] = listQ.data?.items ?? []
  const pages = Math.max(1, Math.ceil(total / limit))

  const onSearch = () => setPage(0) // al cambiar filtros, volver a página 0

  return (
    <div className="mx-auto max-w-6xl p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Catálogo de ejercicios</h2>
      </div>

      <ExerciseFilters value={query} onChange={setQuery} onSearch={onSearch} />

      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800">
        <div className="p-3 text-sm text-gray-600 dark:text-neutral-300">
          {listQ.isLoading ? 'Cargando…' : `Resultados: ${total}`}
          {listQ.isError && <span className="text-red-600 ml-2">Error al cargar ejercicios</span>}
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-50 dark:bg-neutral-800 text-gray-600 dark:text-neutral-300">
              <tr>
                <th className="text-left px-4 py-2">Nombre</th>
                <th className="text-left px-4 py-2">Músculo</th>
                <th className="text-left px-4 py-2">Equipo</th>
                <th className="text-left px-4 py-2">Tags</th>
              </tr>
            </thead>
            <tbody>
              {items.map((e) => (
                <tr key={e.id} className="border-t dark:border-neutral-800">
                  <td className="px-4 py-2">{e.name}</td>
                  <td className="px-4 py-2">{e.primary_muscle ?? '—'}</td>
                  <td className="px-4 py-2">{e.equipment ?? '—'}</td>
                  <td className="px-4 py-2">
                    <div className="flex flex-wrap gap-1">
                      {(e.tags ?? []).map((t) => (
                        <span key={t} className="text-[11px] rounded border px-2 py-0.5 dark:border-neutral-800">{t}</span>
                      ))}
                    </div>
                  </td>
                </tr>
              ))}
              {!listQ.isLoading && items.length === 0 && (
                <tr><td colSpan={4} className="px-4 py-6 text-center text-gray-500">Sin resultados</td></tr>
              )}
            </tbody>
          </table>
        </div>

        <div className="p-3 flex items-center justify-between">
          <div className="text-xs text-gray-600 dark:text-neutral-300">
            Página {page + 1} de {pages}
          </div>
          <div className="flex gap-2">
            <button
              className="rounded border px-2 py-1 text-sm disabled:opacity-50 dark:border-neutral-800"
              disabled={page === 0}
              onClick={() => setPage(p => Math.max(0, p - 1))}
            >Anterior</button>
            <button
              className="rounded border px-2 py-1 text-sm disabled:opacity-50 dark:border-neutral-800"
              disabled={(page + 1) >= pages}
              onClick={() => setPage(p => p + 1)}
            >Siguiente</button>
          </div>
        </div>
      </div>
    </div>
  )
}
