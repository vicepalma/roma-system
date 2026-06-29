import { NavLink } from 'react-router-dom'
import { clsx } from 'clsx'
import { useSessionStore, useSessionHydrated } from '@/store/session'
import useAuth from '@/store/auth'

export default function Sidebar() {
  const hydrated = useSessionHydrated()
  const currentSessionId = useSessionStore(s => s.currentSessionId)
  const currentDiscipleName = useSessionStore(s => s.currentDiscipleName)
  const role = useAuth(s => s.user?.role)

  const hasSession = typeof currentSessionId === 'string' && currentSessionId.length > 0

  const link = ({ isActive }: { isActive: boolean }) =>
    clsx(
      'block rounded-md px-3 py-2 text-sm',
      isActive ? 'bg-black text-white' : 'text-gray-700 hover:bg-gray-100',
      'dark:text-gray-200 dark:hover:bg-neutral-800'
    )

  if (!hydrated) {
    return (
      <aside className="hidden md:block w-60 shrink-0 border-r bg-white dark:bg-neutral-900 dark:border-neutral-800">
        <div className="p-4 text-sm font-semibold">ROMA System</div>
        <nav className="px-3 space-y-1">
          <div className="h-7 bg-gray-100 dark:bg-neutral-800 rounded" />
          <div className="h-7 bg-gray-100 dark:bg-neutral-800 rounded" />
          <div className="h-7 bg-gray-100 dark:bg-neutral-800 rounded" />
        </nav>
      </aside>
    )
  }

  return (
    <aside className="hidden md:block w-60 shrink-0 border-r bg-white dark:bg-neutral-900 dark:border-neutral-800">
      <div className="p-4 text-sm font-semibold">ROMA System</div>
      <nav className="px-3 space-y-1">
        {role === 'coach' && <NavLink to="/dashboard" className={link}>Dashboard</NavLink>}
        {role === 'coach' && <NavLink to="/assignments" className={link}>Asignaciones</NavLink>}
        <NavLink to="/exercises" className={link}>Ejercicios</NavLink>
        {(role === 'coach' || role === 'disciple') && (
          <NavLink to="/programs" className={link}>{role === 'disciple' ? 'Mis rutinas' : 'Programas'}</NavLink>
        )}
        <NavLink to="/history" className={link}>Historial</NavLink>

        <NavLink to={hasSession ? `/sessions/${currentSessionId}` : '/sessions'} className={link}>
          {hasSession ? `Continuar sesión${currentDiscipleName ? ` — ${currentDiscipleName}` : ''}` : 'Entrenar'}
        </NavLink>
      </nav>
    </aside>
  )
}
