import { useEffect, useMemo, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useToast } from '@/components/toast/ToastProvider'
import { getCoachDisciples } from '@/services/coach'
import { getCoachAssignments, createAssignment } from '@/services/assignments'
import type { AssignmentListRow } from '@/types/assignments'
import AssignmentForm from '@/components/forms/AssignmentForm'

function StatusBadge({ active }: { active: boolean }) {
  const cls = active
    ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
    : 'bg-gray-100 text-gray-700 dark:bg-neutral-800 dark:text-neutral-300'
  return (
    <span className={`inline-flex items-center rounded px-2 py-0.5 text-xs font-medium ${cls}`}>
      {active ? 'Activo' : 'Inactivo'}
    </span>
  )
}

function fmtDate(d?: string | null) {
  if (!d) return '—'
  try { return new Date(d).toLocaleDateString('es-CL') } catch { return d }
}

export default function Assignments() {
  const { show } = useToast()
  const qc = useQueryClient()
  const [showAll, setShowAll] = useState(false)

  // 1) Para el formulario: discípulos
  const disciplesQ = useQuery({
    queryKey: ['coach', 'disciples'],
    queryFn: getCoachDisciples,
    staleTime: 60_000,
  })

  // 2) Lista de asignaciones
  const assignmentsQ = useQuery({
    queryKey: ['coach', 'assignments'],
    queryFn: getCoachAssignments,
    staleTime: 30_000,
  })

  useEffect(() => {
    if (disciplesQ.isError) show({ type: 'error', message: 'Error al cargar discípulos' })
  }, [disciplesQ.isError, show])
  useEffect(() => {
    if (assignmentsQ.isError) show({ type: 'error', message: 'Error al cargar asignaciones' })
  }, [assignmentsQ.isError, show])

  const createM = useMutation({
    mutationFn: createAssignment,
    onSuccess: () => {
      show({ type: 'success', message: 'Asignación creada' })
      qc.invalidateQueries({ queryKey: ['coach', 'assignments'] })
    },
    onError: () => show({ type: 'error', message: 'No se pudo crear la asignación' }),
  })

  // Datos crudos
  const items: AssignmentListRow[] = assignmentsQ.data ?? []

  // all = lista completa; collapsed = deduplicada por (disciple_id, program_id, program_version)
  const { all, collapsed } = useMemo(() => {
    const allRows = items

    const keyOf = (a: AssignmentListRow) =>
      `${a.disciple_id}:${a.program_id}:${a.program_version}`

    const bestByKey = new Map<string, AssignmentListRow>()
    for (const a of allRows) {
      const k = keyOf(a)
      const prev = bestByKey.get(k)
      if (!prev) {
        bestByKey.set(k, a)
      } else {
        // más nueva por created_at (fallback a start_date)
        const prevTs = new Date(prev.created_at || prev.start_date).getTime()
        const currTs = new Date(a.created_at || a.start_date).getTime()
        if (currTs > prevTs) bestByKey.set(k, a)
      }
    }

    const deduped = Array.from(bestByKey.values())

    // Orden común: activas primero, luego por fecha de inicio desc
    const sortFn = (x: AssignmentListRow, y: AssignmentListRow) => {
      const rx = x.is_active ? 0 : 1
      const ry = y.is_active ? 0 : 1
      if (rx !== ry) return rx - ry
      const dx = new Date(x.start_date).getTime()
      const dy = new Date(y.start_date).getTime()
      return dy - dx
    }

    const allSorted = [...allRows].sort(sortFn)
    const dedupedSorted = [...deduped].sort(sortFn)

    return { all: allSorted, collapsed: dedupedSorted }
  }, [items])

  // Selección según toggle
  const rows = showAll ? all : collapsed

  return (
    <div className="space-y-6">
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
        <div className="font-semibold mb-2">Nueva asignación</div>
        <AssignmentForm
          disciples={disciplesQ.data ?? []}
          onSubmit={(vals) => createM.mutateAsync(vals)}
          submitting={createM.isPending}
        />
      </div>

      <div className="flex items-center justify-between px-1 sm:px-4 pt-2">
        <div className="text-sm text-gray-600 dark:text-neutral-300">
          {showAll
            ? `Mostrando ${all.length} asignaciones`
            : `Mostrando ${collapsed.length} (colapsadas de ${all.length})`}
        </div>
        <label className="text-sm inline-flex items-center gap-2 select-none">
          <input
            type="checkbox"
            className="accent-blue-600"
            checked={showAll}
            onChange={(e) => setShowAll(e.target.checked)}
          />
          Ver sin colapsar
        </label>
      </div>

      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800">
        <div className="p-4 font-semibold">Asignaciones</div>
        <div className="overflow-x-auto">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-50 dark:bg-neutral-800 text-gray-600 dark:text-neutral-300">
              <tr>
                <th className="text-left px-4 py-2">Discípulo</th>
                <th className="text-left px-4 py-2">Programa</th>
                <th className="text-left px-4 py-2">Versión</th>
                <th className="text-left px-4 py-2">Inicio</th>
                <th className="text-left px-4 py-2">Fin</th>
                <th className="text-left px-4 py-2">Estado</th>
                <th className="text-left px-4 py-2">ID</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((a) => (
                <tr key={a.id} className="border-t dark:border-neutral-800">
                  <td className="px-4 py-2">
                    <div className="font-medium">{a.disciple_name}</div>
                    <div className="text-xs text-gray-500">{a.disciple_email}</div>
                  </td>
                  <td className="px-4 py-2">
                    <div className="font-medium">{a.program_title || a.program_id}</div>
                    <div className="text-[11px] text-gray-500 font-mono">{a.program_id}</div>
                  </td>
                  <td className="px-4 py-2">{a.program_version}</td>
                  <td className="px-4 py-2">{fmtDate(a.start_date)}</td>
                  <td className="px-4 py-2">{fmtDate(a.end_date ?? null)}</td>
                  <td className="px-4 py-2"><StatusBadge active={a.is_active} /></td>
                  <td className="px-4 py-2 font-mono text-[11px]">{a.id}</td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr>
                  <td colSpan={7} className="px-4 py-6 text-center text-gray-500">
                    No hay asignaciones
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
