// src/App.tsx
import { Outlet, useNavigate, NavLink } from 'react-router-dom'
import useAuth from './store/auth'
import { Button } from './components/ui/button'
import Sidebar from './components/layout/Sidebar'
import { useTheme } from './store/theme'
import AuthEvents from '@/components/auth/AuthEvents'

export default function App() {
  const { isAuthenticated, logout } = useAuth()
  const navigate = useNavigate()
  const { theme, toggle } = useTheme()

  return (
    <div className="min-h-screen flex bg-gray-50 dark:bg-neutral-900">
      <Sidebar />
      <div className="flex-1 flex flex-col">

        <AuthEvents />

        <header className="border-b bg-white dark:bg-neutral-900 dark:border-neutral-800">
          <div className="mx-auto max-w-6xl px-4 py-3 flex items-center justify-between">
            <nav className="flex items-center gap-4">
              <NavLink to="/dashboard" className="text-sm font-medium">Inicio</NavLink>
            </nav>
            <div className="flex items-center gap-2">
              <Button variant="outline" onClick={toggle}>
                {theme === 'dark' ? 'Tema claro' : 'Tema oscuro'}
              </Button>
              {isAuthenticated ? (
                <Button variant="outline" onClick={() => { logout(); navigate('/auth/login') }}>
                  Cerrar sesión
                </Button>
              ) : (
                <Button onClick={() => navigate('/auth/login')}>Iniciar sesión</Button>
              )}
            </div>
          </div>
        </header>

        <main className="flex-1">
          <div className="mx-auto max-w-6xl p-6">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  )
}
