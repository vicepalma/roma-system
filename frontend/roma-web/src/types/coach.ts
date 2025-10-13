// Lista plana de discípulos (shape actual de GET /api/coach/disciples)
export type CoachDisciple = {
  id: string
  name: string
  email: string
}

// Respuesta para la lista plana
export type CoachDisciplesListResponse = {
  items: CoachDisciple[]
}

// Vínculos coach–discípulo (si el backend expone links incoming/outgoing)
export type CoachLink = {
  id: string
  coach_id: string
  disciple_id: string
  status: 'pending' | 'accepted' | 'rejected' | string
  created_at: string
  updated_at: string
}

// Respuesta para links (incoming/outgoing)
export type CoachLinksResponse = {
  incoming: CoachLink[]
  outgoing: CoachLink[]
}
