import { ButtonHTMLAttributes } from 'react'
import { clsx } from 'clsx'

type Props = ButtonHTMLAttributes<HTMLButtonElement> & { variant?: 'default' | 'outline' }
export function Button({ className, variant = 'default', ...props }: Props) {
  return (
    <button
      className={clsx(
        'inline-flex items-center justify-center rounded-md text-sm font-medium h-10 px-4 py-2 transition-colors',
        variant === 'default' && 'bg-black text-white hover:bg-black/90',
        variant === 'outline' && 'border border-gray-300 bg-white hover:bg-gray-50',
        className
      )}
      {...props}
    />
  )
}
