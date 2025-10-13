export function RangePicker({ value, onChange }: { value: number; onChange: (n:number)=>void }) {
  const opts = [7, 14, 30]
  return (
    <div className="inline-flex gap-2">
      {opts.map(n => (
        <button
          key={n}
          onClick={() => onChange(n)}
          className={
            value === n
              ? 'text-sm px-3 py-1 rounded border bg-black text-white'
              : 'text-sm px-3 py-1 rounded border bg-white hover:bg-gray-50 dark:bg-neutral-800 dark:border-neutral-700'
          }
        >
          {n}d
        </button>
      ))}
    </div>
  )
}
