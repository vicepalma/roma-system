import { useEffect, useMemo, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useToast } from '@/components/toast/ToastProvider'
import {
  listMyPrograms, createProgram, type Program,
  listWeeks, addWeek, type ProgramWeek,
  listDays, addDay, deleteDay, type ProgramDay,
  listPrescriptions, addPrescription, deletePrescription, type DayPrescription,
} from '@/services/programs'
import Modal from '@/components/ui/Modal'
import PrescriptionForm from '@/components/forms/PrescriptionForm'

export default function Programs() {
  const { show } = useToast()
  const qc = useQueryClient()

  // ---- Queries base ----
  const programsQ = useQuery({ queryKey: ['programs', 'mine'], queryFn: listMyPrograms, staleTime: 60_000 })

  const [selectedProgram, setSelectedProgram] = useState<Program | null>(null)
  const [selectedWeek, setSelectedWeek] = useState<ProgramWeek | null>(null)
  const [selectedDay, setSelectedDay] = useState<ProgramDay | null>(null)

  // Al seleccionar programa, cargar semanas
  const weeksQ = useQuery({
    queryKey: ['programs', selectedProgram?.id, 'weeks'],
    queryFn: () => listWeeks(selectedProgram!.id),
    enabled: !!selectedProgram?.id,
    staleTime: 30_000,
  })

  // Al seleccionar semana, cargar días
  const daysQ = useQuery({
    queryKey: ['programs', selectedProgram?.id, 'weeks', selectedWeek?.id, 'days'],
    queryFn: () => listDays(selectedProgram!.id, selectedWeek!.id),
    enabled: !!selectedProgram?.id && !!selectedWeek?.id,
    staleTime: 30_000,
  })

  // Al seleccionar día, cargar prescripciones
  const prescQ = useQuery({
    queryKey: ['days', selectedDay?.id, 'prescriptions'],
    queryFn: () => listPrescriptions(selectedDay!.id),
    enabled: !!selectedDay?.id,
    staleTime: 15_000,
  })

  useEffect(() => {
    if (programsQ.isError) show({ type: 'error', message: 'No se pudieron cargar los programas' })
  }, [programsQ.isError, show])

  // ---- Mutations ----
  const createProgramM = useMutation({
    mutationFn: createProgram,
    onSuccess: async (p) => {
      await qc.invalidateQueries({ queryKey: ['programs', 'mine'] })
      setSelectedProgram(p)
      show({ type: 'success', message: 'Programa creado' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo crear el programa' }),
  })

  const addWeekM = useMutation({
    mutationFn: ({ programId, index, title }: { programId: string; index: number; title?: string | null }) =>
      addWeek(programId, { index, title }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['programs', selectedProgram?.id, 'weeks'] })
      show({ type: 'success', message: 'Semana agregada' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo agregar la semana' }),
  })

  const addDayM = useMutation({
    mutationFn: ({ programId, weekId, day_index, notes }: { programId: string; weekId: string; day_index: number; notes?: string | null }) =>
      addDay(programId, weekId, { day_index, notes }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['programs', selectedProgram?.id, 'weeks', selectedWeek?.id, 'days'] })
      show({ type: 'success', message: 'Día agregado' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo agregar el día' }),
  })

  const delDayM = useMutation({
    mutationFn: ({ programId, weekId, dayId }: { programId: string; weekId: string; dayId: string }) =>
      deleteDay(programId, weekId, dayId),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['programs', selectedProgram?.id, 'weeks', selectedWeek?.id, 'days'] })
      setSelectedDay(null)
      show({ type: 'success', message: 'Día eliminado' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo eliminar el día' }),
  })

  const addPrescM = useMutation({
    mutationFn: ({ dayId, values }: { dayId: string; values: Parameters<typeof addPrescription>[1] }) =>
      addPrescription(dayId, values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['days', selectedDay?.id, 'prescriptions'] })
      setOpenNewPresc(false)
      show({ type: 'success', message: 'Prescripción agregada' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo agregar la prescripción' }),
  })

  const delPrescM = useMutation({
    mutationFn: (id: string) => deletePrescription(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['days', selectedDay?.id, 'prescriptions'] })
      show({ type: 'success', message: 'Prescripción eliminada' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo eliminar la prescripción' }),
  })

  // ---- UI State ----
  const [openNewProgram, setOpenNewProgram] = useState(false)
  const [openNewWeek, setOpenNewWeek] = useState(false)
  const [openNewDay, setOpenNewDay] = useState(false)
  const [openNewPresc, setOpenNewPresc] = useState(false)

  const weeks = weeksQ.data ?? []
  const days = daysQ.data ?? []
  const presc = prescQ.data ?? []

  // Próximos índices sugeridos
  const nextWeekIndex = useMemo(() => (weeks.length ? Math.max(...weeks.map(w => w.index)) + 1 : 1), [weeks])
  const nextDayIndex = useMemo(() => (days.length ? Math.max(...days.map(d => d.day_index)) + 1 : 1), [days])
  const nextPosition = useMemo(() => (presc.length ? Math.max(...presc.map(p => p.position)) + 1 : 1), [presc])

  return (
    <div className="grid md:grid-cols-[320px,1fr] gap-4">
      {/* Columna izquierda: Programas */}
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4 space-y-3">
        <div className="flex items-center justify-between">
          <div className="font-semibold">Programas</div>
          <button
            onClick={() => setOpenNewProgram(true)}
            className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
          >
            Nuevo
          </button>
        </div>

        {programsQ.isLoading && <div className="h-24 bg-gray-100 dark:bg-neutral-800 rounded" />}
        {programsQ.isError && <div className="text-red-600 text-sm">Error al cargar programas</div>}

        <ul className="space-y-1">
          {(programsQ.data ?? []).map(p => (
            <li key={p.id}>
              <button
                onClick={() => { setSelectedProgram(p); setSelectedWeek(null); setSelectedDay(null) }}
                className={`w-full text-left rounded px-2 py-1 text-sm ${
                  selectedProgram?.id === p.id ? 'bg-black text-white' : 'hover:bg-gray-100 dark:hover:bg-neutral-800'
                }`}
              >
                {p.title}
              </button>
            </li>
          ))}
          {!programsQ.isLoading && (programsQ.data ?? []).length === 0 && (
            <li className="text-xs text-gray-500">Sin programas</li>
          )}
        </ul>
      </div>

      {/* Columna derecha: detalle */}
      <div className="space-y-4">
        {/* Semanas */}
        <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
          <div className="flex items-center justify-between">
            <div className="font-semibold">Semanas {selectedProgram ? `— ${selectedProgram.title}` : ''}</div>
            {selectedProgram && (
              <button
                onClick={() => setOpenNewWeek(true)}
                className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
              >
                Agregar semana
              </button>
            )}
          </div>

          {!selectedProgram && <div className="text-sm text-gray-500">Selecciona o crea un programa</div>}
          {selectedProgram && weeksQ.isLoading && <div className="h-16 bg-gray-100 dark:bg-neutral-800 rounded" />}
          {selectedProgram && weeksQ.isError && <div className="text-red-600 text-sm">Error al cargar semanas</div>}

          {selectedProgram && weeks.length > 0 && (
            <div className="flex flex-wrap gap-2 mt-2">
              {weeks.map(w => (
                <button
                  key={w.id}
                  onClick={() => { setSelectedWeek(w); setSelectedDay(null) }}
                  className={`rounded px-3 py-1 text-sm border dark:border-neutral-800 ${
                    selectedWeek?.id === w.id ? 'bg-black text-white' : 'bg-white hover:bg-gray-50 dark:bg-neutral-900'
                  }`}
                >
                  Semana {w.index}{w.title ? ` — ${w.title}` : ''}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Días */}
        <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
          <div className="flex items-center justify-between">
            <div className="font-semibold">Días {selectedWeek ? `— Semana ${selectedWeek.index}` : ''}</div>
            {selectedWeek && (
              <button
                onClick={() => setOpenNewDay(true)}
                className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
              >
                Agregar día
              </button>
            )}
          </div>

          {!selectedWeek && <div className="text-sm text-gray-500">Selecciona una semana</div>}
          {selectedWeek && daysQ.isLoading && <div className="h-16 bg-gray-100 dark:bg-neutral-800 rounded" />}

          {selectedWeek && days.length > 0 && (
            <ul className="mt-2 space-y-2">
              {days.map(d => (
                <li key={d.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                  <div className="flex items-center justify-between">
                    <button
                      onClick={() => setSelectedDay(d)}
                      className={`text-sm font-medium ${selectedDay?.id === d.id ? 'text-blue-600' : ''}`}
                    >
                      Día {d.day_index}
                    </button>
                    <button
                      onClick={() => delDayM.mutate({ programId: selectedProgram!.id, weekId: selectedWeek!.id, dayId: d.id })}
                      className="text-xs text-red-600 rounded px-2 py-0.5 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                    >
                      Eliminar
                    </button>
                  </div>
                  {d.notes && <div className="text-xs text-gray-600 mt-1">{d.notes}</div>}
                </li>
              ))}
            </ul>
          )}
          {selectedWeek && !daysQ.isLoading && days.length === 0 && (
            <div className="text-sm text-gray-500 mt-2">Sin días</div>
          )}
        </div>

        {/* Prescripciones del día */}
        <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
          <div className="flex items-center justify-between">
            <div className="font-semibold">Prescripciones {selectedDay ? `— Día ${selectedDay.day_index}` : ''}</div>
            {selectedDay && (
              <button
                onClick={() => setOpenNewPresc(true)}
                className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
              >
                Agregar prescripción
              </button>
            )}
          </div>

          {!selectedDay && <div className="text-sm text-gray-500">Selecciona un día</div>}
          {selectedDay && prescQ.isLoading && <div className="h-16 bg-gray-100 dark:bg-neutral-800 rounded" />}

          {selectedDay && presc.length > 0 && (
            <ul className="mt-2 space-y-2">
              {presc.map(p => (
                <li key={p.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                  <div className="flex items-center justify-between">
                    <div className="text-sm font-medium">
                      {p.exercise_name ?? p.exercise_id}
                    </div>
                    <button
                      onClick={() => delPrescM.mutate(p.id)}
                      className="text-xs text-red-600 rounded px-2 py-0.5 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                    >
                      Eliminar
                    </button>
                  </div>
                  <div className="text-xs text-gray-600 mt-1">
                    Series: {p.series} · Reps: {p.reps} · Posición: {p.position}{p.rest_sec ? ` · Descanso: ${p.rest_sec}s` : ''}
                    {p.to_failure ? ' · A fallo' : ''}
                  </div>
                </li>
              ))}
            </ul>
          )}
          {selectedDay && !prescQ.isLoading && presc.length === 0 && (
            <div className="text-sm text-gray-500 mt-2">Sin prescripciones</div>
          )}
        </div>
      </div>

      {/* Modales */}
      <Modal open={openNewProgram} onClose={() => setOpenNewProgram(false)} title="Nuevo programa">
        <form
          className="space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            const fd = new FormData(e.currentTarget as HTMLFormElement)
            createProgramM.mutate({ title: String(fd.get('title') || ''), description: String(fd.get('description') || '') || null })
          }}
        >
          <label className="text-sm block">
            <div className="mb-1">Título *</div>
            <input name="title" required className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
          </label>
          <label className="text-sm block">
            <div className="mb-1">Descripción</div>
            <textarea name="description" rows={3} className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
          </label>
          <div className="flex items-center gap-2 pt-1">
            <button type="submit" disabled={createProgramM.isPending} className="text-sm rounded px-3 py-1 border">
              {createProgramM.isPending ? 'Creando…' : 'Crear'}
            </button>
            <button type="button" onClick={() => setOpenNewProgram(false)} className="text-sm text-gray-600 hover:underline">
              Cancelar
            </button>
          </div>
        </form>
      </Modal>

      <Modal open={openNewWeek} onClose={() => setOpenNewWeek(false)} title="Agregar semana">
        <form
          className="space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            if (!selectedProgram) return
            const fd = new FormData(e.currentTarget as HTMLFormElement)
            addWeekM.mutate({
              programId: selectedProgram.id,
              index: Number(fd.get('index') || nextWeekIndex),
              title: String(fd.get('title') || '') || null,
            })
          }}
        >
          <label className="text-sm block">
            <div className="mb-1">Índice *</div>
            <input name="index" type="number" min={1} defaultValue={nextWeekIndex} required className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
          </label>
          <label className="text-sm block">
            <div className="mb-1">Título</div>
            <input name="title" className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
          </label>
          <div className="flex items-center gap-2 pt-1">
            <button type="submit" disabled={addWeekM.isPending} className="text-sm rounded px-3 py-1 border">
              {addWeekM.isPending ? 'Agregando…' : 'Agregar'}
            </button>
            <button type="button" onClick={() => setOpenNewWeek(false)} className="text-sm text-gray-600 hover:underline">
              Cancelar
            </button>
          </div>
        </form>
      </Modal>

      <Modal open={openNewDay} onClose={() => setOpenNewDay(false)} title="Agregar día">
        <form
          className="space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            if (!selectedProgram || !selectedWeek) return
            const fd = new FormData(e.currentTarget as HTMLFormElement)
            addDayM.mutate({
              programId: selectedProgram.id,
              weekId: selectedWeek.id,
              day_index: Number(fd.get('day_index') || nextDayIndex),
              notes: String(fd.get('notes') || '') || null,
            })
          }}
        >
          <label className="text-sm block">
            <div className="mb-1">Índice *</div>
            <input name="day_index" type="number" min={1} defaultValue={nextDayIndex} required className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
          </label>
          <label className="text-sm block">
            <div className="mb-1">Notas</div>
            <textarea name="notes" rows={2} className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
          </label>
          <div className="flex items-center gap-2 pt-1">
            <button type="submit" disabled={addDayM.isPending} className="text-sm rounded px-3 py-1 border">
              {addDayM.isPending ? 'Agregando…' : 'Agregar'}
            </button>
            <button type="button" onClick={() => setOpenNewDay(false)} className="text-sm text-gray-600 hover:underline">
              Cancelar
            </button>
          </div>
        </form>
      </Modal>

      <Modal open={openNewPresc} onClose={() => setOpenNewPresc(false)} title="Agregar prescripción">
        {selectedDay ? (
          <PrescriptionForm
            defaultValues={{ position: nextPosition }}
            submitting={addPrescM.isPending}
            onCancel={() => setOpenNewPresc(false)}
            onSubmit={async  (vals) => { await addPrescM.mutateAsync({ dayId: selectedDay.id, values: vals }) }}
          />
        ) : (
          <div className="text-sm text-gray-500">Selecciona un día</div>
        )}
      </Modal>
    </div>
  )
}
