import { create } from 'zustand'
type Theme = 'light' | 'dark'
const prefersDark = () => window.matchMedia?.('(prefers-color-scheme: dark)').matches
export const useTheme = create<{ theme: Theme; toggle: () => void }>((set, get) => ({
  theme: (localStorage.getItem('roma.theme') as Theme) || (prefersDark() ? 'dark' : 'light'),
  toggle: () => {
    const next = get().theme === 'dark' ? 'light' : 'dark'
    localStorage.setItem('roma.theme', next)
    set({ theme: next })
  },
}))
