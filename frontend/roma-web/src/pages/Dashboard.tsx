import { useEffect, useMemo, useRef } from 'react'
import { useToast } from '@/components/toast/ToastProvider'
import { useQuery } from '@tanstack/react-query'
import { getCoachDisciples, getCoachLinks } from '@/services/coach'
import type { CoachDisciple, CoachLink } from '@/types/coach'
import { NavLink } from 'react-router-dom'

function StatusBadge({ status }: { status?: string }) {
  const s = (status ?? 'unknown').toLowerCase()
  const cls =
    s === 'accepted' ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300' :
    s === 'pending'  ? 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300' :
    s === 'rejected' ? 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300' :
                       'bg-gray-100 text-gray-700 dark:bg-neutral-800 dark:text-neutral-300'
  const label =
    s === 'accepted' ? 'Activo' :
    s === 'pending'  ? 'Pendiente' :
    s === 'rejected' ? 'Rechazado' : 'Desconocido'
  return <span className={`inline-flex items-center rounded px-2 py-0.5 text-xs font-medium ${cls}`}>{label}</span>
}

function formatWhen(s?: string) {
  if (!s) return '—'
  try { return new Date(s).toLocaleString('es-CL') } catch { return s }
}

export default function Dashboard() {
  const { show } = useToast()

  const disciplesQ = useQuery({
    queryKey: ['coach', 'disciples'],
    queryFn: getCoachDisciples,
    staleTime: 60_000,
  })

  const linksQ = useQuery({
    queryKey: ['coach', 'links'],
    queryFn: getCoachLinks,
    staleTime: 60_000,
  })

  // Evitar loops de toast
  const prevDiscStatus = useRef<string | null>(null)
  useEffect(() => {
    if (disciplesQ.status === 'error' && prevDiscStatus.current !== 'error') {
      show({ type: 'error', message: 'Error al cargar discípulos' })
    }
    prevDiscStatus.current = disciplesQ.status
  }, [disciplesQ.status, show])

  const prevLinksStatus = useRef<string | null>(null)
  useEffect(() => {
    if (linksQ.status === 'error' && prevLinksStatus.current !== 'error') {
      show({ type: 'error', message: 'No se pudieron cargar los vínculos' })
    }
    prevLinksStatus.current = linksQ.status
  }, [linksQ.status, show])

  // Index de links por disciple_id
  const linkByDiscipleId = useMemo(() => {
    const map = new Map<string, CoachLink>()
    for (const l of (linksQ.data ?? [])) map.set(l.disciple_id, l)
    return map
  }, [linksQ.data])

  // Unir y ORDENAR (activos, luego pendientes, luego rechazados, luego desconocidos)
  const rows = useMemo(() => {
    const base: CoachDisciple[] = disciplesQ.data ?? []
    const arr = base.map(d => {
      const link = linkByDiscipleId.get(d.id)
      return { ...d, linkStatus: link?.status, linkedAt: link?.created_at }
    })
    const rank = (s?: string) => s === 'accepted' ? 0 : s === 'pending' ? 1 : s === 'rejected' ? 2 : 3
    arr.sort((a, b) => rank(a.linkStatus) - rank(b.linkStatus))
    return arr
  }, [disciplesQ.data, linkByDiscipleId])

  if (disciplesQ.isLoading) {
    return (
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4 space-y-3">
        <div className="h-5 w-40 bg-gray-200 dark:bg-neutral-800 rounded" />
        <div className="h-32 w-full bg-gray-100 dark:bg-neutral-800 rounded" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="overflow-x-auto rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800">
        <table className="min-w-full text-sm">
          <thead className="bg-gray-50 dark:bg-neutral-800 text-gray-600 dark:text-neutral-300">
            <tr>
              <th className="text-left px-4 py-2">Nombre</th>
              <th className="text-left px-4 py-2">Email</th>
              <th className="text-left px-4 py-2">Estado</th>
              <th className="text-left px-4 py-2">Vinculado</th>
              <th className="text-left px-4 py-2">ID</th>
              <th className="px-4 py-2"></th>
            </tr>
          </thead>
          <tbody>
            {rows.map((d) => (
              <tr key={d.id} className="border-t dark:border-neutral-800">
                <td className="px-4 py-2">{d.name}</td>
                <td className="px-4 py-2">{d.email}</td>
                <td className="px-4 py-2"><StatusBadge status={d.linkStatus} /></td>
                <td className="px-4 py-2">{formatWhen(d.linkedAt)}</td>
                <td className="px-4 py-2 font-mono text-[11px]">{d.id}</td>
                <td className="px-4 py-2 text-right">
                  <NavLink
                    to={`/disciples/${d.id}`}
                    className="text-blue-600 hover:underline"
                    state={{ name: d.name, email: d.email }}
                  >
                    Ver detalle
                  </NavLink>
                </td>
              </tr>
            ))}
            {rows.length === 0 && (
              <tr>
                <td colSpan={6} className="px-4 py-6 text-center text-gray-500">
                  Sin discípulos
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
