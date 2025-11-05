import { useNavigate } from 'react-router-dom'
import { useState, FormEvent } from 'react'
import { login } from '@/services/auth'
import useAuth from '@/store/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card } from '@/components/ui/card'
import { CardHeader } from '@/components/ui/card-header'
import { CardTitle } from '@/components/ui/card-title'
import { CardContent } from '@/components/ui/card-content'
import { Bus } from '@/lib/bus'

export default function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { setTokens } = useAuth()
  const navigate = useNavigate()

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setLoading(true); setError(null)
    try {
      const tokens = await login({ email, password })
      setTokens(tokens)
      Bus.emit('auth:login')
      navigate('/sessions', { replace: true })
    } catch (err: any) {
      const msg = err?.response?.data?.message || 'Credenciales inválidas o servidor no disponible'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-[calc(100vh-56px)] grid place-items-center px-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Inicia sesión</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input id="email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
            </div>
            {error && <div className="text-sm text-red-600 border border-red-100 bg-red-50 rounded p-2">{error}</div>}
            <Button type="submit" disabled={loading} className="w-full">{loading ? 'Ingresando...' : 'Entrar'}</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
