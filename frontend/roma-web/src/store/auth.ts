import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { AuthUser, LoginResponse } from '@/types/auth'
import { saveTokens, clearTokens } from '@/lib/storage'

interface AuthState {
  access: string | null
  refresh: string | null
  user: AuthUser | null
  isAuthenticated: boolean
  setTokens: (t: LoginResponse) => void
  setUser: (u: AuthUser | null) => void
  logout: () => void
}

const useAuth = create<AuthState>()(
  persist(
    (set) => ({
      access: null,
      refresh: null,
      user: null,
      isAuthenticated: false,
      setUser: (user) => set({ user }),
      setTokens: ({ access, refresh }) => {
        saveTokens(access, refresh)
        set({ access, refresh, isAuthenticated: true })
      },
      logout: () => {
        clearTokens()
        set({ access: null, refresh: null, user: null, isAuthenticated: false })
      },
    }),
    { name: 'roma.auth' }
  )
)

export default useAuth
