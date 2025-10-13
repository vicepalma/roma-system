import type { CoachLink } from '@/types/coach'

export default function LinksGrid({ links }: { links: CoachLink[] }) {
  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {links.map((l, idx) => (
        <a
          key={idx}
          href={l.href}
          className="rounded-lg border bg-white p-4 hover:bg-gray-50 transition"
        >
          <div className="font-semibold">{l.title}</div>
          {l.description && <div className="text-sm text-gray-500 mt-1">{l.description}</div>}
        </a>
      ))}
    </div>
  )
}
