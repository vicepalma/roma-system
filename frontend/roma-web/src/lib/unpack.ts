export function unpackList<T>(payload: unknown, keys: string[]): T[] {
  if (Array.isArray(payload)) return payload as T[]
  if (payload && typeof payload === 'object') {
    for (const k of keys) {
      const v = (payload as any)[k]
      if (Array.isArray(v)) return v as T[]
    }
  }
  if (payload && typeof payload === 'object' && (payload as any).results) {
    const r = (payload as any).results
    if (Array.isArray(r)) return r as T[]
    if (r && typeof r === 'object') {
      for (const k of keys) {
        const v = r[k]
        if (Array.isArray(v)) return v as T[]
      }
    }
  }
  return [] // fallback silencioso
}
