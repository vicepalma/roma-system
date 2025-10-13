import { NavLink } from 'react-router-dom'
import { clsx } from 'clsx'

export default function Sidebar() {
  const link = ({ isActive }: { isActive: boolean }) =>
    clsx(
      'block rounded-md px-3 py-2 text-sm',
      isActive
        ? 'bg-black text-white'
        : 'text-gray-700 hover:bg-gray-100',
      'dark:text-gray-200 dark:hover:bg-neutral-800'
    )

  return (
    <aside className="hidden md:block w-60 shrink-0 border-r bg-white dark:bg-neutral-900 dark:border-neutral-800">
      <div className="p-4 text-sm font-semibold">ROMA System</div>
      <nav className="px-3 space-y-1">
        <NavLink to="/dashboard" className={link}>Dashboard</NavLink>
        <NavLink to="/assignments" className={link}>Asignaciones</NavLink>
        <NavLink to="/profile" className={link}>Perfil</NavLink>
      </nav>
    </aside>
  )
}
