import { NavLink } from 'react-router-dom'
import { clsx } from 'clsx'
import { useSessionHydrated } from '@/store/session'
import useAuth from '@/store/auth'

type SidebarProps = {
  mobileOpen?: boolean
  onClose?: () => void
}

export default function Sidebar({ mobileOpen = false, onClose }: SidebarProps) {
  const hydrated = useSessionHydrated()
  const role = useAuth(s => s.user?.role)

  const link = ({ isActive }: { isActive: boolean }) =>
    clsx(
      'block min-h-11 rounded-md px-3 py-3 text-sm',
      isActive ? 'bg-black text-white' : 'text-gray-700 hover:bg-gray-100',
      'dark:text-gray-200 dark:hover:bg-neutral-800'
    )

  const links = (close?: () => void) => (
    <>
      {role === 'coach' && (
        <>
          <NavLink to="/dashboard" className={link} onClick={close}>Dashboard</NavLink>
          <NavLink to="/assignments" className={link} onClick={close}>Asignaciones</NavLink>
        </>
      )}
      {(role === 'disciple' || role === 'coach') && (
        <>
          {role === 'coach' && <div className="px-3 pt-3 pb-1 text-xs font-medium text-gray-500 dark:text-neutral-400">Personal</div>}
          <NavLink to="/sessions" className={link} onClick={close}>Entrenar</NavLink>
          <NavLink to="/programs" className={link} onClick={close}>Mis rutinas</NavLink>
          <NavLink to="/templates" className={link} onClick={close}>Biblioteca</NavLink>
          <NavLink to="/exercises" className={link} onClick={close}>Ejercicios</NavLink>
          <NavLink to="/history" className={link} onClick={close}>Historial</NavLink>
          <NavLink to="/checkins" className={link} onClick={close}>Check-ins</NavLink>
        </>
      )}
      {role !== 'disciple' && role !== 'coach' && (
        <>
          <NavLink to="/exercises" className={link} onClick={close}>Ejercicios</NavLink>
          <NavLink to="/history" className={link} onClick={close}>Historial</NavLink>
        </>
      )}
    </>
  )

  const loading = !hydrated
  return (
    <>
      <aside className="hidden w-60 shrink-0 border-r bg-white dark:border-neutral-800 dark:bg-neutral-900 md:block">
        <div className="p-4 text-sm font-semibold">ROMA System</div>
        <nav className="space-y-1 px-3">
          {loading ? (
            <>
              <div className="h-7 rounded bg-gray-100 dark:bg-neutral-800" />
              <div className="h-7 rounded bg-gray-100 dark:bg-neutral-800" />
              <div className="h-7 rounded bg-gray-100 dark:bg-neutral-800" />
            </>
          ) : links()}
        </nav>
      </aside>

      {mobileOpen && (
        <div className="md:hidden">
          <button type="button" aria-label="Cerrar navegación" className="fixed inset-0 z-40 bg-black/40" onClick={onClose} />
          <aside className="fixed inset-y-0 left-0 z-50 w-[min(18rem,calc(100vw-2rem))] overflow-y-auto border-r bg-white shadow-xl dark:border-neutral-800 dark:bg-neutral-900">
            <div className="flex min-h-16 items-center justify-between border-b px-4 dark:border-neutral-800">
              <div className="text-sm font-semibold">ROMA System</div>
              <button type="button" aria-label="Cerrar navegación" onClick={onClose} className="min-h-11 min-w-11 rounded-md border text-lg dark:border-neutral-700">×</button>
            </div>
            <nav className="space-y-1 p-3">{loading ? null : links(onClose)}</nav>
          </aside>
        </div>
      )}
    </>
  )
}
