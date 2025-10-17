export type CoachLink = {
  id: string
  coach_id: string
  disciple_id: string
  status: string          // 'pending' | 'accepted' | 'rejected' | ...
  created_at: string
  updated_at: string
}

// Lo que lista el backend para /api/coach/disciples
export type CoachDisciple = {
  id: string
  name: string
  email: string
}

// Si necesitas respuestas envoltorio:
export type CoachLinksResponse = {
  incoming: CoachLink[]
  outgoing: CoachLink[]
}

export type CoachDisciplesResponse = {
  items: CoachDisciple[]
}
