import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { useToast } from '@/components/toast/ToastProvider'
import useAuth from '@/store/auth'
import { listMyPrograms } from '@/services/programs'
import { createTemplate, importTemplate, listTemplates, publishTemplate } from '@/services/templates'

export default function Templates() {
  const role = useAuth(s => s.user?.role)
  const userID = useAuth(s => s.user?.id)
  const { show } = useToast(); const qc = useQueryClient()
  const [programId, setProgramId] = useState(''); const [title, setTitle] = useState(''); const [notes, setNotes] = useState('')
  const templatesQ = useQuery({ queryKey: ['templates'], queryFn: listTemplates })
  const programsQ = useQuery({ queryKey: ['programs', 'list', 'mine'], queryFn: listMyPrograms, enabled: role === 'coach' })
  const createM = useMutation({ mutationFn: createTemplate, onSuccess: async () => { setTitle(''); setNotes(''); await qc.invalidateQueries({ queryKey: ['templates'] }); show({ type: 'success', message: 'Plantilla creada como borrador' }) }, onError: () => show({ type: 'error', message: 'No se pudo crear la plantilla' }) })
  const publishM = useMutation({ mutationFn: publishTemplate, onSuccess: async () => { await qc.invalidateQueries({ queryKey: ['templates'] }); show({ type: 'success', message: 'Plantilla publicada' }) }, onError: () => show({ type: 'error', message: 'No se pudo publicar la plantilla' }) })
  const importM = useMutation({ mutationFn: importTemplate, onSuccess: async () => { await qc.invalidateQueries({ queryKey: ['programs'], exact: false }); show({ type: 'success', message: 'Rutina importada en Mis rutinas' }) }, onError: () => show({ type: 'error', message: 'No se pudo importar la plantilla' }) })
  const ownPrograms = (programsQ.data ?? []).filter(p => p.owner_id === userID && p.kind === 'coach_program')
  const templates = templatesQ.data ?? []
  return <div className="mx-auto max-w-4xl space-y-4">
    <div><h1 className="text-xl font-semibold">Biblioteca de rutinas</h1><p className="text-sm text-gray-600 dark:text-neutral-300">Importa rutinas publicadas como copias personales.</p></div>
    {role === 'coach' && <section className="rounded border bg-white p-4 dark:bg-neutral-900 dark:border-neutral-800 space-y-3"><div className="font-semibold">Publicar una plantilla</div><select value={programId} onChange={e => setProgramId(e.target.value)} className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"><option value="">Selecciona un programa propio</option>{ownPrograms.map(p => <option key={p.id} value={p.id}>{p.title}</option>)}</select><input value={title} onChange={e => setTitle(e.target.value)} placeholder="Título de plantilla" className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800" /><textarea value={notes} onChange={e => setNotes(e.target.value)} placeholder="Descripción opcional" rows={3} className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800" /><button disabled={!programId || title.trim().length < 2 || createM.isPending} onClick={() => createM.mutate({ program_id: programId, title: title.trim(), notes: notes || null })} className="rounded bg-black px-3 py-2 text-sm text-white disabled:opacity-50">Crear borrador</button></section>}
    <section className="space-y-2">{templatesQ.isLoading && <div className="rounded border p-4">Cargando biblioteca…</div>}{!templatesQ.isLoading && templates.length === 0 && <div className="rounded border border-dashed p-4 text-sm text-gray-500">No hay plantillas disponibles.</div>}{templates.map(t => <article key={t.id} className="rounded border bg-white p-4 dark:bg-neutral-900 dark:border-neutral-800"><div className="flex flex-wrap items-start justify-between gap-3"><div><h2 className="font-semibold">{t.title}</h2><p className="text-sm text-gray-600 dark:text-neutral-300">{t.notes || 'Sin descripción'}</p><span className="mt-2 inline-block rounded border px-2 py-1 text-xs">{t.status === 'published' ? 'Publicada' : 'Borrador'}</span></div><div className="flex flex-wrap gap-2">{t.status === 'published' && <button onClick={() => importM.mutate(t.id)} disabled={importM.isPending} className="rounded bg-black px-3 py-2 text-sm text-white disabled:opacity-50">Importar</button>}{role === 'coach' && t.status === 'draft' && t.owner_id === userID && <button onClick={() => publishM.mutate(t.id)} disabled={publishM.isPending} className="rounded border px-3 py-2 text-sm">Publicar</button>}</div></div></article>)}</section>
  </div>
}
