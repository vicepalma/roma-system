import type { Today, Prescription, NullableInt, NullableString } from '@/types/disciples'

function fromNullableString(v: NullableString): string | undefined {
  if (v == null) return undefined
  if (typeof v === 'string') return v
  const s = v as { String?: string; Valid?: boolean }
  if (s.Valid && typeof s.String === 'string') return s.String
  return undefined
}

function fromNullableInt(v: NullableInt): number | undefined {
  if (v == null) return undefined
  if (typeof v === 'number') return v
  const n = v as { Int32?: number; Valid?: boolean }
  if (n.Valid && typeof n.Int32 === 'number') return n.Int32
  return undefined
}

export function normalizeToday(input: any): Today {
  if (!input || typeof input !== 'object') return {}

  const d = input.day || {}
  const day = {
    id: d.ID ?? d.id ?? '',
    week_id: d.WeekID ?? d.week_id ?? '',
    day_index: d.DayIndex ?? d.day_index ?? 0,
    notes: fromNullableString(d.Notes),
  }

  const prescriptions: Prescription[] = Array.isArray(input.prescriptions)
    ? input.prescriptions.map((p: any) => ({
        id: p.ID ?? p.id,
        day_id: p.DayID ?? p.day_id,
        exercise_id: p.ExerciseID ?? p.exercise_id,
        series: p.Series ?? p.series ?? 0,
        reps: p.Reps ?? p.reps ?? '',
        rest_sec: fromNullableInt(p.RestSec),
        to_failure: p.ToFailure ?? p.to_failure ?? false,
        position: p.Position ?? p.position,
        exercise_name: p.ExerciseName ?? p.exercise_name ?? '',
        primary_muscle: p.PrimaryMuscle ?? p.primary_muscle,
        equipment: fromNullableString(p.Equipment),
      }))
    : []

  return {
    assignment_id: input.assignment_id ?? input.assignmentId,
    day,
    prescriptions,
  }
}
