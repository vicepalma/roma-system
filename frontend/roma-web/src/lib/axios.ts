import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios'
import { getAccess, getRefresh, saveTokens, clearTokens } from './storage'
import { refreshTokens } from '@/services/auth'
import { Bus } from './bus'

// --- Tipado para _retry en el request config ---
declare module 'axios' {
  // eslint-disable-next-line @typescript-eslint/no-empty-interface
  export interface InternalAxiosRequestConfig<D = any> {
    _retry?: boolean
  }
}

function isUsable(token: string | null | undefined) {
  return !!token && token !== 'undefined' && token !== 'null'
}

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || 'http://localhost:8080',
  withCredentials: false,
  timeout: 15000,
  headers: {
    Accept: 'application/json',
  },
})

// --- Request: agrega Authorization si hay access token válido ---
api.interceptors.request.use((config) => {
  const token = getAccess()
  // Puedes omitir enviar Authorization en login/register si quieres:
  // const isAuthRoute = config.url?.startsWith('/auth/')
  // if (!isAuthRoute && isUsable(token)) { ... }
  if (isUsable(token)) {
    config.headers = config.headers ?? {}
    ;(config.headers as any)['Authorization'] = `Bearer ${token}`
  }
  return config
})

// --- Manejamos cola mientras se refresca ---
let isRefreshing = false
let queue: Array<{ resolve: (v: unknown) => void; reject: (r?: any) => void }> = []

function enqueue<T = unknown>() {
  return new Promise<T>((resolve, reject) => queue.push({ resolve, reject }))
}

function flushQueue(err?: any) {
  if (err) {
    queue.forEach(({ reject }) => reject(err))
  } else {
    queue.forEach(({ resolve }) => resolve(true))
  }
  queue = []
}

// --- Response: intenta refresh en 401 y reintenta original ---
api.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const original = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    // Si no hay response o no es 401, deja pasar el error
    if (!error.response || error.response.status !== 401) {
      return Promise.reject(error)
    }

    // Evita bucle si el 401 viene del refresh o si ya reintentamos
    const url = original?.url || ''
    const isRefreshCall = url.includes('/auth/refresh')
    if (isRefreshCall || original?._retry) {
      // Sesión inválida → logout global
      clearTokens()
      Bus.emit('auth:logout', { reason: 'unauthorized' })
      return Promise.reject(error)
    }

    // Si ya se está refrescando, encola y reintenta cuando termine
    if (isRefreshing) {
      try {
        await enqueue()
        return api(original) // reintenta con el nuevo token
      } catch (e) {
        return Promise.reject(e)
      }
    }

    // Marca y refresca
    original._retry = true
    isRefreshing = true
    try {
      const r = getRefresh()
      if (!isUsable(r)) throw new Error('No refresh token')

      const { access, refresh } = await refreshTokens(r as string)
      saveTokens(access, refresh)

      flushQueue() // libera a los que esperaban
      return api(original) // reintenta petición original
    } catch (e) {
      flushQueue(e)
      clearTokens()
      Bus.emit('auth:logout', { reason: 'unauthorized' })
      return Promise.reject(error)
    } finally {
      isRefreshing = false
    }
  }
)

export default api
