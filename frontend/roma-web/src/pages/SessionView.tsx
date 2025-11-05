import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useEffect, useMemo, useState } from 'react'
import { useToast } from '@/components/toast/ToastProvider'
import { getSession, addSet, deleteSet, patchSession, endSession } from '@/services/sessions'
import { listPrescriptions } from '@/services/programs'
import { AssignmentDay } from '@/types/assignments'
import { listAssignmentDays, getMyActiveAssignment } from '@/services/assignments'
import Modal from '@/components/ui/Modal'
import LogSetForm from '@/components/forms/LogSetForm'
import RestTimer from '@/components/timers/RestTimer'
import type { SessionSet } from '@/types/sessions'

export default function SessionView() {
  const { id: sessionIdParam = '' } = useParams()
  const navigate = useNavigate()
  const { show } = useToast()
  const qc = useQueryClient()

  const [openLog, setOpenLog] = useState<{ open: boolean; prescription_id?: string }>({ open: false })
  const [openPicker, setOpenPicker] = useState(false)
  const [openPickOther, setOpenPickOther] = useState(false)
  const [otherDayId, setOtherDayId] = useState<string | null>(null)

  const disabledBtn = "opacity-60 cursor-not-allowed"
  const baseBtn = "text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"

  const sessQ = useQuery({
    queryKey: ['session', sessionIdParam],
    queryFn: () => getSession(sessionIdParam),
    enabled: !!sessionIdParam,
  })

  const sess = sessQ.data?.session ?? null
  const sets = sessQ.data?.sets ?? []
  const cardio = sessQ.data?.cardio ?? []
  const isClosed = sess?.status === 'closed'

  const activeQ = useQuery({
    queryKey: ['me', 'assignment', 'active'],
    queryFn: getMyActiveAssignment,
    staleTime: 30_000,
  })
  const active = activeQ.data || null

  const showMismatchBanner = Boolean(
    sess && active?.id && sess.assignment_id !== active.id
  )

  const daysQ = useQuery({
    queryKey: ['assignment', active?.id, 'days'],
    queryFn: () => listAssignmentDays(active!.id),
    enabled: !!active?.id,
    staleTime: 30_000,
  })
  const days = daysQ.data ?? []

  const sessionDayId = sess?.day_id ?? null
  const [selectedDayId, setSelectedDayId] = useState<string | null>(null)

  useEffect(() => {
    if (sessionDayId) setSelectedDayId(sessionDayId)
    else if (!sessionDayId && days.length) {
      const pref = days.find(d => (d as any).is_session_day) ?? days[0]
      setSelectedDayId(pref?.id ?? null)
    }
  }, [sessionDayId, days])

  const effectiveDayId = selectedDayId ?? sessionDayId ?? undefined

  const prescQ = useQuery({
    queryKey: ['day', effectiveDayId, 'prescriptions'],
    queryFn: () => listPrescriptions(effectiveDayId!),
    enabled: !!effectiveDayId,
    staleTime: 30_000,
  })
  const presc = prescQ.data ?? []



  const safeAdd = async (vals: {
    prescription_id: string; set_index: number; reps: number; weight?: number | null; rpe?: number | null; to_failure?: boolean;
  }) => {
    if (isClosed) return
    await mAdd.mutateAsync(vals)
  }
  const safeDel = (id: string) => {
    if (isClosed) return
    mDel.mutate(id)
  }

  // 1) Registrar set — si es el primero y no hay performed_at, lo establecemos
  const mAdd = useMutation({
    mutationFn: async (vals: {
      prescription_id: string
      set_index: number
      reps: number
      weight?: number | null
      rpe?: number | null
      to_failure?: boolean
    }) => {
      const payload = {
        prescription_id: vals.prescription_id,
        set_index: vals.set_index,
        reps: vals.reps,
        weight: vals.weight ?? null,
        rpe: vals.rpe ?? null,
        to_failure: !!vals.to_failure,
      } satisfies Omit<SessionSet, 'id' | 'session_id' | 'created_at'>

      // lazy start: set performed_at si está vacío
      if (sess && !sess.performed_at) {
        await patchSession(sessionIdParam, { performed_at: new Date().toISOString() })
      }
      return addSet(sessionIdParam, payload)
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['session', sessionIdParam] })
      setOpenLog({ open: false })
      show({ type: 'success', message: 'Set registrado' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo registrar el set' }),
  })

  const mDel = useMutation({
    mutationFn: (setId: string) => deleteSet(sessionIdParam, setId),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['session', sessionIdParam] })
      show({ type: 'success', message: 'Set eliminado' })
    },
  })

  // 2) Cambiar día base (sólo si no hay sets). Opcional
  const mPatchDay = useMutation({
    mutationFn: (newDayId: string) => patchSession(sessionIdParam, { day_id: newDayId }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ['session', sessionIdParam] })
      await qc.invalidateQueries({ queryKey: ['day', effectiveDayId, 'prescriptions'] })
      setOpenPicker(false)
      show({ type: 'success', message: 'Día actualizado' })
    },
    onError: () => show({ type: 'error', message: 'No se pudo cambiar el día' }),
  })

  // 3) Terminar entrenamiento
  const mEnd = useMutation({
    mutationFn: () => endSession(sessionIdParam),
    onSuccess: async () => {
      show({ type: 'success', message: 'Sesión finalizada' })
      // limpia caches clave
      await qc.invalidateQueries({ queryKey: ['me', 'session', 'active'] })
      await qc.invalidateQueries({ queryKey: ['session', sessionIdParam] })
      navigate('/sessions')
    },
    onError: () => show({ type: 'error', message: 'No se pudo finalizar la sesión' }),
  })

  const selectedDay = useMemo(
    () => days.find(d => d.id === (effectiveDayId ?? '')) || null,
    [days, effectiveDayId]
  )

  const prescIdsOfEffectiveDay = new Set(presc.map((p: any) => p.id))
  const isOffPlan = (prescriptionId: string) => !prescIdsOfEffectiveDay.has(prescriptionId)

  return (
    <div className="space-y-4">
      {/* Banner de sesión cerrada */}
      {isClosed && (
        <div className="rounded-lg border bg-gray-50 dark:bg-neutral-800/50 dark:border-neutral-700 p-3 flex items-center justify-between">
          <div className="text-sm">
            <b>Sesión cerrada.</b> Esta fue tu última sesión registrada y está en modo lectura.
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => navigate('/sessions')}
              className={baseBtn}
            >
              Ir a Sesiones
            </button>
          </div>
        </div>
      )}

      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="text-sm text-gray-600 dark:text-neutral-300">
            {sess?.performed_at ? new Date(sess.performed_at).toLocaleString() : '—'}
            {isClosed && sess?.ended_at ? ` · Finalizada: ${new Date(sess.ended_at).toLocaleString()}` : ''}
          </div>
          <h1 className="text-xl font-semibold">Sesión</h1>
          {selectedDay && (
            <div className="text-sm text-gray-600 dark:text-neutral-300">
              Día: {selectedDay.title?.trim() ? selectedDay.title : `Día ${selectedDay.day_index}`}
            </div>
          )}
        </div>
        {/* RestTimer sólo si está abierta */}
        {!isClosed && <RestTimer initial={60} key={sets.length} />}
      </div>

      {/* Banner mismatch de programa: se puede seguir mostrando; si cerrada, sólo permite ir a /sessions */}
      {showMismatchBanner && (
        <div className="rounded-lg border bg-amber-50 dark:bg-amber-900/20 dark:border-amber-900 p-3 text-sm flex items-center justify-between gap-3">
          <div>
            Estás viendo una sesión de <b>otro programa</b>.
            {isClosed
              ? ' Como la sesión está cerrada, inicia una nueva desde la vista de Sesiones.'
              : ' Termina la sesión o inicia una nueva con el programa activo.'}
          </div>
          {isClosed ? (
            <button
              onClick={() => navigate('/sessions')}
              className={baseBtn}
            >
              Ir a Sesiones
            </button>
          ) : (
            <button
              onClick={() => setOpenPicker(true)}
              className={baseBtn}
            >
              Iniciar con programa activo
            </button>
          )}
        </div>
      )}

      {/* Prescripciones */}
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
        <div className="flex items-center justify-between mb-2">
          <div className="font-semibold">Prescripciones {prescQ.isLoading ? '…' : ''}</div>
          {/* + Set de otro día - deshabilitado si cerrada */}
          <button
            onClick={() => { setOtherDayId(null); setOpenPickOther(true) }}
            className={`${baseBtn} ${isClosed ? disabledBtn : ''}`}
            disabled={isClosed || !days.length}
          >
            + Set de otro día
          </button>
        </div>

        {presc.length === 0 && !prescQ.isLoading && (
          <div className="text-sm text-gray-500">Sin prescripciones en este día</div>
        )}

        <ul className="mt-2 space-y-2">
          {presc.map((p: any) => (
            <li key={p.id} className="rounded border px-3 py-2 dark:border-neutral-800">
              <div className="flex items-center justify-between">
                <div className="text-sm font-medium">{p.exercise_name || p.exercise_id}</div>
                {/* + Set - deshabilitado si cerrada */}
                <button
                  onClick={() => !isClosed && setOpenLog({ open: true, prescription_id: p.id })}
                  className={`${baseBtn} ${isClosed ? disabledBtn : ''}`}
                  disabled={isClosed}
                >
                  + Set
                </button>
              </div>
              <div className="text-xs text-gray-600 mt-1">
                Series: {p.series} · Reps: {p.reps}
                {p.rest_sec ? ` · Descanso: ${p.rest_sec}s` : ''}{p.to_failure ? ' · A fallo' : ''} · Posición: {p.position}
              </div>
            </li>
          ))}
        </ul>
      </div>

      {/* Sets */}
      <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
        <div className="font-semibold mb-2">Sets de la sesión</div>
        {sets.length === 0 ? (
          <div className="text-sm text-gray-500">Aún no hay sets</div>
        ) : (
          <ul className="mt-2 space-y-2">
            {sets.map((s: any) => (
              <li key={s.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                <div className="flex items-center justify-between">
                  <div className="text-sm">
                    <span className="font-medium">Set {s.set_index}</span>
                    {s.exercise_name ? ` — ${s.exercise_name}` : ''}
                    {` — Reps: ${s.reps}`}
                    {s.weight != null ? ` · Peso: ${s.weight}` : ''}
                    {s.rpe != null ? ` · RPE: ${s.rpe}` : ''}
                    {s.to_failure ? ' · A fallo' : ''}
                    {isOffPlan(s.prescription_id) && (
                      <span className="ml-2 text-[11px] px-1.5 py-0.5 rounded bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-200">
                        Fuera del día
                      </span>
                    )}
                  </div>
                  {/* Eliminar - deshabilitado si cerrada */}
                  <button
                    onClick={() => safeDel(s.id)}
                    className={`text-xs text-red-600 rounded px-2 py-0.5 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800 ${isClosed ? disabledBtn : ''}`}
                    disabled={isClosed}
                  >
                    Eliminar
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}

        {/* Entrenamiento terminado - oculto si cerrada */}
        {!isClosed && (
          <div className="pt-3">
            <button
              onClick={() => mEnd.mutate()}
              disabled={mEnd.isPending}
              className={baseBtn}
            >
              {mEnd.isPending ? 'Cerrando…' : 'Entrenamiento terminado'}
            </button>
          </div>
        )}
      </div>

      {/* Cardio igual en lectura */}
      {!!cardio.length && (
        <div className="rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4">
          <div className="font-semibold mb-2">Cardio</div>
          <ul className="mt-2 space-y-1 text-sm">
            {cardio.map((c: any) => (
              <li key={c.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                {c.modality} — {c.minutes} min
                {c.target_hr_min ? ` · HR ${c.target_hr_min}-${c.target_hr_max ?? ''}` : ''}
                {c.notes ? ` · ${c.notes}` : ''}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Modal registrar set — bloquea submit si cerrada */}
      <Modal open={openLog.open} onClose={() => setOpenLog({ open: false })} title="Registrar set">
        {openLog.prescription_id ? (
          <LogSetForm
            defaultValues={{ reps: 10 }}
            onCancel={() => setOpenLog({ open: false })}
            onSubmit={async (vals) => {
              if (isClosed) return
              const prescId = openLog.prescription_id!
              const nextIndex = (sets?.filter((s: any) => s.prescription_id === prescId).length ?? 0) + 1
              await safeAdd({
                prescription_id: prescId,
                set_index: nextIndex,
                reps: vals.reps,
                weight: vals.weight ?? null,
                rpe: vals.rpe ?? null,
                to_failure: !!vals.to_failure,
              })
            }}
          />
        ) : (
          <div className="text-sm text-gray-500">Selecciona una prescripción</div>
        )}
      </Modal>

      {/* Modal cambiar día — botón que lo abre ya va disabled si isClosed */}
      <Modal open={openPicker} onClose={() => setOpenPicker(false)} title="Seleccionar día del programa activo">
        {!days.length ? (
          <div className="text-sm text-gray-500">No hay días disponibles.</div>
        ) : (
          <ul className="space-y-2">
            {days.map((d: AssignmentDay) => (
              <li key={d.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                <div className="flex items-center justify-between">
                  <div>
                    <div className="text-sm font-medium">
                      {d.title?.trim() ? d.title : `Día ${d.day_index}`}
                    </div>
                    <div className="text-xs text-gray-600 dark:text-neutral-300">
                      {d.exercise_names?.length ? d.exercise_names.join(', ') : '—'}
                    </div>
                  </div>
                  <button
                    onClick={() => mPatchDay.mutate(d.id)}
                    disabled={sets.length > 0 || mPatchDay.isPending}
                    className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                  >
                    {mPatchDay.isPending ? 'Cambiando…' : 'Usar este día'}
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </Modal>

      {/* Modal + Set de otro día — abrirlo está bloqueado arriba si cerrada */}
      <Modal open={openPickOther} onClose={() => setOpenPickOther(false)} title="Agregar set de otro día">
        {!days.length ? (
          <div className="text-sm text-gray-500">No hay días disponibles.</div>
        ) : (
          <OtherDayPicker
            excludeDayId={effectiveDayId}
            days={days}
            onPickPresc={(prescId) => {
              setOpenPickOther(false)
              setOpenLog({ open: true, prescription_id: prescId })
            }}
          />
        )}
      </Modal>
    </div>
  )
}

function OtherDayPicker({
  days, excludeDayId, onPickPresc,
}: { days: AssignmentDay[]; excludeDayId?: string; onPickPresc: (id: string) => void }) {
  const [dayId, setDayId] = useState<string>('')
  const prescQ = useQuery({
    queryKey: ['day', dayId, 'prescriptions'],
    queryFn: () => listPrescriptions(dayId),
    enabled: !!dayId,
    staleTime: 30_000,
  })
  const items = prescQ.data ?? []

  return (
    <div className="space-y-3">
      <label className="text-sm block">
        <div className="mb-1">Elegir día</div>
        <select
          className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800"
          value={dayId}
          onChange={(e) => setDayId(e.target.value)}
        >
          <option value="">Selecciona…</option>
          {days.filter(d => d.id !== excludeDayId).map(d => (
            <option key={d.id} value={d.id}>
              {d.title?.trim() ? d.title : `Día ${d.day_index}`}
            </option>
          ))}
        </select>
      </label>

      {dayId && (
        prescQ.isLoading ? (
          <div className="text-sm text-gray-500">Cargando prescripciones…</div>
        ) : !items.length ? (
          <div className="text-sm text-gray-500">Este día no tiene prescripciones.</div>
        ) : (
          <ul className="space-y-2">
            {items.map((p: any) => (
              <li key={p.id} className="rounded border px-3 py-2 dark:border-neutral-800">
                <div className="flex items-center justify-between">
                  <div className="text-sm">
                    <div className="font-medium">{p.exercise_name || p.exercise_id}</div>
                    <div className="text-xs text-gray-600">
                      Series: {p.series} · Reps: {p.reps}{p.rest_sec ? ` · Descanso: ${p.rest_sec}s` : ''}{p.to_failure ? ' · A fallo' : ''}
                    </div>
                  </div>
                  <button
                    onClick={() => onPickPresc(p.id)}
                    className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
                  >
                    Usar
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )
      )}
    </div>
  )
}
