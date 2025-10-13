import { ReactNode } from 'react'

export default function Modal({
  open, onClose, title, children,
}: { open: boolean; onClose: () => void; title?: string; children: ReactNode }) {
  if (!open) return null
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="absolute inset-0 bg-black/40" onClick={onClose} />
      <div className="relative z-10 w-full max-w-md rounded-lg border bg-white dark:bg-neutral-900 dark:border-neutral-800 p-4 shadow-lg">
        {title && <div className="text-sm font-semibold mb-2">{title}</div>}
        {children}
      </div>
    </div>
  )
}
