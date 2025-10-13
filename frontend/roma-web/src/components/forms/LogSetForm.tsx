import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'

const schema = z.object({
  reps: z.coerce.number().int().min(1).max(200),
  weight: z.coerce.number().min(0).max(2000).optional().nullable(),
  rpe: z.coerce.number().min(1).max(10).optional().nullable(),
  notes: z.string().max(500).optional().nullable(),
})
type Values = z.infer<typeof schema>

export default function LogSetForm({
  defaultValues,
  onSubmit,
  onCancel,
}: {
  defaultValues?: Partial<Values>
  onSubmit: (values: Values) => void | Promise<void>
  onCancel: () => void
}) {
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: { reps: 10, weight: undefined, rpe: undefined, notes: '', ...defaultValues },
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
      <div>
        <label className="block text-sm mb-1">Reps</label>
        <input type="number" {...register('reps')}
          className="w-full rounded-md border px-3 py-2 bg-white dark:bg-neutral-900 dark:border-neutral-800" />
        {errors.reps && <div className="text-xs text-red-600 mt-1">Reps inválidas</div>}
      </div>
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="block text-sm mb-1">Peso (kg)</label>
          <input type="number" step="0.5" {...register('weight')}
            className="w-full rounded-md border px-3 py-2 bg-white dark:bg-neutral-900 dark:border-neutral-800" />
          {errors.weight && <div className="text-xs text-red-600 mt-1">Peso inválido</div>}
        </div>
        <div>
          <label className="block text-sm mb-1">RPE (1-10)</label>
          <input type="number" step="0.5" {...register('rpe')}
            className="w-full rounded-md border px-3 py-2 bg-white dark:bg-neutral-900 dark:border-neutral-800" />
          {errors.rpe && <div className="text-xs text-red-600 mt-1">RPE inválido</div>}
        </div>
      </div>
      <div>
        <label className="block text-sm mb-1">Notas</label>
        <textarea rows={3} {...register('notes')}
          className="w-full rounded-md border px-3 py-2 bg-white dark:bg-neutral-900 dark:border-neutral-800" />
        {errors.notes && <div className="text-xs text-red-600 mt-1">Notas demasiado largas</div>}
      </div>

      <div className="flex justify-end gap-2 pt-2">
        <button type="button" onClick={onCancel}
          className="rounded-md border px-3 py-1.5 text-sm bg-white dark:bg-neutral-900 dark:border-neutral-800">
          Cancelar
        </button>
        <button type="submit" disabled={isSubmitting}
          className="rounded-md border px-3 py-1.5 text-sm bg-black text-white">
          {isSubmitting ? 'Guardando…' : 'Guardar'}
        </button>
      </div>
    </form>
  )
}
