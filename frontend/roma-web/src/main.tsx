// src/main.tsx
import React, { useEffect } from 'react'
import ReactDOM from 'react-dom/client'
import { createBrowserRouter, RouterProvider, Navigate } from 'react-router-dom'
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
import { useTheme } from './store/theme'
import { queryClient } from './lib/query'
import Exercises from './pages/Exercises'
import Programs from './pages/Programs'
import ProgramDetail from './pages/ProgramDetail'

import SessionView from './pages/SessionView'
import History from './pages/History'
import SessionsIndex from './pages/SessionsIndex'
import RootRedirect from './routes/RootRedirect'

const router = createBrowserRouter([
  // Login fullscreen fuera del layout
  {
    path: '/auth/login',
    element: <Login />,
    errorElement: <ErrorBoundary />,
  },

  // Raíz decide a dónde ir según autenticación
  {
    path: '/',
    element: <RootRedirect />,
    errorElement: <ErrorBoundary />,
  },

  // Rutas con layout (App) y contenido
  {
    path: '/',
    element: <App />,
    errorElement: <ErrorBoundary />,
    children: [
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
      {
        path: '/exercises',
        element: (
          <ProtectedRoute>
            <Exercises />
          </ProtectedRoute>
        ),
      },
      {
        path: '/programs',
        element: (
          <ProtectedRoute>
            <Programs />
          </ProtectedRoute>
        ),
      },
      {
        path: '/programs/:id',
        element: (
          <ProtectedRoute>
            <ProgramDetail />
          </ProtectedRoute>
        ),
      },
      {
        path: '/sessions',
        element: (
          <ProtectedRoute>
            <SessionsIndex />
          </ProtectedRoute>
        ),
      },
      {
        path: '/sessions/:id',
        element: (
          <ProtectedRoute>
            <SessionView />
          </ProtectedRoute>
        ),
      },
      {
        path: '/history',
        element: (
          <ProtectedRoute>
            <History />
          </ProtectedRoute>
        ),
      },
      { path: '*', element: <NotFound /> },
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
