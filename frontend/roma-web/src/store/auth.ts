import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { LoginResponse } from '@/types/auth'
import { saveTokens, clearTokens } from '@/lib/storage'

interface AuthState {
  access: string | null
  refresh: string | null
  isAuthenticated: boolean
  setTokens: (t: LoginResponse) => void
  logout: () => void
}

const useAuth = create<AuthState>()(
  persist(
    (set) => ({
      access: null,
      refresh: null,
      isAuthenticated: false,
      setTokens: ({ access, refresh }) => {
        saveTokens(access, refresh)
        set({ access, refresh, isAuthenticated: true })
      },
      logout: () => {
        clearTokens()
        set({ access: null, refresh: null, isAuthenticated: false })
      },
    }),
    { name: 'roma.auth' }
  )
)

export default useAuth
