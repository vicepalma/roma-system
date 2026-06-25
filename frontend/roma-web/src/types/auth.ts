export type LoginRequest = { email: string; password: string }
export type LoginResponse = { access: string; refresh: string }
export type UserRole = 'coach' | 'disciple'
export type AuthUser = { id: string; email: string; name: string; role: UserRole }

// Respuesta real del backend (puede venir anidada)
export type LoginApiResponse =
  | { tokens: LoginResponse; user?: any }
  | { access: string; refresh: string; user?: any }
