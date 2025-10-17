import { useForm, type SubmitHandler, type Resolver } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'

const schema = z.object({
  reps: z.coerce.number().int().min(1).max(200),
  weight: z.coerce.number().optional().nullable(),
  rpe: z.coerce.number().min(0).max(10).optional().nullable(),
  notes: z.string().optional().nullable(),
})
export type Values = z.infer<typeof schema>

type Props = {
  defaultValues?: Partial<Values>
  onSubmit: (values: Values) => void | Promise<void>
  onCancel: () => void
}

export default function LogSetForm({ defaultValues, onSubmit, onCancel }: Props) {
      const { register, handleSubmit, formState: { isSubmitting, errors } } =
    useForm<Values>({
      resolver: zodResolver(schema) as unknown as Resolver<Values>,
      defaultValues,
    })

const onValid: SubmitHandler<Values> = (values) => onSubmit(values)

  return (
    <form onSubmit={handleSubmit(onValid)} className="space-y-3">
      <label className="text-sm block">
        <div className="mb-1">Reps *</div>
        <input type="number" min={1} {...register('reps', { valueAsNumber: true })} className="w-full border rounded px-2 py-1" />
        {errors.reps && <div className="text-xs text-red-600">Repeticiones inválidas</div>}
      </label>

      <div className="grid grid-cols-2 gap-3">
        <label className="text-sm block">
          <div className="mb-1">Peso</div>
          <input type="number" step="0.5" {...register('weight', { setValueAs: v => (v === '' ? null : Number(v)) })} className="w-full border rounded px-2 py-1" />
        </label>
        <label className="text-sm block">
          <div className="mb-1">RPE</div>
          <input type="number" step="0.5" min={0} max={10} {...register('rpe', { setValueAs: v => (v === '' ? null : Number(v)) })} className="w-full border rounded px-2 py-1" />
        </label>
      </div>

      <label className="text-sm block">
        <div className="mb-1">Notas</div>
        <textarea rows={2} {...register('notes')} className="w-full border rounded px-2 py-1" />
      </label>

      <div className="flex items-center gap-2">
        <button type="submit" disabled={isSubmitting}>Guardar</button>
        <button type="button" onClick={onCancel} className="text-sm text-gray-600 hover:underline">
          Cancelar
        </button>
      </div>
    </form>
  )
}
