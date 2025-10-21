import { useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import type { Exercise } from '@/types/exercises'

const schema = z.object({
  name: z.string().min(2, 'Mínimo 2 caracteres'),
  primary_muscle: z.string().min(2, 'Requerido'),
  equipment: z.string().optional().nullable(),
  notes: z.string().optional().nullable(),
  // Campo de UI para escribir en una sola línea:
  tagsText: z.string().optional(), // "press, chest, machine"
})

type Values = z.infer<typeof schema>

type Props = {
  defaultValues?: Partial<Exercise>
  onSubmit: (v: {
    name: string
    primary_muscle: string
    equipment?: string | null
    notes?: string | null
    tags?: string[]
  }) => void | Promise<void>
  onCancel: () => void
  submitting?: boolean
}

export default function ExerciseForm({ defaultValues, onSubmit, onCancel, submitting }: Props) {
  const initial: Values = useMemo(() => {
    const tagsText = Array.isArray(defaultValues?.tags) ? defaultValues!.tags!.join(', ') : ''
    return {
      name: defaultValues?.name ?? '',
      primary_muscle: defaultValues?.primary_muscle ?? '',
      equipment: defaultValues?.equipment ?? '',
      notes: defaultValues?.notes ?? '',
      tagsText,
    }
  }, [defaultValues])

  const { register, handleSubmit, formState: { errors } } = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: initial,
  })

  const submit = (v: Values) => {
    const tags = (v.tagsText ?? '')
      .split(',')
      .map(s => s.trim())
      .filter(Boolean)

    onSubmit({
      name: v.name.trim(),
      primary_muscle: v.primary_muscle.trim(),
      equipment: v.equipment ? v.equipment.trim() : null,
      notes: v.notes ? v.notes.trim() : null,
      tags,
    })
  }

  return (
    <form onSubmit={handleSubmit(submit)} className="space-y-3">
      <div>
        <label className="block text-sm font-medium mb-1">Nombre *</label>
        <input {...register('name')} className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
        {errors.name && <p className="text-xs text-red-600 mt-1">{errors.name.message}</p>}
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">Músculo principal *</label>
        <input {...register('primary_muscle')} className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
        {errors.primary_muscle && <p className="text-xs text-red-600 mt-1">{errors.primary_muscle.message}</p>}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div>
          <label className="block text-sm font-medium mb-1">Equipo</label>
          <input {...register('equipment')} className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Tags (separadas por coma)</label>
          <input {...register('tagsText')} placeholder="press, machine, chest" className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium mb-1">Notas</label>
        <textarea rows={3} {...register('notes')} className="w-full border rounded px-2 py-1 text-sm dark:bg-neutral-900 dark:border-neutral-800" />
      </div>

      <div className="flex items-center justify-end gap-2 pt-2">
        <button type="button" onClick={onCancel} className="text-sm rounded px-3 py-1 border bg-white hover:bg-gray-50 dark:bg-neutral-900 dark:border-neutral-800">
          Cancelar
        </button>
        <button type="submit" disabled={submitting} className="text-sm rounded px-3 py-1 border bg-black text-white disabled:opacity-60">
          {submitting ? 'Guardando…' : 'Guardar'}
        </button>
      </div>
    </form>
  )
}
