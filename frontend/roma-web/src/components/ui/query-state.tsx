type QueryStateProps = {
  title: string
  detail?: string
  tone?: 'muted' | 'error'
  className?: string
}

export function QueryState({ title, detail, tone = 'muted', className = '' }: QueryStateProps) {
  const toneClass = tone === 'error'
    ? 'border-red-200 bg-red-50 text-red-700 dark:border-red-900/60 dark:bg-red-950/20 dark:text-red-300'
    : 'border-gray-200 bg-gray-50 text-gray-700 dark:border-neutral-800 dark:bg-neutral-950/40 dark:text-neutral-300'

  return (
    <div className={`rounded border px-3 py-3 text-sm ${toneClass} ${className}`}>
      <div className="font-medium">{title}</div>
      {detail && <div className="mt-1 text-xs opacity-80">{detail}</div>}
    </div>
  )
}
