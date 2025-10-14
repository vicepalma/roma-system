import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listPrograms, createProgram } from '@/services/programs'
import { useState } from 'react'
import { NavLink } from 'react-router-dom'

export default function Programs() {
  const qc = useQueryClient()
  const listQ = useQuery({ queryKey: ['programs'], queryFn: listPrograms, staleTime: 30_000 })
  const [title, setTitle] = useState('')

  const createM = useMutation({
    mutationFn: createProgram,
    onSuccess: () => { setTitle(''); qc.invalidateQueries({ queryKey: ['programs'] }) }
  })

  const rows = listQ.data ?? []
  return (
    <div className="mx-auto max-w-5xl p-6 space-y-4">
      <h2 className="text-xl font-semibold">Programas</h2>

      <div className="rounded border p-3 dark:bg-neutral-900 dark:border-neutral-800">
        <div className="text-sm font-medium mb-2">Nuevo programa</div>
        <div className="flex gap-2">
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Título"
            className="rounded border px-3 py-2 flex-1 dark:bg-neutral-900 dark:border-neutral-800"
          />
          <button
            onClick={() => title && createM.mutateAsync({ title })}
            className="rounded border px-3 py-2 dark:border-neutral-800"
          >Crear</button>
        </div>
      </div>

      <div className="rounded border dark:bg-neutral-900 dark:border-neutral-800 overflow-x-auto">
        <table className="min-w-full text-sm">
          <thead className="bg-gray-50 dark:bg-neutral-800 text-gray-600 dark:text-neutral-300">
            <tr>
              <th className="text-left px-4 py-2">Título</th>
              <th className="text-left px-4 py-2">Versión</th>
              <th className="px-4 py-2"></th>
            </tr>
          </thead>
          <tbody>
            {rows.map((p: any) => (
              <tr key={p.id} className="border-t dark:border-neutral-800">
                <td className="px-4 py-2">{p.title}</td>
                <td className="px-4 py-2">{p.version}</td>
                <td className="px-4 py-2 text-right">
                  <NavLink to={`/programs/${p.id}`} className="text-blue-600 hover:underline">
                    Abrir
                  </NavLink>
                </td>
              </tr>
            ))}
            {rows.length === 0 && (
              <tr><td colSpan={3} className="px-4 py-6 text-center text-gray-500">No hay programas</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
