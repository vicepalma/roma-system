export type Exercise = {
  id: string
  name: string
  primary_muscle?: string | null
  equipment?: string | null
  tags?: string[] | null
}
export type SearchExercisesResponse = {
  items: Exercise[]
  total: number
  limit: number
  offset: number
}
