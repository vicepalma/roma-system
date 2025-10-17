import type { LoggedSet, SessionSet } from '@/types/sessions'
export function groupSetsByPrescription(
  sets: LoggedSet[],
  nameByPrescriptionId: (pid: string) => string
) {
  const map = new Map<string, { name: string; items: LoggedSet[] }>()
  for (const s of sets) {
    const key = s.prescription_id
    if (!map.has(key)) {
      map.set(key, { name: nameByPrescriptionId(key), items: [] })
    }
    map.get(key)!.items.push(s)
  }
  // orden por prescripción (alfa) y set_index asc
  const groups = Array.from(map.entries()).map(([id, g]) => ({
    prescription_id: id,
    name: g.name,
    items: g.items.sort((a, b) => a.set_index - b.set_index),
  }))
  groups.sort((a, b) => a.name.localeCompare(b.name))
  return groups
}

export function nextSetIndexForPrescription(existing: SessionSet[], prescriptionId: string): number {
  return existing
    .filter(s => s.prescription_id === prescriptionId)
    .reduce((m, s) => Math.max(m, s.set_index), 0) + 1
}
