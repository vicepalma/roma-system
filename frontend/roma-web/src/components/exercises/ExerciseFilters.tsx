import { useState } from 'react'

type Props = {
  value: { q: string; tags: string[]; match: 'any'|'all' }
  onChange: (v: Props['value']) => void
  onSearch: () => void
}
export default function ExerciseFilters({ value, onChange, onSearch }: Props) {
  const [tagInput, setTagInput] = useState('')

  const addTag = () => {
    const t = tagInput.trim()
    if (!t) return
    if (!value.tags.includes(t)) onChange({ ...value, tags: [...value.tags, t] })
    setTagInput('')
  }
  const removeTag = (t: string) =>
    onChange({ ...value, tags: value.tags.filter(x => x !== t) })

  return (
    <div className="flex flex-col gap-3 md:flex-row md:items-end">
      <div className="flex-1">
        <label className="block text-xs text-gray-600 dark:text-neutral-300 mb-1">Buscar</label>
        <input
          value={value.q}
          onChange={(e) => onChange({ ...value, q: e.target.value })}
          placeholder="Nombre, músculo, etc."
          className="w-full rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
        />
      </div>

      <div>
        <label className="block text-xs text-gray-600 dark:text-neutral-300 mb-1">Tags</label>
        <div className="flex gap-2">
          <input
            value={tagInput}
            onChange={(e) => setTagInput(e.target.value)}
            placeholder="agrega un tag"
            className="rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
          />
          <button onClick={addTag} className="rounded border px-3 py-2 text-sm bg-white dark:bg-neutral-900 dark:border-neutral-800">
            Agregar
          </button>
        </div>
        <div className="mt-2 flex flex-wrap gap-2">
          {value.tags.map(t => (
            <span key={t} className="text-xs rounded border px-2 py-0.5 dark:border-neutral-800">
              {t}{' '}
              <button onClick={() => removeTag(t)} className="text-blue-600 hover:underline">x</button>
            </span>
          ))}
        </div>
      </div>

      <div>
        <label className="block text-xs text-gray-600 dark:text-neutral-300 mb-1">Match</label>
        <select
          value={value.match}
          onChange={(e) => onChange({ ...value, match: e.target.value as 'any'|'all' })}
          className="rounded border px-3 py-2 dark:bg-neutral-900 dark:border-neutral-800"
        >
          <option value="any">Cualquiera</option>
          <option value="all">Todas</option>
        </select>
      </div>

      <button onClick={onSearch} className="rounded border px-4 py-2 bg-white dark:bg-neutral-900 dark:border-neutral-800">
        Buscar
      </button>
    </div>
  )
}
