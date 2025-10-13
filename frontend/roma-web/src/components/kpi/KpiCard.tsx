export function KpiCard({ label, value, sub }: { label: string; value: string; sub?: string }) {
  return (
    <div className="rounded-lg border bg-white dark:bg-neutral-800 dark:border-neutral-700 p-4">
      <div className="text-sm text-gray-500 dark:text-neutral-400">{label}</div>
      <div className="text-2xl font-semibold">{value}</div>
      {sub && <div className="text-xs text-gray-500 dark:text-neutral-400">{sub}</div>}
    </div>
  )
}
