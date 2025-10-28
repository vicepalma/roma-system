import { useEffect, useMemo, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useToast } from '@/components/toast/ToastProvider'
import {
  listMyPrograms, createProgram, type Program,
  listWeeks, addWeek, type ProgramWeek,
  listDays, addDay, deleteDay, type ProgramDay,
  listPrescriptions, addPrescription, deletePrescription, deleteProgram, deleteWeek, type DayPrescription,
} from '@/services/programs'
import Modal from '@/components/ui/Modal'
import PrescriptionForm from '@/components/forms/PrescriptionForm'
import clsx from 'clsx'
import { listExercises } from '@/services/exercises'
import type { Exercise } from '@/types/exercises'

export default function Programs() {
  const { show } = useToast()
  const qc = useQueryClient()
  const [selectedProgram, setSelectedProgram] = useState<Program | null>(null)
  const [selectedWeek, setSelectedWeek] = useState<ProgramWeek | null>(null)
  const [selectedDay, setSelectedDay] = useState<ProgramDay | null>(null)
  const canAddDay = !!selectedProgram?.id && !!selectedWeek?.id;
  const [selectedExerciseId, setSelectedExerciseId] = useState<string>('')

  // ---- Queries base ----
  const programsQ = useQuery({
    queryKey: ['programs', 'list', 'mine'],
    queryFn: listMyPrograms,
  })

  // Al seleccionar programa, cargar semanas
  const weeksQ = useQuery({
    queryKey: ['programs', selectedProgram?.id, 'weeks'],
    queryFn: () => listWeeks(selectedProgram!.id),
    enabled: !!selectedProgram,
  })


  // Al seleccionar semana, cargar días
  const daysQ = useQuery({
    queryKey: ['programs', selectedProgram?.id, 'weeks', selectedWeek?.id, 'days'],
    queryFn: () => listDays(selectedProgram!.id, selectedWeek!.id),
    enabled: !!selectedProgram?.id && !!selectedWeek?.id,
    staleTime: 0,
  })

  // Al seleccionar día, cargar prescripciones
  const prescQ = useQuery({
    queryKey: ['days', selectedDay?.id, 'prescriptions'],
    queryFn: () => listPrescriptions(selectedDay!.id),
    enabled: !!selectedDay?.id,
    staleTime: 15_000,
  })

  const exQ = useQuery({
    queryKey: ['exercises', 'list'],
    queryFn: listExercises,
    staleTime: 60_000,
  })
  const exMap = useMemo(() => {
  const m = new Map<string, Exercise>()
  ;(exQ.data ?? []).forEach((e: Exercise) => m.set(e.id, e))
  return m
  }, [exQ.data])

  const exercises: Exercise[] = exQ.data ?? []

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
      setOpenNewProgram(false)
      await qc.invalidateQueries({ queryKey: ['programs'], exact: false, refetchType: 'active' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo crear el programa' }),
  })

  const delProgramM = useMutation({
    mutationFn: (id: string) => deleteProgram(id),

    // 1) Optimista
    onMutate: async (id: string) => {
      await qc.cancelQueries({ queryKey: ['programs'] })

      const prevList = qc.getQueryData<Program[]>(['programs', 'list']) ?? []

      qc.setQueryData<Program[]>(['programs', 'list'], prev =>
        (prev ?? []).filter(p => p.id !== id)
      )

      // Si tienes cache por detalle: ['programs', id]
      qc.removeQueries({ queryKey: ['programs', id], exact: true })
      if (selectedProgram?.id === id) setSelectedProgram(null)
      return { prevList }
    },

    // 2) Rollback si falla
    onError: (_err, _id, ctx) => {
      if (ctx?.prevList) qc.setQueryData(['programs', 'list'], ctx.prevList)
    },

    // 3) Refetch confiable de TODO lo que empiece por 'programs'
    onSettled: async () => {
      await qc.invalidateQueries({ queryKey: ['programs'], exact: false, refetchType: 'active' })
    },
  })

  const addWeekM = useMutation({
    mutationFn: (p: { programId: string; week_index: number; title?: string | null }) =>
      addWeek(p.programId, { week_index: p.week_index, title: p.title ?? null }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['programs', selectedProgram?.id, 'weeks'] })
      setOpenNewWeek(false)
      show({ type: 'success', message: 'Semana agregada' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo agregar la semana' }),
  })

  const delWeekM = useMutation({
    mutationFn: (p: { programId: string; weekId: string }) =>
      deleteWeek(p.programId, p.weekId),
    onSuccess: async (_data, vars) => {
      if (selectedWeek?.id === vars.weekId) setSelectedWeek(null)
      await qc.invalidateQueries({ queryKey: ['programs', vars.programId, 'weeks'] })
      show({ type: 'success', message: 'Semana eliminada' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo eliminar la semana' }),
  })


  const addDayM = useMutation({
    mutationFn: (v: { programId: string; weekId: string; day_index: number; notes?: string | null }) =>
      addDay(v.programId, v.weekId, { day_index: v.day_index, notes: v.notes }),
    onSuccess: async (_created, v) => {
      await qc.invalidateQueries({ queryKey: ['programs', v.programId, 'weeks', v.weekId, 'days'] })
      show({ type: 'success', message: 'Día agregado' })
      setOpenNewDay(false)
    },
    onError: (err: any) => {
      const msg = err?.response?.data?.detail || err?.response?.data?.error || 'Error'
      show({ type: 'error', message: `No se pudo agregar el día: ${msg}` })
    },
  })

  const delDayM = useMutation({
    mutationFn: ({ programId, weekId, dayId }: { programId: string; weekId: string; dayId: string }) =>
      deleteDay(programId, weekId, dayId),
    onSuccess: async () => {
      show({ type: 'success', message: 'Día eliminado' })
      await qc.invalidateQueries({ queryKey: ['programs', selectedProgram?.id, 'weeks', selectedWeek?.id, 'days'] })
    },
    onError: () => show({ type: 'error', message: 'No se pudo eliminar el día' }),
  })


  const addPrescM = useMutation({
    mutationFn: ({ dayId, values }: { dayId: string; values: Parameters<typeof addPrescription>[1] }) =>
      addPrescription(dayId, { ...values, method_id: null }),
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
  const nextWeekIndex = useMemo(() => (weeks.length ? Math.max(...weeks.map(w => w.week_index)) + 1 : 1), [weeks])
  const nextDayIndex = useMemo(() => (days.length ? Math.max(...days.map(d => d.day_index)) + 1 : 1), [days])
  const nextPosition = useMemo(() => (presc.length ? Math.max(...presc.map(p => p.position)) + 1 : 1), [presc])

  useEffect(() => {
    if (weeks?.length && !selectedWeek) setSelectedWeek(weeks[0])
  }, [weeks, selectedWeek])

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
              <div className="flex items-center justify-between hover:bg-gray-300 rounded-lg">
                <button
                  onClick={() => {
                    setSelectedProgram(p)
                    setSelectedWeek(null)
                    setSelectedDay(null)
                  }}
                  className={`text-left rounded px-2 py-1 text-sm ${selectedProgram?.id === p.id
                    ? 'bg-black text-white'
                    : 'hover:bg-gray-100 dark:hover:bg-neutral-800'
                    }`}
                >
                  {p.title}
                </button>

                <button
                  onClick={() => {
                    if (confirm('¿Eliminar este programa? Esta acción no se puede deshacer.')) {
                      delProgramM.mutate(p.id)
                    }
                  }}
                  disabled={delProgramM.isPending}
                  className="text-xs rounded px-2 py-1 border text-red-600 bg-white 
hover:bg-red-50 active:bg-red-100 
dark:bg-neutral-900 dark:border-neutral-800 
transition-colors duration-150"
                >
                  Eliminar
                </button>
              </div>
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
              {weeks.map((w) => (
                <div key={w.id} className="inline-flex items-center">
                  <button
                    onClick={() => { setSelectedWeek(w); setSelectedDay(null) }}
                    className={clsx(
                      'rounded px-3 py-1 text-sm border dark:border-neutral-800',
                      selectedWeek?.id === w.id
                        ? 'bg-black text-white'
                        : 'bg-white hover:bg-gray-50 dark:bg-neutral-900'
                    )}
                  >
                    Semana {w.week_index}{w.title ? ` — ${w.title}` : ''}
                  </button>

                  <button
                    onClick={() => {
                      console.log(w)
                      delWeekM.mutate({ programId: selectedProgram.id, weekId: w.id })
                    }}
                    disabled={delWeekM.isPending}
                    className={clsx(
                      'text-xs rounded px-2 py-1 border text-red-600 bg-white ',
                      'hover:bg-red-50 active:bg-red-100',
                      'dark:bg-neutral-900 dark:border-neutral-800',
                      'transition-colors duration-150'
                    )}
                    title="Eliminar semana"
                  >
                    {delWeekM.isPending ? '…' : 'Eliminar'}
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Días */}
        <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
          <div className="flex items-center justify-between">
            <div className="font-semibold">
              Días {selectedWeek ? `— Semana ${selectedWeek.week_index}` : ''}
            </div>
            {selectedWeek && (
              <button
                onClick={() => setOpenNewDay(true)}
                className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
              >
                Agregar día
              </button>
            )}
          </div>

          {!selectedWeek && (
            <div className="text-sm text-gray-500">Selecciona una semana</div>
          )}

          {selectedWeek && daysQ.isLoading && (
            <div className="h-16 bg-gray-100 dark:bg-neutral-800 rounded" />
          )}

          {selectedWeek && !daysQ.isLoading && days.length > 0 && (
            <ul className="mt-2 space-y-2">
              {days.map((d) => {
                const isActive = selectedDay?.id === d.id
                return (
                  <li key={d.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                    <div className="flex items-center justify-between">
                      <button
                        onClick={() => setSelectedDay(d)}
                        className={[
                          'text-sm font-medium rounded px-2 py-1 border dark:border-neutral-800',
                          isActive
                            ? 'bg-black text-white'
                            : 'bg-white hover:bg-gray-50 dark:bg-neutral-900',
                        ].join(' ')}
                      >
                        Día {d.day_index}
                      </button>

                      <button
                        onClick={() => delDayM.mutate({
                          programId: selectedProgram!.id,
                          weekId: selectedWeek!.id,
                          dayId: d.id,
                        })}
                        className="text-xs text-red-600 rounded px-2 py-0.5 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                      >
                        Eliminar
                      </button>
                    </div>

                    {d.notes && (
                      <div className="text-xs text-gray-600 mt-1">{d.notes}</div>
                    )}
                  </li>
                )
              })}
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
{presc.map(p => {
  const ex = exMap.get(p.exercise_id)
  const title = p.exercise_name && p.exercise_name.trim() !== ''
    ? p.exercise_name
    : ex?.name ?? p.exercise_id

  const chipMuscle = p.primary_muscle ?? ex?.primary_muscle
  const chipEquip = (p as any).equipment ?? ex?.equipment // por si backend no lo trae en este endpoint

  return (
    <li key={p.id} className="rounded border px-3 py-2 dark:border-neutral-800">
      <div className="flex items-center justify-between">
        <div>
          <div className="text-sm font-medium">{title}</div>

          <div className="mt-1 flex flex-wrap items-center gap-2 text-xs text-gray-600 dark:text-neutral-300">
            {chipMuscle && (
              <span className="inline-flex items-center rounded border px-1.5 py-0.5 dark:border-neutral-800">
                {chipMuscle}
              </span>
            )}
            {chipEquip && (
              <span className="inline-flex items-center rounded border px-1.5 py-0.5 dark:border-neutral-800">
                {chipEquip}
              </span>
            )}
            <span className="inline-flex items-center rounded border px-1.5 py-0.5 dark:border-neutral-800">
              Posición {p.position}
            </span>
          </div>
        </div>

        <button
          onClick={() => delPrescM.mutate(p.id)}
          className="text-xs text-red-600 rounded px-2 py-0.5 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
        >
          Eliminar
        </button>
      </div>

      <div className="text-xs text-gray-700 dark:text-neutral-200 mt-2">
        Series: <span className="font-medium">{p.series}</span>
        <span className="mx-2">·</span>
        Reps: <span className="font-medium">{p.reps}</span>
        {typeof p.rest_sec === 'number' && (
          <>
            <span className="mx-2">·</span>
            Descanso: <span className="font-medium">{p.rest_sec}s</span>
          </>
        )}
        {p.to_failure && (
          <>
            <span className="mx-2">·</span>
            A fallo
          </>
        )}
      </div>
    </li>
  )
})}

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

      <Modal open={openNewWeek} onClose={() => setOpenNewWeek(false)} title="Nueva semana">
        <form
          className="space-y-3"
          onSubmit={(e) => {
            e.preventDefault()
            if (!selectedProgram) return
            const fd = new FormData(e.currentTarget as HTMLFormElement)
            addWeekM.mutate({
              programId: selectedProgram.id,
              week_index: Number(fd.get('week_index') || nextWeekIndex),
              title: (fd.get('title') as string) || null,
            })
          }}
        >
          <label className="text-sm block">
            <div className="mb-1">Índice *</div>
            <input
              name="week_index"
              type="number"
              min={1}
              defaultValue={(weeksQ.data?.length ?? 0) + 1}
              required
              className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
            />
          </label>
          <label className="text-sm block">
            <div className="mb-1">Título (opcional)</div>
            <input
              name="title"
              className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
              placeholder="Semana de fuerza"
            />
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
            e.preventDefault();
            if (!selectedProgram?.id || !selectedWeek?.id) {
              show({ type: 'error', message: 'Selecciona una semana primero' });
              return;
            }
            const fd = new FormData(e.currentTarget as HTMLFormElement);
            addDayM.mutate({
              programId: selectedProgram.id,
              weekId: selectedWeek.id,             // UUID de la semana
              day_index: Number(fd.get('day_index') || nextDayIndex),
              notes: (fd.get('notes') as string) || null,
            });
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
            <button type="submit" disabled={!canAddDay || addDayM.isPending} className="text-sm rounded px-3 py-1 border">
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
          <div className="space-y-3">
            <label className="text-sm block">
              <div className="mb-1">Ejercicio *</div>
              <select
                className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
                value={selectedExerciseId}
                onChange={(e) => setSelectedExerciseId(e.target.value)}
                required
              >
                <option value="">— Selecciona —</option>
                {exercises.map(ex => (
                  <option key={ex.id} value={ex.id}>
                    {ex.name} · {ex.primary_muscle}
                  </option>
                ))}
              </select>
            </label>

            <PrescriptionForm
              exerciseId={selectedExerciseId}                 // <- ahora sí existe
              defaultValues={{ position: nextPosition }}
              submitting={addPrescM.isPending}
              onCancel={() => setOpenNewPresc(false)}
              onSubmit={async (vals) => {
                if (!selectedExerciseId) return
                await addPrescM.mutateAsync({ dayId: selectedDay.id, values: vals })
              }}
            />
          </div>
        ) : (
          <div className="text-sm text-gray-500">Selecciona un día</div>
        )}
      </Modal>

    </div>
  )
}
