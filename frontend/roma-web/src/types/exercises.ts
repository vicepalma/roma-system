export type Exercise = {
  id: string
  name: string
  primary_muscle: string
  equipment?: string | null
  tags?: string[]
  notes?: string | null
  created_at?: string
  updated_at?: string
}

export type ExerciseCatalogItem = {
  id: string
  name: string
}

export type SearchExercisesResponse = {
  items: Exercise[]
  total: number
  limit: number
  offset: number
}
