// src/components/auth/AuthEvents.tsx
import { useEffect } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import useAuth from '@/store/auth'
import { useToast } from '@/components/toast/ToastProvider'
import { Bus } from '@/lib/bus'

export default function AuthEvents() {
  const navigate = useNavigate()
  const { pathname } = useLocation()
  const { isAuthenticated, logout } = useAuth()
  const { show } = useToast()

  // LOGOUT
  useEffect(() => {
    return Bus.on('auth:logout', () => {
      logout()
      show({ type: 'error', message: 'Tu sesión expiró. Vuelve a iniciar sesión.' })
      setTimeout(() => navigate('/auth/login', { replace: true }), 0)
    })
  }, [logout, navigate, show])

  // LOGIN (evento) -> navega a /sessions
  useEffect(() => {
    return Bus.on('auth:login', () => {
      show({ type: 'success', message: '¡Bienvenido! Listo para entrenar.' })
      // microtask para no competir con setTokens
      setTimeout(() => navigate('/sessions', { replace: true }), 0)
    })
  }, [navigate, show])

  // Fallback reactivo: si ya estás autenticado y estás en "/" o en "/auth/*", redirige a /sessions
  useEffect(() => {
    if (!isAuthenticated) return
    if (pathname === '/' || pathname.startsWith('/auth')) {
      setTimeout(() => navigate('/sessions', { replace: true }), 0)
    }
  }, [isAuthenticated, pathname, navigate])

  // Refresh token falló
  useEffect(() => {
    return Bus.on('auth:refresh:failed', () => {
      logout()
      show({ type: 'error', message: 'Tu sesión expiró. Inicia sesión nuevamente.' })
      setTimeout(() => navigate('/auth/login', { replace: true }), 0)
    })
  }, [logout, navigate, show])

  return null
}
