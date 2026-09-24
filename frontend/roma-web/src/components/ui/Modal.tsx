import { ReactNode } from 'react'

export default function Modal({
  open, onClose, title, children,
}: { open: boolean; onClose: () => void; title?: string; children: ReactNode }) {
  if (!open) return null
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />
      <div className="relative z-10 mx-3 max-h-[calc(100vh-2rem)] w-full max-w-md overflow-y-auto rounded-lg border bg-white p-3 shadow-lg dark:border-neutral-800 dark:bg-neutral-900 sm:mx-0 sm:p-4">
        {title && <div className="text-sm font-semibold mb-2">{title}</div>}
        {children}
      </div>
    </div>
  )
}
