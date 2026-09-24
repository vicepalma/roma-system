import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { Button } from './components/ui/button'
import Sidebar from './components/layout/Sidebar'
import { useTheme } from './store/theme'
import AuthEvents from '@/components/auth/AuthEvents'
import useAuth from '@/store/auth'

export default function App() {
  const navigate = useNavigate()
  const location = useLocation()
  const { theme, toggle } = useTheme()
  const { logout } = useAuth()
  const [mobileNavOpen, setMobileNavOpen] = useState(false)

  useEffect(() => {
    setMobileNavOpen(false)
  }, [location.pathname])

  const handleLogout = () => {
    logout()
    setMobileNavOpen(false)
    navigate('/auth/login', { replace: true })
  }

  return (
    <div className="min-h-screen min-w-0 max-w-full overflow-x-hidden bg-gray-50 dark:bg-neutral-900 md:flex">
      <Sidebar mobileOpen={mobileNavOpen} onClose={() => setMobileNavOpen(false)} />
      <div className="flex min-w-0 flex-1 flex-col">
        <AuthEvents />
        <header className="border-b bg-white dark:border-neutral-800 dark:bg-neutral-900">
          <div className="mx-auto flex min-h-16 w-full max-w-6xl items-center justify-between gap-2 px-3 py-2 sm:px-4">
            <nav className="flex min-w-0 items-center gap-2">
              <button type="button" aria-label="Abrir navegación" onClick={() => setMobileNavOpen(true)} className="min-h-11 min-w-11 rounded-md border text-lg md:hidden dark:border-neutral-700">☰</button>
              <button type="button" className="min-h-11 rounded-md px-2 text-sm font-medium sm:px-3" onClick={() => navigate('/sessions')}>Inicio</button>
            </nav>
            <div className="flex shrink-0 items-center gap-1 sm:gap-2">
              <Button variant="outline" className="min-h-11 whitespace-nowrap px-2 text-xs sm:px-3 sm:text-sm" onClick={toggle}>
                <span className="sm:hidden">Tema</span><span className="hidden sm:inline">{theme === 'dark' ? 'Tema claro' : 'Tema oscuro'}</span>
              </Button>
              <Button variant="outline" className="min-h-11 whitespace-nowrap px-2 text-xs sm:px-3 sm:text-sm" onClick={handleLogout}>
                <span className="sm:hidden">Salir</span><span className="hidden sm:inline">Cerrar sesión</span>
              </Button>
            </div>
          </div>
        </header>
        <main className="min-w-0 flex-1">
          <div className="mx-auto min-w-0 w-full max-w-6xl overflow-x-hidden p-3 sm:p-4 md:p-6">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  )
}
