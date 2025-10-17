export type Invitation = {
  id: string
  code: string
  coach_id: string
  email: string
  name?: string | null
  status: string
  expires_at: string
  accepted_by?: string | null
  accepted_at?: string | null
  created_at: string
}