import { createContext, useContext, useEffect, useMemo, useState, ReactNode } from 'react'

export type Toast = { id: string; title?: string; message: string; type?: 'info' | 'success' | 'error'; timeout?: number }

type ToastCtx = {
  toasts: Toast[]
  show: (t: Omit<Toast, 'id'>) => void
  remove: (id: string) => void
}

const Ctx = createContext<ToastCtx | null>(null)

export function useToast() {
  const ctx = useContext(Ctx)
  if (!ctx) throw new Error('ToastProvider missing')
  return ctx
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([])

  const show: ToastCtx['show'] = (t) => {
    const id = crypto.randomUUID()
    const toast: Toast = { id, type: 'info', timeout: 3500, ...t }
    setToasts((prev) => [...prev, toast])
    if (toast.timeout && toast.timeout > 0) {
      setTimeout(() => remove(id), toast.timeout)
    }
  }

  const remove = (id: string) => setToasts((prev) => prev.filter((t) => t.id !== id))

  const value = useMemo(() => ({ toasts, show, remove }), [toasts])

  return (
    <Ctx.Provider value={value}>
      {children}
      <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
        {toasts.map((t) => (
          <div
            key={t.id}
            className={[
              'min-w-[260px] max-w-sm rounded-md border bg-white px-3 py-2 shadow',
              t.type === 'success' ? 'border-green-200' :
              t.type === 'error'   ? 'border-red-200'   : 'border-gray-200'
            ].join(' ')}
          >
            {t.title && <div className="text-sm font-medium">{t.title}</div>}
            <div className="text-sm text-gray-700">{t.message}</div>
          </div>
        ))}
      </div>
    </Ctx.Provider>
  )
}
