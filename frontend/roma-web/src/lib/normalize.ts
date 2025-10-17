import { unboxInt, unboxString } from '@/types/backend.common'
import type { MeTodayDay, MeTodayPrescription, Today } from '@/types/disciples'

export function normalizeToday(input: any): Today {
  // si viene mal, devuelve shape vacío mínimo válido
  if (!input || typeof input !== 'object') {
    return {
      assignment_id: '', // evita TS error
      day: undefined,
      prescriptions: [],
    }
  }

  const day = input.day as MeTodayDay | undefined
  const rawPresc = (input.prescriptions ?? []) as MeTodayPrescription[]

  return {
    assignment_id: String(input.assignment_id ?? ''),
    current_session_id: input.current_session_id ?? null,
    current_session_sets_count: Number(input.current_session_sets_count ?? 0),
    current_session_started_at: input.current_session_started_at ?? null,
    day: day
      ? {
          id: day.id,
          week_id: day.week_id,
          day_index: day.day_index,
          notes: unboxString(day.notes),
        }
      : undefined,
    prescriptions: rawPresc.map(p => ({
      id: p.id,
      day_id: p.day_id,
      exercise_id: p.exercise_id,
      series: p.series,
      reps: p.reps,
      rest_sec: unboxInt(p.rest_sec),
      to_failure: p.to_failure,
      position: p.position,
      exercise_name: p.exercise_name,
      primary_muscle: p.primary_muscle,
      equipment: unboxString(p.equipment),
    })),
  }
}
