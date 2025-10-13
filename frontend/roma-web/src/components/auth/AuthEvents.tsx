import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Bus } from '@/lib/bus'
import useAuth from '@/store/auth'
import { useToast } from '@/components/toast/ToastProvider'

export default function AuthEvents() {
  const navigate = useNavigate()
  const { logout } = useAuth()
  const { show } = useToast()

  useEffect(() => {
    return Bus.on('auth:logout', () => {
      logout()
      show({ type: 'error', message: 'Tu sesión expiró. Vuelve a iniciar sesión.' })
      navigate('/auth/login')
    })
  }, [logout, navigate, show])

  return null
}
