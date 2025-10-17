// src/types/backend.common.ts
export type NullString = { String?: string; Valid?: boolean }
export type NullInt32  = { Int32?: number; Valid?: boolean }

// Helpers (útiles en normalización UI)
export const unboxString = (n?: NullString) =>
  n?.Valid ? (n.String ?? null) : null

export const unboxInt = (n?: NullInt32) =>
  n?.Valid ? (typeof n.Int32 === 'number' ? n.Int32 : null) : null
