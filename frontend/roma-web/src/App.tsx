// src/App.tsx
import { Outlet, useNavigate } from 'react-router-dom'
import { Button } from './components/ui/button'
import Sidebar from './components/layout/Sidebar'
import { useTheme } from './store/theme'
import AuthEvents from '@/components/auth/AuthEvents'
import useAuth from '@/store/auth'

export default function App() {
  const navigate = useNavigate()
  const { theme, toggle } = useTheme()
  const { logout } = useAuth()

  const handleLogout = () => {
    // limpia estado/tokens y navega al login
    logout()
    navigate('/auth/login', { replace: true })
  }

  return (
    <div className="min-h-screen flex bg-gray-50 dark:bg-neutral-900">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <AuthEvents />
        <header className="border-b bg-white dark:bg-neutral-900 dark:border-neutral-800">
          <div className="mx-auto max-w-6xl px-4 py-3 flex items-center justify-between">
            <nav className="flex items-center gap-4">
              <button
                className="text-sm font-medium"
                onClick={() => navigate('/sessions')}
              >
                Inicio
              </button>
            </nav>
            <div className="flex items-center gap-2">
              <Button variant="outline" onClick={toggle}>
                {theme === 'dark' ? 'Tema claro' : 'Tema oscuro'}
              </Button>
              <Button variant="outline" onClick={handleLogout}>
                Cerrar sesión
              </Button>
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
