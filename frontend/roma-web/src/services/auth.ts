import api from '@/lib/axios'
import type { LoginRequest, LoginResponse, LoginApiResponse } from '@/types/auth'

function extractTokens(data: LoginApiResponse): LoginResponse {
  if ((data as any).tokens?.access && (data as any).tokens?.refresh) {
    return (data as any).tokens as LoginResponse
  }
  if ((data as any).access && (data as any).refresh) {
    return { access: (data as any).access, refresh: (data as any).refresh }
  }
  throw new Error('Formato de respuesta de login/refresh no reconocido')
}

export async function login(payload: LoginRequest): Promise<LoginResponse> {
  const { data } = await api.post<LoginApiResponse>('/auth/login', payload)
  return extractTokens(data)
}

export async function refreshTokens(refresh: string): Promise<LoginResponse> {
  const { data } = await api.post<LoginApiResponse>('/auth/refresh', { refresh })
  return extractTokens(data)
}
