import type { CoachLink } from '@/types/coach'

type Props = { items: CoachLink[] }

export default function LinksGrid({ items }: Props) {
  return (
    <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {items.map(l => (
        <div key={l.id} className="rounded border p-3 dark:border-neutral-800">
          <div className="text-sm font-semibold">Vínculo</div>
          <div className="text-xs text-gray-500 mt-1">Discípulo: {l.disciple_id}</div>
          <div className="mt-2 inline-flex rounded px-2 py-0.5 text-xs
                          bg-gray-100 text-gray-700 dark:bg-neutral-800 dark:text-neutral-300">
            {l.status}
          </div>
          <div className="text-[11px] text-gray-500 font-mono mt-2">#{l.id}</div>
        </div>
      ))}
    </div>
  )
}
