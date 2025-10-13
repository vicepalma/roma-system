import { Navigate, useLocation } from 'react-router-dom'
import useAuth from '@/store/auth'
import { ReactNode } from 'react'

export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth()
  const location = useLocation()
  if (!isAuthenticated) return <Navigate to="/auth/login" state={{ from: location }} replace />
  return <>{children}</>
}
