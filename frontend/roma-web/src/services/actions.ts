import api from '@/lib/axios'
export async function completeSet(sessionId: string, payload: { exercise_id: string; reps: number; weight?: number }) {
  const { data } = await api.post(`/api/coach/sessions/${sessionId}/complete`, payload)
  return data
}

const m = useMutation({
  mutationFn: (p: { sessionId: string; exercise_id: string; reps: number; weight?: number }) => completeSet(p.sessionId, p),
  onSuccess: () => { show({ type: 'success', message: 'Set registrado' }); queryClient.invalidateQueries({ queryKey: ['disciple', id, 'today'] })},
  onError: () => show({ type: 'error', message: 'No se pudo registrar el set' })
})
