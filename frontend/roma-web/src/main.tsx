// src/main.tsx
import React, { useEffect } from 'react'
import ReactDOM from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { QueryClientProvider } from '@tanstack/react-query'
import './index.css'

import App from './App'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import DiscipleDetail from './pages/DiscipleDetail'
import Assignments from './pages/Assignments'
import ProtectedRoute from './router/ProtectedRoute'
import ErrorBoundary from './routes/ErrorBoundary'
import NotFound from './pages/NotFound'

import { ToastProvider } from '@/components/toast/ToastProvider'
import AuthEvents from '@/components/auth/AuthEvents'
import { useTheme } from './store/theme'
import { queryClient } from './lib/query' // <-- defaults centralizados

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    errorElement: <ErrorBoundary />,
    children: [
      { index: true, element: <Login /> },
      { path: '/auth/login', element: <Login /> },
      {
        path: '/dashboard',
        element: (
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        ),
      },
      {
        path: '/disciples/:id',
        element: (
          <ProtectedRoute>
            <DiscipleDetail />
          </ProtectedRoute>
        ),
      },
      {
        path: '/assignments',
        element: (
          <ProtectedRoute>
            <Assignments />
          </ProtectedRoute>
        ),
      },
      { path: '*', element: <NotFound /> }, // <-- al final por claridad
    ],
  },
])

function ThemeEffect() {
  const { theme } = useTheme()
  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark')
  }, [theme])
  return null
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <ToastProvider>
        <ThemeEffect />
        <RouterProvider router={router} />
      </ToastProvider>
    </QueryClientProvider>
  </React.StrictMode>
)
