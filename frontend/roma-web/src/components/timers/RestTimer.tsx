import { useEffect, useMemo, useRef, useState } from 'react'

export default function RestTimer({
  initial = 90,               // segundos por defecto
  autoRestart = false,        // reiniciar al terminar
  onFinish,
}: {
  initial?: number
  autoRestart?: boolean
  onFinish?: () => void
}) {
  const [seconds, setSeconds] = useState<number>(initial)
  const [running, setRunning] = useState<boolean>(false)
  const lastTick = useRef<number | null>(null)
  const rafId = useRef<number | null>(null)

  // mm:ss
  const label = useMemo(() => {
    const m = Math.floor(seconds / 60).toString().padStart(2, '0')
    const s = Math.floor(seconds % 60).toString().padStart(2, '0')
    return `${m}:${s}`
  }, [seconds])

  // ticker usando requestAnimationFrame para suavidad
  useEffect(() => {
    if (!running) return

    const loop = (ts: number) => {
      if (lastTick.current == null) lastTick.current = ts
      const dt = (ts - lastTick.current) / 1000
      lastTick.current = ts
      setSeconds(prev => {
        const next = Math.max(0, prev - dt)
        if (next === 0) {
          // vibración suave si está disponible
          if (navigator?.vibrate) navigator.vibrate(150)
          onFinish?.()
          if (autoRestart) {
            // reinicia con el valor inicial
            lastTick.current = null
            return initial
          } else {
            // detiene
            setRunning(false)
          }
        }
        return next
      })
      rafId.current = requestAnimationFrame(loop)
    }
    rafId.current = requestAnimationFrame(loop)

    return () => {
      if (rafId.current) cancelAnimationFrame(rafId.current)
      rafId.current = null
      lastTick.current = null
    }
  }, [running, autoRestart, initial, onFinish])

  // helpers
  const startPause = () => setRunning(r => !r)
  const reset = () => { setSeconds(initial); setRunning(false) }
  const setPreset = (s: number) => { setSeconds(s); setRunning(false) }

  const bump = (delta: number) => {
    setSeconds(prev => Math.max(0, prev + delta))
  }

  return (
    <div className="inline-flex items-center gap-2">
      <div className="rounded-lg border px-3 py-2 bg-white dark:bg-neutral-900 dark:border-neutral-800">
        <div className="text-xs text-gray-500 dark:text-neutral-300">Descanso</div>
        <div className="text-2xl font-semibold tabular-nums text-center">{label}</div>
        <div className="mt-2 flex items-center gap-2">
          <button
            onClick={startPause}
            className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
          >
            {running ? 'Pausar' : 'Iniciar'}
          </button>
          <button
            onClick={reset}
            className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
          >
            Reiniciar
          </button>
        </div>
        <div className="mt-2 flex items-center gap-2">
          <button
            onClick={() => bump(-10)}
            className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
            title="-10s"
          >
            −10s
          </button>
          <button
            onClick={() => bump(10)}
            className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
            title="+10s"
          >
            +10s
          </button>
        </div>
        <div className="mt-2 flex items-center gap-2">
          {[60, 90, 120].map(p => (
            <button
              key={p}
              onClick={() => setPreset(p)}
              className="text-xs rounded px-2 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800"
            >
              {p}s
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}
