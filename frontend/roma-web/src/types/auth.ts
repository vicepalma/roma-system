export type LoginRequest = { email: string; password: string }
export type LoginResponse = { access: string; refresh: string }

// Respuesta real del backend (puede venir anidada)
export type LoginApiResponse =
  | { tokens: LoginResponse; user?: any }
  | { access: string; refresh: string; user?: any }
